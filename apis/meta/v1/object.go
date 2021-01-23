package v1

// Object represents either a source or destination object for a job.
type Object struct {
	// The actual name for the object being read or written to. In the context of
	// a source object this is pulled from a watch event. In the context of a destination
	// this is computed by the controller from the user supplied configuration.
	Name string `json:"name"`
	// The endpoint and bucket configurations for the object.
	Config *SourceSinkConfig `json:"config"`
	// The type of the stream for this object. Only applies to sinks. For a split transform
	// pipeline there will be an Object for each stream. Otherwise there will be a single
	// object with a StreamTypeAll.
	StreamType StreamType `json:"streamType"`
}

// StreamType represents a type of stream found in a source input, or designated for an output.
type StreamType string

const (
	// StreamTypeAll represents all possible output streams.
	StreamTypeAll StreamType = "all"
	// StreamTypeVideo represents a video stream.
	StreamTypeVideo StreamType = "video"
	// StreamTypeAudio represents an audio stream.
	StreamTypeAudio StreamType = "audio"
)
