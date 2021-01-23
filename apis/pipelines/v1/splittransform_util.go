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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	pipelinesmeta "github.com/tinyzimmer/gst-pipeline-operator/apis/meta/v1"
)

// OwnerReferences returns the OwnerReferences for this pipeline to be placed on jobs.
func (t *SplitTransform) OwnerReferences() []metav1.OwnerReference { return ownerReferences(t) }

// GetPipelineKind satisfies the Pipeline interface and returns the type of the pipeline.
func (t *SplitTransform) GetPipelineKind() pipelinesmeta.PipelineKind {
	return PipelineSplitTransform
}

// GetPipelineConfig returns the PipelineConfig.
func (t *SplitTransform) GetPipelineConfig() *pipelinesmeta.PipelineConfig { return t.Spec.Pipeline }

// GetSrcConfig will return the src config for this pipeline merged with the globals.
func (t *SplitTransform) GetSrcConfig() *pipelinesmeta.SourceSinkConfig {
	return mergeConfigs(t.Spec.Globals, t.Spec.Src)
}

// GetSinkConfig will return the first non-dropped sink config found in this pipeline.
func (t *SplitTransform) GetSinkConfig() *pipelinesmeta.SourceSinkConfig {
	if t.Spec.Video != nil {
		return mergeConfigs(t.Spec.Globals, t.Spec.Video)
	}
	if t.Spec.Audio != nil {
		return mergeConfigs(t.Spec.Globals, t.Spec.Audio)
	}
	return nil
}

// GetSinkObjects returns the sink objects for a pipeline.
func (t *SplitTransform) GetSinkObjects(srcKey string) []*pipelinesmeta.Object {
	objs := make([]*pipelinesmeta.Object, 0)
	if t.Spec.Video != nil {
		videoCfg := mergeConfigs(t.Spec.Globals, t.Spec.Video)
		objs = append(objs, &pipelinesmeta.Object{
			Name:       videoCfg.MinIO.GetDestinationKey(srcKey), // TODO
			Config:     videoCfg,
			StreamType: pipelinesmeta.StreamTypeVideo,
		})
	}
	if t.Spec.Audio != nil {
		audioCfg := mergeConfigs(t.Spec.Globals, t.Spec.Audio)
		objs = append(objs, &pipelinesmeta.Object{
			Name:       audioCfg.MinIO.GetDestinationKey(srcKey), // TODO
			Config:     audioCfg,
			StreamType: pipelinesmeta.StreamTypeAudio,
		})
	}
	return objs
}
