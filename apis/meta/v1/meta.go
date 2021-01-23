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

// PipelineKind represents a type of pipeline.
type PipelineKind string

// PipelineReference is used to refer to the pipeline that holds the configuration
// for a given job.
type PipelineReference struct {
	// Name is the name of the Pipeline CR
	Name string `json:"name"`
	// Kind is the type of the Pipeline CR
	Kind PipelineKind `json:"kind"`
}

// PipelineState represents the state of a Pipeline CR.
type PipelineState string

const (
	// PipelineInSync represents that the pipeline configuration is in sync with the watchers.
	PipelineInSync PipelineState = "InSync"
)
