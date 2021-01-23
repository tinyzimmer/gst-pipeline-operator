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
