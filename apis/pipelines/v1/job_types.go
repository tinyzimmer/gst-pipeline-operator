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
	pipelinesmeta "github.com/tinyzimmer/gst-pipeline-operator/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// JobState represents the state of a pipeline job.
type JobState string

const (
	// JobPending means the pipeline is waiting to be started.
	JobPending JobState = "Pending"
	// JobInProgress means the pipeline is currently running.
	JobInProgress JobState = "InProgress"
	// JobFinished means the pipeline completed without error.
	JobFinished JobState = "Finished"
	// JobFailed means the pipeline completed with an error.
	JobFailed JobState = "Failed"
)

// JobSpec defines the desired state of Job
type JobSpec struct {
	// A reference to the pipeline for this job's configuration.
	PipelineReference pipelinesmeta.PipelineReference `json:"pipelineRef"`
	// The source object for the pipeline.
	Source *pipelinesmeta.Object `json:"src"`
	// The output objects for the pipeline.
	Sinks []*pipelinesmeta.Object `json:"sinks"`
}

// JobStatus defines the observed state of Job
type JobStatus struct {
	// Conditions represent the latest available observations of a job's state
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=`.metadata.creationTimestamp`
// +kubebuilder:printcolumn:name="Src",type="string",JSONPath=`.spec.src.name`
// +kubebuilder:printcolumn:name="Sinks",type="string",JSONPath=`.spec.sinks[*].name`
// +kubebuilder:printcolumn:name="Status",type="string",priority=1,JSONPath=`.status.conditions[-1].message`

// Job is the Schema for the jobs API
type Job struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JobSpec   `json:"spec,omitempty"`
	Status JobStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// JobList contains a list of Job
type JobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Job `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Job{}, &JobList{})
}
