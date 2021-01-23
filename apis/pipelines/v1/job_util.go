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
	"context"

	pipelinesmeta "github.com/tinyzimmer/gst-pipeline-operator/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// OwnerReferences returns the OwnerReferences for this job.
func (j *Job) OwnerReferences() []metav1.OwnerReference { return ownerReferences(j) }

// GetPipelineKind returns the type of the pipeline.
func (j *Job) GetPipelineKind() pipelinesmeta.PipelineKind { return j.Spec.PipelineReference.Kind }

// GetTransformPipeline returns the transform pipeline for this job spec.
func (j *Job) GetTransformPipeline(ctx context.Context, client client.Client) (*Transform, error) {
	nn := types.NamespacedName{
		Name:      j.Spec.PipelineReference.Name,
		Namespace: j.GetNamespace(),
	}
	var pipeline Transform
	return &pipeline, client.Get(ctx, nn, &pipeline)
}

// GetSplitTransformPipeline returns the splittransform pipeline for this job spec.
func (j *Job) GetSplitTransformPipeline(ctx context.Context, client client.Client) (*SplitTransform, error) {
	nn := types.NamespacedName{
		Name:      j.Spec.PipelineReference.Name,
		Namespace: j.GetNamespace(),
	}
	var pipeline SplitTransform
	return &pipeline, client.Get(ctx, nn, &pipeline)
}
