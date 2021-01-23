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

const (
	// PipelineTransform represents a transform pipeline
	PipelineTransform pipelinesmeta.PipelineKind = "Transform"
)

// TransformSpec defines the desired state of Transform
type TransformSpec struct {
	// Global configurations to apply when omitted from the src or sink configurations.
	Globals *pipelinesmeta.SourceSinkConfig `json:"globals,omitempty"`
	// Configurations for src object to the pipeline.
	Src *pipelinesmeta.SourceSinkConfig `json:"src"`
	// Configurations for sink objects from the pipeline.
	Sink *pipelinesmeta.SourceSinkConfig `json:"sink"`
	// The configuration for the processing pipeline
	Pipeline *pipelinesmeta.PipelineConfig `json:"pipeline"`
}

// TransformStatus defines the observed state of Transform
type TransformStatus struct {
	// Conditions represent the latest available observations of a transform's state
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=`.metadata.creationTimestamp`
// +kubebuilder:printcolumn:name="Status",type="string",priority=1,JSONPath=`.status.conditions[-1].message`

// Transform is the Schema for the transforms API
type Transform struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TransformSpec   `json:"spec,omitempty"`
	Status TransformStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TransformList contains a list of Transform
type TransformList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Transform `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Transform{}, &TransformList{})
}
