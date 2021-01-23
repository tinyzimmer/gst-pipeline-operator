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

package pipelines

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	pipelinesv1 "github.com/tinyzimmer/gst-pipeline-operator/apis/pipelines/v1"
	pipelinetypes "github.com/tinyzimmer/gst-pipeline-operator/pkg/types"
)

// JobReconciler reconciles a Job object
type JobReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups=pipelines.gst.io,resources=jobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=pipelines.gst.io,resources=jobs/status,verbs=get;update;patch

// Reconcile reconciles a Job
func (r *JobReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := r.Log.WithValues("job", req.NamespacedName)

	job := &pipelinesv1.Job{}
	err := r.Client.Get(ctx, req.NamespacedName, job)
	if err != nil {
		if client.IgnoreNotFound(err) == nil {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	var pipeline pipelinetypes.Pipeline
	switch job.GetPipelineKind() {
	case pipelinesv1.PipelineTransform:
		reqLogger.Info("Fetching Transform pipeline from Job")
		pipeline, err = job.GetTransformPipeline(ctx, r.Client)
	case pipelinesv1.PipelineSplitTransform:
		pipeline, err = job.GetSplitTransformPipeline(ctx, r.Client)
	default:
		err = fmt.Errorf("Unknown pipeline kind: %s", string(job.GetPipelineKind()))
	}
	if err != nil {
		return ctrl.Result{}, err
	}

	batchjob, err := newPipelineJob(job, pipeline)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, reconcileJob(ctx, reqLogger, r.Client, job, batchjob)
}

// SetupWithManager adds the Job reconciler to the given manager.
func (r *JobReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&pipelinesv1.Job{}).
		Owns(&batchv1.Job{}).
		Complete(r)
}
