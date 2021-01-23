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
	"encoding/json"

	pipelinesmeta "github.com/tinyzimmer/gst-pipeline-operator/apis/meta/v1"
	pipelinesv1 "github.com/tinyzimmer/gst-pipeline-operator/apis/pipelines/v1"
	pipelinetypes "github.com/tinyzimmer/gst-pipeline-operator/pkg/types"

	"github.com/go-logr/logr"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var backoffLimit int32 = 5

func reconcileJob(ctx context.Context, reqLogger logr.Logger, c client.Client, pipelineJob *pipelinesv1.Job, job *batchv1.Job) error {
	nn := types.NamespacedName{
		Name:      job.GetName(),
		Namespace: job.GetNamespace(),
	}

	// Check if job exists
	found := &batchv1.Job{}
	err := c.Get(ctx, nn, found)
	if err != nil {
		if client.IgnoreNotFound(err) != nil {
			return err
		}
		// Need to create job
		reqLogger.Info("Creating new Job", "Name", job.GetName(), "Namespace", job.GetNamespace())
		err := c.Create(ctx, job)
		if err != nil {
			return err
		}
		// Add job pending status condition
		reqLogger.Info("Updating pipeline job status to Pending")
		pipelineJob.Status.Conditions = append(pipelineJob.Status.Conditions, metav1.Condition{
			Type:               string(pipelinesv1.JobPending),
			Status:             metav1.ConditionTrue,
			ObservedGeneration: pipelineJob.GetGeneration(),
			LastTransitionTime: metav1.Now(),
			Reason:             "JobCreated",
			Message:            "The pipeline job has been created",
		})
		return c.Status().Update(ctx, pipelineJob)
	}

	// The job exists

	reqLogger = reqLogger.WithValues("Name", found.GetName(), "Namespace", found.GetNamespace())

	// Check if the job is still pending
	if jobPending(found) {
		reqLogger.Info("Job is currently pending creation")
		// Add job in progress status condition
		if !statusObservedForGeneration(string(pipelinesv1.JobPending), "JobPending", pipelineJob) {
			reqLogger.Info("Job is waiting for active containers")
			pipelineJob.Status.Conditions = append(pipelineJob.Status.Conditions, metav1.Condition{
				Type:               string(pipelinesv1.JobPending),
				Status:             metav1.ConditionTrue,
				ObservedGeneration: pipelineJob.GetGeneration(),
				LastTransitionTime: metav1.Now(),
				Reason:             "JobPending",
				Message:            "Waiting for the job to be scheduled",
			})
			return c.Status().Update(ctx, pipelineJob)
		}
		return nil
	}

	// Check if the job is still in progress
	if jobInProgress(found) {
		reqLogger.Info("Job is currently in progress")
		// Add job in progress status condition
		if !statusObservedForGeneration(string(pipelinesv1.JobInProgress), "JobInProgress", pipelineJob) {
			reqLogger.Info("Job is in progress, updating status")
			pipelineJob.Status.Conditions = append(pipelineJob.Status.Conditions, metav1.Condition{
				Type:               string(pipelinesv1.JobInProgress),
				Status:             metav1.ConditionTrue,
				ObservedGeneration: pipelineJob.GetGeneration(),
				LastTransitionTime: metav1.Now(),
				Reason:             "JobInProgress",
				Message:            "The pipeline job is currently running",
			})
			return c.Status().Update(ctx, pipelineJob)
		}
		return nil
	}

	// Check if the job succeeded
	if jobSucceeded(found) {
		// Add job finished status condition
		if !statusObservedForGeneration(string(pipelinesv1.JobFinished), "JobFinished", pipelineJob) {
			reqLogger.Info("Job finished successfully, updating status")
			pipelineJob.Status.Conditions = append(pipelineJob.Status.Conditions, metav1.Condition{
				Type:               string(pipelinesv1.JobFinished),
				Status:             metav1.ConditionTrue,
				ObservedGeneration: pipelineJob.GetGeneration(),
				LastTransitionTime: metav1.Now(),
				Reason:             "JobFinished",
				Message:            "The pipeline job completed successfully",
			})
			return c.Status().Update(ctx, pipelineJob)
		}
		return nil
	}

	// Check if the job failed
	if jobFailed(found) {
		// Add job failed status condition
		if !statusObservedForGeneration(string(pipelinesv1.JobFailed), "JobFailed", pipelineJob) {
			reqLogger.Info("Job failed, updating status")
			pipelineJob.Status.Conditions = append(pipelineJob.Status.Conditions, metav1.Condition{
				Type:               string(pipelinesv1.JobFailed),
				Status:             metav1.ConditionTrue,
				ObservedGeneration: pipelineJob.GetGeneration(),
				LastTransitionTime: metav1.Now(),
				Reason:             "JobFailed",
				Message:            "The pipeline job failed to complete",
			})
			return c.Status().Update(ctx, pipelineJob)
		}
		return nil
	}

	return nil
}

