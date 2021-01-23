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
