/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package managers

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"path"
	"sync"

	minio "github.com/minio/minio-go/v7"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	pipelinesmeta "github.com/tinyzimmer/gst-pipeline-operator/apis/meta/v1"
	pipelinesv1 "github.com/tinyzimmer/gst-pipeline-operator/apis/pipelines/v1"
	pipelinetypes "github.com/tinyzimmer/gst-pipeline-operator/pkg/types"
	"github.com/tinyzimmer/gst-pipeline-operator/pkg/util"
)

var log = ctrl.Log.WithName("pipeline-manager")

var backoffLimit int32 = 5

// A globally held map of the current pipeline managers running
var managers = make(map[types.UID]*PipelineManager)
var managersMutex sync.Mutex

// GetManagerForPipeline returns a PipelineManager for the given transformation pipeline.
// If one already exists globally, it is returned.
func GetManagerForPipeline(client client.Client, pipeline pipelinetypes.Pipeline) *PipelineManager {
	managersMutex.Lock()
	defer managersMutex.Unlock()

	if manager, ok := managers[pipeline.GetUID()]; ok {
		return manager
	}
	managers[pipeline.GetUID()] = &PipelineManager{
		client:     client,
		pipeline:   pipeline,
		reloadChan: make(chan struct{}),
		stopChan:   make(chan struct{}),
	}
	return managers[pipeline.GetUID()]
}

// PipelineManager is an object for watching MinIO buckets for changes and queuing
// processing in a pipeline. It exports a method for reloading configuration changes.
type PipelineManager struct {
	client     client.Client
	pipeline   pipelinetypes.Pipeline
	reloadChan chan struct{}
	stopChan   chan struct{}
	running    bool
	mux        sync.Mutex
}

var marker = ".gst-watch"

// Start starts the pipeline manager.
func (p *PipelineManager) Start() error {
	p.mux.Lock()
	defer p.mux.Unlock()

	if p.running {
		return errors.New("pipeline manager is already running")
	}

	srcConfigFull := p.pipeline.GetSrcConfig()
	if srcConfigFull.MinIO == nil {
		return errors.New("Non-MinIO sources are not yet implemented")
	}
	srcConfig := srcConfigFull.MinIO

	client, err := util.GetMinIOClient(srcConfig, util.MinIOWatchCredentialsFromCR(p.client, p.pipeline))
	if err != nil {
		return err
	}

	markerName := path.Join(srcConfig.GetPrefix(), marker)

	// Check for a marker in the prefix we are watching. This checks for the existence of the
	// bucket as well as ensure the subsequent watch works correctly.
	obj, err := client.GetObject(context.TODO(), srcConfig.GetBucket(), markerName, minio.GetObjectOptions{})
	if err != nil {
		return err
	}

	// The errors we care about would be while trying to read it
	if _, err := ioutil.ReadAll(obj); err != nil {
		if resErr, ok := err.(minio.ErrorResponse); ok {
			switch resErr.Code {
			case "NoSuchKey":
				log.Info("Laying watch marker in bucket prefix", "Bucket", srcConfig.GetBucket(), "Prefix", srcConfig.GetPrefix())
				if _, err := client.PutObject(context.TODO(), srcConfig.GetBucket(), markerName, bytes.NewReader([]byte{}), 0, minio.PutObjectOptions{}); err != nil {
					return err
				}
			default:
				return err
			}
		} else {
			return err
		}
	}

	go p.watchSrcBucket(srcConfig, client)
	p.running = true
	return nil
}

// IsRunning returns true if the pipeline manager is already running.
func (p *PipelineManager) IsRunning() bool { return p.running }

// Reload reloads the bucket watchers with the given pipeline configuration.
func (p *PipelineManager) Reload(cfg pipelinetypes.Pipeline) {
	p.mux.Lock()
	defer p.mux.Unlock()

	p.pipeline = cfg
	p.reloadChan <- struct{}{}
}

// Stop stops the bucket watching goroutine.
func (p *PipelineManager) Stop() {
	p.mux.Lock()
	defer p.mux.Unlock()

	p.stopChan <- struct{}{}
	p.running = false
}

func (p *PipelineManager) watchSrcBucket(srcConfig *pipelinesmeta.MinIOConfig, client *minio.Client) {
	log.Info("Watching for object created events", "Bucket", srcConfig.GetBucket(), "Prefix", srcConfig.GetPrefix())
	eventChan := client.ListenBucketNotification(context.Background(), srcConfig.GetBucket(), srcConfig.GetPrefix(), "", []string{"s3:ObjectCreated:*"})
	excludeRegex := srcConfig.GetExcludeRegex()
	for {
		select {
		case event := <-eventChan:
			for _, record := range event.Records {
				log.Info("Processing record from MinIO event", "Record", record)
				if excludeRegex != nil && excludeRegex.MatchString(record.S3.Object.Key) {
					log.Info("Skipping processing for item matching exclude regex", "Object", record.S3.Object.Key)
					continue
				}
				p.createJob(srcConfig, record.S3.Object.Key)
			}
		case <-p.reloadChan:
			srcConfig = p.pipeline.GetSrcConfig().MinIO // TODO
			excludeRegex = srcConfig.GetExcludeRegex()
			log.Info("Reloading event channel", "Bucket", srcConfig.GetBucket(), "Prefix", srcConfig.GetPrefix())
			eventChan = client.ListenBucketNotification(context.Background(), srcConfig.GetBucket(), srcConfig.GetPrefix(), "", []string{"s3:ObjectCreated:*"})
		case <-p.stopChan:
			return
		}
	}
}

func (p *PipelineManager) createJob(srcConfig *pipelinesmeta.MinIOConfig, object string) {
	log.Info("Creating pipeline job", "Bucket", srcConfig.GetBucket(), "Key", object)
	job := p.newJobForObject(object)
	if err := p.client.Create(context.TODO(), job); err != nil {
		log.Error(err, "Failed to create processing job for object")
	}
}

func (p *PipelineManager) newJobForObject(key string) *pipelinesv1.Job {
	job := &pipelinesv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName:    p.pipeline.GetName(),
			Namespace:       p.pipeline.GetNamespace(),
			Labels:          pipelinesv1.GetJobLabels(p.pipeline, key),
			OwnerReferences: p.pipeline.OwnerReferences(),
		},
		Spec: pipelinesv1.JobSpec{
			PipelineReference: pipelinesmeta.PipelineReference{
				Name: p.pipeline.GetName(),
				Kind: p.pipeline.GetPipelineKind(),
			},
			Source: &pipelinesmeta.Object{
				Name:   key,
				Config: p.pipeline.GetSrcConfig(),
			},
			Sinks: p.pipeline.GetSinkObjects(key),
		},
	}
	return job
}