func jobSucceeded(job *batchv1.Job) bool  { return job.Status.Succeeded == 1 }
func jobFailed(job *batchv1.Job) bool     { return job.Status.Failed == 1 }
func jobInProgress(job *batchv1.Job) bool { return job.Status.Succeeded == 0 && job.Status.Failed == 0 }
func jobPending(job *batchv1.Job) bool    { return job.Status.Active == 0 && jobInProgress(job) }

func newPipelineJob(pipelineJob *pipelinesv1.Job, pipeline pipelinetypes.Pipeline) (*batchv1.Job, error) {
	// TODO
	srcConfig := pipeline.GetSrcConfig().MinIO
	sinkConfig := pipeline.GetSinkConfig().MinIO
	srcSecret, err := srcConfig.GetCredentialsSecret()
	if err != nil {
		return nil, err
	}
	sinkSecret, err := sinkConfig.GetCredentialsSecret()
	if err != nil {
		return nil, err
	}
	pipelineCfg := pipeline.GetPipelineConfig()
	marshaledConfig, err := json.Marshal(pipelineCfg)
	if err != nil {
		return nil, err
	}
	marshaledSrc, err := json.Marshal(pipelineJob.Spec.Source)
	if err != nil {
		return nil, err
	}
	marshaledSinks, err := json.Marshal(pipelineJob.Spec.Sinks)
	if err != nil {
		return nil, err
	}
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:            pipelineJob.GetName(),
			Namespace:       pipelineJob.GetNamespace(),
			Labels:          pipelineJob.GetLabels(),
			OwnerReferences: pipelineJob.OwnerReferences(),
		},
		Spec: batchv1.JobSpec{
			BackoffLimit: &backoffLimit,
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyOnFailure,
					Containers: []corev1.Container{
						{
							Name:      "gstreamer",
							Image:     pipelineCfg.GetImage(),
							Resources: pipelineCfg.Resources,
							Env: []corev1.EnvVar{
								{
									Name:  "GST_DEBUG",
									Value: pipelineCfg.GetGSTDebug(),
								},
								{
									Name:  pipelinesmeta.JobSrcObjectsEnvVar,
									Value: string(marshaledSrc),
								},
								{
									Name:  pipelinesmeta.JobSinkObjectsEnvVar,
									Value: string(marshaledSinks),
								},
								{
									Name:  pipelinesmeta.JobPipelineConfigEnvVar,
									Value: string(marshaledConfig),
								},
								{
									Name: pipelinesmeta.MinIOSrcAccessKeyIDEnvVar,
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: srcSecret,
											},
											Key: pipelinesmeta.AccessKeyIDKey,
										},
									},
								},
								{
									Name: pipelinesmeta.MinIOSrcSecretAccessKeyEnvVar,
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: srcSecret,
											},
											Key: pipelinesmeta.SecretAccessKeyKey,
										},
									},
								},
								{
									Name: pipelinesmeta.MinIOSinkAccessKeyIDEnvVar,
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: sinkSecret,
											},
											Key: pipelinesmeta.AccessKeyIDKey,
										},
									},
								},
								{
									Name: pipelinesmeta.MinIOSinkSecretAccessKeyEnvVar,
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: sinkSecret,
											},
											Key: pipelinesmeta.SecretAccessKeyKey,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return job, nil
}
