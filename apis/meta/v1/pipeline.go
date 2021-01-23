package v1

import (
	"fmt"
	"path"
	"strconv"
	"time"

	"github.com/tinyzimmer/gst-pipeline-operator/pkg/version"
	corev1 "k8s.io/api/core/v1"
)

// PipelineConfig represents a series of elements through which to pass the contents of
// processed objects.
type PipelineConfig struct {
	// The image to use to run a/v processing pipelines.
	Image string `json:"image,omitempty"`
	// Debug configurations for the pipeline
	Debug *DebugConfig `json:"debug,omitempty"`
	// A list of element configurations in the order they will be used in the pipeline.
	// Using these is mutually exclusive with a decodebin configuration. This only really
	// works for linear pipelines. That is to say, not the syntax used by `gst-launch-1.0` that
	// allows naming elements and referencing them later in the pipeline. For complex handling
	// of multiple streams decodebin will still be better to work with for now, despite its
	// shortcomings.
	Elements []*ElementConfig `json:"elements,omitempty"`
	// Resource restraints to place on jobs created for this pipeline.
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
}

// DebugConfig represents debug configurations for a GStreamer pipeline.
type DebugConfig struct {
	// The level of log output to produce from the gstreamer process. This value gets set to
	// the GST_DEBUG variable. Defaults to INFO level (4). Higher numbers mean more output.
	LogLevel int `json:"logLevel,omitempty"`
	// Dot specifies to dump a dot file of the pipeline layout for debugging.
	Dot *DotConfig `json:"dot,omitempty"`
}

// DotConfig represents a configuration for the dot output of a pipeline.
type DotConfig struct {
	// The path to save files. The configuration other than the path is assumed to be that of
	// the source of the pipeline. For example, for a MinIO source, this should be a prefix in
	// the same bucket as the source (but not overlapping with the watch prefix otherwise an infinite
	// loop will happen). The files will be saved in directories matching the source object's name with
	// the _debug suffix.
	Path string `json:"path,omitempty"`
	// Specify to also render the pipeline graph to images in the given format. Accepted formats are
	// png, svg, or jpg.
	Render string `json:"render,omitempty"`
	// Whether to save timestamped versions of the pipeline layout. This will produce a new graph for every
	// interval specified by Interval. The default is to only keep the latest graph.
	Timestamped bool `json:"timestamped,omitempty"`
	// The interval in seconds to save pipeline graphs. Defaults to every 3 seconds.
	Interval int `json:"interval,omitempty"`
}

// GetGSTDebug returns the string value of the level to set to GST_DEBUG.
func (p *PipelineConfig) GetGSTDebug() string {
	if p.Debug == nil || p.Debug.LogLevel == 0 {
		return DefaultGSTDebug
	}
	return strconv.Itoa(p.Debug.LogLevel)
}

// DoDotDump returns true if the pipeline has DOT debugging enabled.
func (p *PipelineConfig) DoDotDump() bool {
	return p.Debug != nil && p.Debug.Dot != nil
}

// TimestampDotGraphs returns true if timestamped dot images should be saved.
func (p *PipelineConfig) TimestampDotGraphs() bool {
	if p.Debug == nil || p.Debug.Dot == nil {
		return false
	}
	return p.Debug.Dot.Timestamped
}

// GetDotInterval returns the interval in seconds to query for pipeline graphs.
func (p *PipelineConfig) GetDotInterval() time.Duration {
	if p.Debug == nil || p.Debug.Dot == nil || p.Debug.Dot.Interval == 0 {
		return time.Duration(DefaultDotInterval) * time.Second
	}
	return time.Duration(p.Debug.Dot.Interval) * time.Second
}

// GetDotPath returns the path to save dot graphs based on the given source key.
func (p *PipelineConfig) GetDotPath(srcKey string) string {
	if p.Debug == nil || p.Debug.Dot == nil {
		return ""
	}
	return path.Join(p.Debug.Dot.Path, fmt.Sprintf("%s_debug", path.Base(srcKey)))
}

// GetDotRenderFormat returns the image format that the dot graphs should be encoded to
// when uploading alongside the raw format.
func (p *PipelineConfig) GetDotRenderFormat() string {
	if p.Debug == nil || p.Debug.Dot == nil {
		return ""
	}
	return p.Debug.Dot.Render
}

// GetImage returns the container image to use for the gstreamer pipelines.
func (p *PipelineConfig) GetImage() string {
	if p.Image != "" {
		return p.Image
	}
	return fmt.Sprintf("ghcr.io/tinyzimmer/gst-pipeline-operator/gstreamer:%s", version.Version)
}
