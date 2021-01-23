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

package types

import (
	pipelinesmeta "github.com/tinyzimmer/gst-pipeline-operator/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Pipeline is a generic interface implemented by the different Pipeline types.
type Pipeline interface {
	// Extends the metav1.Object interface
	metav1.Object

	// OwnerReferences should return the owner references that can be used to apply
	// ownership to this pipeline.
	OwnerReferences() []metav1.OwnerReference
	// GetPipelineKind should return the type of the pipeline.
	GetPipelineKind() pipelinesmeta.PipelineKind
	// GetPipelineConfig should return the element configurations for the pipeline.
	GetPipelineConfig() *pipelinesmeta.PipelineConfig
	// GetSrcConfig should return the source configuration for the pipeline.
	GetSrcConfig() *pipelinesmeta.SourceSinkConfig
	// GetSinkConfig should return a sink configuration for the pipeline. This method
	// is primarily used to retrieve any required credentials when constructing a
	// pipeline job.
	GetSinkConfig() *pipelinesmeta.SourceSinkConfig
	// GetSinkObjects should compute the sink objects for a pipeline based on a given
	// source key.
	GetSinkObjects(srcKey string) []*pipelinesmeta.Object
}
