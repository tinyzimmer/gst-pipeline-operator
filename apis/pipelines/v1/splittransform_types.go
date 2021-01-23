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
	// PipelineSplitTransform represents a splittransform pipeline
	PipelineSplitTransform pipelinesmeta.PipelineKind = "SplitTransform"
)

// SplitTransformSpec defines the desired state of SplitTransform. Note that due to current
// implementation, the various streams can be directed to different buckets, but they have to
// be buckets accessible via the same MinIO/S3 server(s).
type SplitTransformSpec struct {
	// Global configurations to apply when omitted from the src or sink configurations.
	Globals *pipelinesmeta.SourceSinkConfig `json:"globals,omitempty"`
	// Configurations for src object to the pipeline.
	Src *pipelinesmeta.SourceSinkConfig `json:"src"`
	// Configurations for video stream outputs. The linkto field in the pipeline config
	// should be present with the value `video-out` to direct an element to this output.
	Video *pipelinesmeta.SourceSinkConfig `json:"video,omitempty"`
	// Configurations for audio stream outputs. The linkto field in the pipeline config
	// should be present with the value `audio-out` to direct an element to this output.
	Audio *pipelinesmeta.SourceSinkConfig `json:"audio,omitempty"`
	// The configuration for the processing pipeline
	Pipeline *pipelinesmeta.PipelineConfig `json:"pipeline"`
}

// SplitTransformStatus defines the observed state of SplitTransform
type SplitTransformStatus struct {
	// Conditions represent the latest available observations of a splittransform's state
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=`.metadata.creationTimestamp`
// +kubebuilder:printcolumn:name="Status",type="string",priority=1,JSONPath=`.status.conditions[-1].message`

// SplitTransform is the Schema for the splittransforms API
type SplitTransform struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SplitTransformSpec   `json:"spec,omitempty"`
	Status SplitTransformStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SplitTransformList contains a list of SplitTransform
type SplitTransformList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SplitTransform `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SplitTransform{}, &SplitTransformList{})
}
