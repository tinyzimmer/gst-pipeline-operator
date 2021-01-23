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

package main

import (
	"errors"
	"fmt"

	"github.com/tinyzimmer/go-gst/gst"
	pipelinesmeta "github.com/tinyzimmer/gst-pipeline-operator/apis/meta/v1"
)

func buildPipelineFromCR(cfg *pipelinesmeta.PipelineConfig, srcObject *pipelinesmeta.Object, sinkObjects []*pipelinesmeta.Object) (*gst.Pipeline, error) {
	// Create a new pipeline
	pipeline, err := gst.NewPipeline("")
	if err != nil {
		return nil, err
	}

	pipelineCfg := cfg.GetElements()

	// Create the source element
	src, err := makeSrcElement(srcObject)
	if err != nil {
		return nil, err
	}

	pipeline.Add(src)

	var last *gst.Element = src
	var lastCfg *pipelinesmeta.GstElementConfig
	var staticSinks bool

	for _, elementCfg := range pipelineCfg {

		// If we are jumping in the pipeline - set the last pointers to the appropriate element and config
		if elementCfg.GoTo != "" {
			// We are jumping in the pipeline
			lastCfg = pipelineCfg.GetByAlias(elementCfg.GoTo)
			if lastCfg == nil {
				return nil, fmt.Errorf("No configuration referenced by alias %s", elementCfg.GoTo)
			}
			// Set the last element to this one
			last, err = elementForPipeline(pipeline, lastCfg)
			if err != nil {
				return nil, err
			}
			continue
		}

		// If we are linking the previous element - perform the links depending on the alias.
		// Sets the last pointers as well, but at this point the user is probably doing a goto
		// next or this is the end of the pipeline.
		if elementCfg.LinkTo != "" {
			var thisElem *gst.Element
			var thisCfg *pipelinesmeta.GstElementConfig
			var err error
			if elementCfg.LinkTo == pipelinesmeta.LinkToVideoOut {
				// Check if this is a split pipeline and we are creating a sink for video
				staticSinks = true
				sinkobj := objectByStreamType(pipelinesmeta.StreamTypeVideo, sinkObjects)
				if sinkobj == nil {
					return nil, errors.New("No video sink configured for pipeline")
				}
				thisElem, thisCfg, err = makeSinkElement(sinkobj)
				if err != nil {
					return nil, err
				}
				pipeline.Add(thisElem)
			} else if elementCfg.LinkTo == pipelinesmeta.LinkToAudioOut {
				// Check if this is a split pipeline and we are creating a sink for audio
				staticSinks = true
				sinkobj := objectByStreamType(pipelinesmeta.StreamTypeAudio, sinkObjects)
				if sinkobj == nil {
					return nil, errors.New("No audio sink configured for pipeline")
				}
				thisElem, thisCfg, err = makeSinkElement(sinkobj)
				if err != nil {
					return nil, err
				}
				pipeline.Add(thisElem)
			} else {
				thisCfg = pipelineCfg.GetByAlias(elementCfg.LinkTo)
				thisElem, err = elementForPipeline(pipeline, thisCfg)
				if err != nil {
					return nil, err
				}
			}
			if err := linkLast(pipeline, last, lastCfg, thisElem, thisCfg); err != nil {
				return nil, err
			}

			last = thisElem
			lastCfg = thisCfg
			continue
		}

		// Neither of the conditions apply we are creating a new element and linking
		// the previous one to it.
		element, err := elementForPipeline(pipeline, elementCfg)
		if err != nil {
			return nil, err
		}

		if err := linkLast(pipeline, last, lastCfg, element, elementCfg); err != nil {
			return nil, err
		}

		last = element
		lastCfg = elementCfg
	}

	// If we did not add static sinks while building the pipeline (i.e. this is a regular Transform pipeline)
	// then create a link to the assumed only sink object
	if !staticSinks {
		sinkobj := objectByStreamType(pipelinesmeta.StreamTypeAll, sinkObjects)
		if sinkobj == nil {
			return nil, errors.New("No sink configured for pipeline")
		}
		sink, sinkCfg, err := makeSinkElement(sinkobj)
		if err != nil {
			return nil, err
		}
		pipeline.Add(sink)
		if err := linkLast(pipeline, last, lastCfg, sink, sinkCfg); err != nil {
			return nil, err
		}
	}

	return pipeline, nil
}

func linkLast(pipeline *gst.Pipeline, last *gst.Element, lastCfg *pipelinesmeta.GstElementConfig, element *gst.Element, elementCfg *pipelinesmeta.GstElementConfig) error {
	// If the last element has a static src pad, link it to this element
	// and continue
	if srcpad := last.GetStaticPad("src"); srcpad != nil {
		return last.Link(element)
	}

	// The last element provides dynamic src pads (we hope - user will find out quick if they messed up)
	lastCfg.AddPeer(elementCfg)
	// weakLastCfg := lastCfg
	last.Connect("no-more-pads", func(self *gst.Element) {
		pads, err := self.GetPads()
		if err != nil {
			self.ErrorMessage(gst.DomainLibrary, gst.LibraryErrorFailed, err.Error(), "")
			return
		}
	Pads:
		for _, srcpad := range pads {
			// Skip already linked and non-src pads
			if srcpad.IsLinked() || srcpad.Direction() != gst.PadDirectionSource {
				continue Pads
			}
			for _, peer := range lastCfg.GetPeers() {
				peerElem, err := elementForPipeline(pipeline, peer)
				if err != nil {
					self.ErrorMessage(gst.DomainLibrary, gst.LibraryErrorFailed, err.Error(), "")
					return
				}
				peersink := peerElem.GetStaticPad("sink")
				if peersink == nil {
					self.ErrorMessage(gst.DomainLibrary, gst.LibraryErrorFailed, fmt.Sprintf("peer %s does not have a static sink pad", peer.Name), "")
					return
				}
				if srcpad.CanLink(peersink) {
					srcpad.Link(peersink)
					continue Pads
				}
			}
		}
	})

	return nil
}

func elementForPipeline(pipeline *gst.Pipeline, cfg *pipelinesmeta.GstElementConfig) (thiselem *gst.Element, err error) {
	// Ensure the element is added to the pipeline
	if name := cfg.GetPipelineName(); name != "" {
		// the element was already created because it was referenced elsewhere
		thiselem, err = pipeline.GetElementByName(name)
		if err != nil {
			return
		}
	} else {
		thiselem, err = makeElement(cfg)
		if err != nil {
			return
		}
		pipeline.Add(thiselem)
		cfg.SetPipelineName(thiselem.GetName())
	}
	return
}

func objectByStreamType(t pipelinesmeta.StreamType, objs []*pipelinesmeta.Object) *pipelinesmeta.Object {
	for _, o := range objs {
		if o.StreamType == t {
			return o
		}
	}
	return nil
}
