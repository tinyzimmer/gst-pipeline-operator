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
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"

	pipelinesmeta "github.com/tinyzimmer/gst-pipeline-operator/apis/meta/v1"
	"github.com/tinyzimmer/gst-pipeline-operator/pkg/types"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

func mergeConfigs(into, from *pipelinesmeta.SourceSinkConfig) *pipelinesmeta.SourceSinkConfig {
	if into == nil {
		return from
	}
	cfgBody, err := json.Marshal(from)
	if err != nil {
		fmt.Println("Failed to marshal global config to json:", err)
		return from
	}
	intoCopy := into.DeepCopy()
	if err := json.Unmarshal(cfgBody, intoCopy); err != nil {
		fmt.Println("Error unmarshaling configuration on top of globals:", err)
		return from
	}
	return intoCopy
}

// GetJobLabels returns the labels to apply to a new job issued from this pipeline.
func GetJobLabels(pipeline types.Pipeline, key string) map[string]string {
	h := md5.New()
	io.WriteString(h, key)
	labels := map[string]string{
		pipelinesmeta.JobPipelineLabel:     pipeline.GetName(),
		pipelinesmeta.JobPipelineKindLabel: string(pipeline.GetPipelineKind()),
		pipelinesmeta.JobObjectLabel:       fmt.Sprintf("%x", h.Sum(nil)),
	}
	src := pipeline.GetSrcConfig()
	if src != nil && src.MinIO != nil {
		labels[pipelinesmeta.JobBucketLabel] = src.MinIO.GetBucket()
	}
	return labels
}

func ownerReferences(obj runtime.Object) []metav1.OwnerReference {
	return []metav1.OwnerReference{*metav1.NewControllerRef(obj.(metav1.Object), obj.GetObjectKind().GroupVersionKind())}
}
