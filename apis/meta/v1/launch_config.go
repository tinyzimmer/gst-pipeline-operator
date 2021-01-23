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

// GstLaunchConfig is a slice of ElementConfigs that contain internal fields used for
// dynamic linking.
type GstLaunchConfig []*GstElementConfig

// GetElements returns the elements for this pipeline.
func (p *PipelineConfig) GetElements() GstLaunchConfig {
	out := make(GstLaunchConfig, len(p.Elements))
	for idx, elem := range p.Elements {
		out[idx] = &GstElementConfig{ElementConfig: elem}
	}
	return out
}

// GetByAlias returns the configuration for the element at the given alias
func (g GstLaunchConfig) GetByAlias(alias string) *GstElementConfig {
	for _, elem := range g {
		if elem.Alias == alias {
			return elem
		}
	}
	return nil
}

// GstElementConfig is an extension of the ElementConfig struct providing
// private fields for internal tracking while building a dynamic pipeline.
type GstElementConfig struct {
	*ElementConfig
	pipelineName string
	peers        []*GstElementConfig
}

// SetPipelineName sets the name that was assigned to this element by the pipeline for
// later reference.
func (e *GstElementConfig) SetPipelineName(name string) { e.pipelineName = name }

// GetPipelineName returns the name that was assigned to this element by the pipeline.
func (e *GstElementConfig) GetPipelineName() string { return e.pipelineName }

// AddPeer will add a peer to this configuration. It is used for determining which
// sink pads to pair with dynamically added src pads.
func (e *GstElementConfig) AddPeer(peer *GstElementConfig) {
	if e.peers == nil {
		e.peers = make([]*GstElementConfig, 0)
	}
	e.peers = append(e.peers, peer)
}

// GetPeers returns the peers registered for this element.
func (e *GstElementConfig) GetPeers() []*GstElementConfig { return e.peers }
