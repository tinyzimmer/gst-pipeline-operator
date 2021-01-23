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

import pipelinesv1 "github.com/tinyzimmer/gst-pipeline-operator/apis/pipelines/v1"

func statusObservedForGeneration(status, reason string, job *pipelinesv1.Job) bool {
	for _, cond := range job.Status.Conditions {
		if cond.Type == status && cond.Reason == reason && cond.ObservedGeneration == job.GetGeneration() {
			return true
		}
	}
	return false
}
