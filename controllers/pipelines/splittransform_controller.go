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

// SplitTransformReconciler reconciles a SplitTransform object
type SplitTransformReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=pipelines.gst.io,resources=splittransforms,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=pipelines.gst.io,resources=splittransforms/status,verbs=get;update;patch

// Reconcile reconciles a splittransform pipeline.
func (r *SplitTransformReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := r.Log.WithValues("splittransform", req.NamespacedName)

	// Fetch the object for the request
	pipeline := &pipelinesv1.SplitTransform{}
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

// SetupWithManager adds the SplitTransformReconciler to the given manager.
func (r *SplitTransformReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&pipelinesv1.SplitTransform{}).
		Complete(r)
}

func (r *SplitTransformReconciler) generationObserved(pipeline *pipelinesv1.SplitTransform) bool {
	for _, cond := range pipeline.Status.Conditions {
		if cond.ObservedGeneration == pipeline.GetGeneration() {
			return true
		}
	}
	return false
}

func (r *SplitTransformReconciler) removeFinalizers(ctx context.Context, reqLogger logr.Logger, pipeline *pipelinesv1.SplitTransform) error {
	pipeline.SetFinalizers([]string{})
	return r.Client.Update(ctx, pipeline)
}

func (r *SplitTransformReconciler) ensureFinalizers(ctx context.Context, reqLogger logr.Logger, pipeline *pipelinesv1.SplitTransform) error {
	finalizers := pipeline.GetFinalizers()
	if !contains(finalizers, pipelineFinalizer) {
		reqLogger.Info("Setting finalizers on pipeline")
		pipeline.SetFinalizers([]string{pipelineFinalizer})
		return r.Client.Update(ctx, pipeline)
	}
	return nil
}
