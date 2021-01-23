package v1

// SourceSinkConfig is used to declare configurations related to the retrieval or
// saving of pipeline objects.
type SourceSinkConfig struct {
	// Configurations for a MinIO source or sink
	MinIO *MinIOConfig `json:"minio,omitempty"`
}
