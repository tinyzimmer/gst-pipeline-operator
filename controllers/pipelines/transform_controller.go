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

	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	pipelinesmeta "github.com/tinyzimmer/gst-pipeline-operator/apis/meta/v1"
	pipelinesv1 "github.com/tinyzimmer/gst-pipeline-operator/apis/pipelines/v1"
	"github.com/tinyzimmer/gst-pipeline-operator/pkg/managers"
)

var pipelineFinalizer = "pipelines.gst.io/finalize"

// TransformReconciler reconciles a Transform object
type TransformReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=coordination.k8s.io,resources=leases,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=pipelines.gst.io,resources=transforms,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=pipelines.gst.io,resources=transforms/status,verbs=get;update;patch

// Reconcile reconciles a Transform pipeline
func (r *TransformReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := r.Log.WithValues("transform", req.NamespacedName)

	// Fetch the object for the request
	pipeline := &pipelinesv1.Transform{}
	err := r.Client.Get(ctx, req.NamespacedName, pipeline)
	if err != nil {
		if client.IgnoreNotFound(err) == nil {
			// Object was deleted
			return ctrl.Result{}, nil
		}
		// Requeue any other error
		return ctrl.Result{}, err
	}

	// Get the controller for this pipeline
	controller := managers.GetManagerForPipeline(r.Client, pipeline)

	// Check if we are running finalizers
	if pipeline.GetDeletionTimestamp() != nil {
		if controller.IsRunning() {
			controller.Stop()
		}
		return ctrl.Result{}, r.removeFinalizers(ctx, reqLogger, pipeline)
	}

	if !controller.IsRunning() {
		reqLogger.Info("Starting PipelineManager")
		if err := controller.Start(); err != nil {
			return ctrl.Result{}, err
		}
	} else {
		reqLogger.Info("PipelineManager is already running, reloading config")
		controller.Reload(pipeline)
	}

	if err := r.ensureFinalizers(ctx, reqLogger, pipeline); err != nil {
		return ctrl.Result{}, nil
	}

	if !r.generationObserved(pipeline) {
		pipeline.Status.Conditions = append(pipeline.Status.Conditions, metav1.Condition{
			Type:               string(pipelinesmeta.PipelineInSync),
			Status:             metav1.ConditionTrue,
			ObservedGeneration: pipeline.GetGeneration(),
			LastTransitionTime: metav1.Now(),
			Reason:             string(pipelinesmeta.PipelineInSync),
			Message:            "The pipeline configuration is in-sync",
		})
		if err := r.Client.Status().Update(ctx, pipeline); err != nil {
			return ctrl.Result{}, err
		}
	}

	reqLogger.Info("Reconcile finished")
	return ctrl.Result{}, nil
}

func (r *TransformReconciler) generationObserved(pipeline *pipelinesv1.Transform) bool {
	for _, cond := range pipeline.Status.Conditions {
		if cond.ObservedGeneration == pipeline.GetGeneration() {
			return true
		}
	}
	return false
}

func (r *TransformReconciler) removeFinalizers(ctx context.Context, reqLogger logr.Logger, pipeline *pipelinesv1.Transform) error {
	pipeline.SetFinalizers([]string{})
	return r.Client.Update(ctx, pipeline)
}

func (r *TransformReconciler) ensureFinalizers(ctx context.Context, reqLogger logr.Logger, pipeline *pipelinesv1.Transform) error {
	finalizers := pipeline.GetFinalizers()
	if !contains(finalizers, pipelineFinalizer) {
		reqLogger.Info("Setting finalizers on pipeline")
		pipeline.SetFinalizers([]string{pipelineFinalizer})
		return r.Client.Update(ctx, pipeline)
	}
	return nil
}

func contains(ss []string, s string) bool {
	for _, x := range ss {
		if x == s {
			return true
		}
	}
	return false
}

// SetupWithManager adds the Transform pipeline controller to the manager.
func (r *TransformReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&pipelinesv1.Transform{}).
		Owns(&pipelinesv1.Job{}).
		Complete(r)
}
