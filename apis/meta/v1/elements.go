package v1

// ElementConfig represents the configuration of a single element in a transform pipeline.
type ElementConfig struct {
	// The name of the element. See the GStreamer plugin documentation for a comprehensive
	// list of all the plugins available. Custom pipeline images can also be used that are
	// prebaked with additional plugins.
	Name string `json:"name,omitempty"`
	// Applies an alias to this element in the pipeline configuration. This allows you to specify an
	// element block with this value as the name and have it act as a "goto" or "linkto" while building
	// the pipeline. Note that the aliases "video-out" and "audio-out" are reserved for internal use.
	Alias string `json:"alias,omitempty"`
	// The alias to an element to treat as this configuration. Useful for directing the output of elements
	// with multiple src pads, such as decodebin.
	GoTo string `json:"goto,omitempty"`
	// The alias to an element to link the previous element's sink pad to. Useful for directing the branches of
	// a multi-stream pipeline to a muxer. A linkto almost always needs to be followed by a goto, except when
	// the element being linked to is next in the pipeline, in which case you can omit the linkto entirely.
	LinkTo string `json:"linkto,omitempty"`
	// Optional properties to apply to this element. To not piss off the CRD generator values are
	// declared as a string, but almost anything that can be passed to gst-launch-1.0 will work.
	// Caps will be parsed from their string representation.
	Properties map[string]string `json:"properties,omitempty"`
}

// LinkToVideoOut is used during split pipelines to designate the src of a video sink
const LinkToVideoOut = "video-out"

// LinkToAudioOut is used during split pipelines to designate the src of an audio sink
const LinkToAudioOut = "audio-out"
