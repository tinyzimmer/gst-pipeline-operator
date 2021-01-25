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

// Default Values
const (
	// DefaultRegion if none is provided for any configurations
	DefaultRegion = "us-east-1"
	// AccessKeyIDKey is the key in secrets where the access key ID is stored.
	AccessKeyIDKey = "access-key-id"
	// SecretAccessKeyKey is the key in secrets where the secret access key is stored.
	SecretAccessKeyKey = "secret-access-key"
	// DefaultGSTDebug is the default GST_DEBUG value set in the environment for pipelines.
	DefaultGSTDebug = "4"
	// DefaultDotInterval is the default interval to query a pipeline for graphs.
	DefaultDotInterval = 3
)

// Annotations
const (
	JobCreationSpecAnnotation = "pipelines.gst.io/creation-spec"
)

// Labels
const (
	// JobPipelineLabel is the label on a job to denote the Transform pipeline that initiated
	// it.
	JobPipelineLabel = "pipelines.gst.io/pipeline"
	// JobPipelineKindLabel is the label where the type of the pipeline is stored.
	JobPipelineKindLabel = "pipelines.gst.io/kind"
	// JobObjectLabel is the label on a job to denote the object key it is processing.
	JobObjectLabel = "pipelines.gst.io/object"
	// JobBucketLabel is the label on a job to denote the bucket where the object is that
	// is being processed.
	JobBucketLabel = "pipelines.gst.io/bucket"
)

// Environment Variables
const (
	// The environment variable where the access key id for the src bucket is stored.
	MinIOSrcAccessKeyIDEnvVar = "MINIO_SRC_ACCESS_KEY_ID"
	// The environment variable where the secret access key for the src bucket is stored.
	MinIOSrcSecretAccessKeyEnvVar = "MINIO_SRC_SECRET_ACCESS_KEY"
	// The environment variable where the access key id for the sink bucket is stored.
	MinIOSinkAccessKeyIDEnvVar = "MINIO_SINK_ACCESS_KEY_ID"
	// The environment variable where the secret access key for the sink bucket is stored.
	MinIOSinkSecretAccessKeyEnvVar = "MINIO_SINK_SECRET_ACCESS_KEY"
	// The environment variable where the pipeline config is serialized and set.
	JobPipelineConfigEnvVar = "GST_PIPELINE_CONFIG"
	// The environment variable where the source object is serialized and set.
	JobSrcObjectsEnvVar = "GST_PIPELINE_SRC_OBJECT"
	// The environment variable where the sink objects are serialized and set.
	JobSinkObjectsEnvVar = "GST_PIPELINE_SINK_OBJECTS"
	// The environment variable where the name of the pipeline being watched is set for watcher
	// processes.
	WatcherPipelineNameEnvVar = "GST_WATCH_PIPELINE_NAME"
	// The environment variable where the pipeline kind is set for the watcher processes.
	WatcherPipelineKindEnvVar = "GST_WATCH_PIPELINE_KIND"
)
