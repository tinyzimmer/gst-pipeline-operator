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
func (t *Transform) OwnerReferences() []metav1.OwnerReference { return ownerReferences(t) }

// GetPipelineKind satisfies the Pipeline interface and returns the type of the pipeline.
func (t *Transform) GetPipelineKind() pipelinesmeta.PipelineKind {
	return PipelineTransform
}

// GetPipelineConfig returns the PipelineConfig.
func (t *Transform) GetPipelineConfig() *pipelinesmeta.PipelineConfig { return t.Spec.Pipeline }

// GetSrcConfig will return the src config for this pipeline merged with the globals.
func (t *Transform) GetSrcConfig() *pipelinesmeta.SourceSinkConfig {
	return mergeConfigs(t.Spec.Globals, t.Spec.Src)
}

// GetSinkConfig will return the sink config for this pipeline merged with the globals.
func (t *Transform) GetSinkConfig() *pipelinesmeta.SourceSinkConfig {
	return mergeConfigs(t.Spec.Globals, t.Spec.Sink)
}

// GetSinkObjects returns the sink objects for a pipeline.
func (t *Transform) GetSinkObjects(srcKey string) []*pipelinesmeta.Object {
	return []*pipelinesmeta.Object{
		{
			Name:       t.GetSinkConfig().MinIO.GetDestinationKey(srcKey), // TODO
			Config:     t.GetSinkConfig(),
			StreamType: pipelinesmeta.StreamTypeAll,
		},
	}
}
