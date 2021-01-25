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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/goccy/go-graphviz"
	"github.com/minio/minio-go/v7"
	"github.com/tinyzimmer/go-glib/glib"
	"github.com/tinyzimmer/go-gst/gst"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	pipelinesmeta "github.com/tinyzimmer/gst-pipeline-operator/apis/meta/v1"
	"github.com/tinyzimmer/gst-pipeline-operator/pkg/util"
)

var log = zap.New(zap.UseDevMode(true)).WithName("gst-runner")

func init() { gst.Init(nil) }

func main() {
	mainLoop := glib.NewMainLoop(glib.MainContextDefault(), false)

	cfg, srcobject, sinkobjects, err := getPipelineCfgAndObjects()
	if err != nil {
		log.Error(err, "Failed to retrieve job spec from environment")
		os.Exit(1)
	}

	pipeline, err := buildPipelineFromCR(cfg, srcobject, sinkobjects)
	if err != nil {
		log.Error(err, "Failed to build pipeline from job spec")
		os.Exit(2)
	}

	pipeline.GetPipelineBus().AddWatch(func(msg *gst.Message) bool {
		switch msg.Type() {
		case gst.MessageEOS:
			log.Info("Received EOS, setting pipeline state to NULL")
			pipeline.BlockSetState(gst.StateNull)
			mainLoop.Quit()
			return false
		case gst.MessageError:
			err := msg.ParseError()
			log.Error(err, err.DebugString())
			os.Exit(3)
		}

		log.Info(msg.String())
		return true
	})

	pipeline.BlockSetState(gst.StatePlaying)

	go func() {
		var mc *minio.Client
		var err error
		var outbucket string
		var outpath string
		var g *graphviz.Graphviz

		for cfg.DoDotDump() {
			mc, err = util.GetMinIOClient(srcobject.Config.MinIO, util.MinIOSrcCredentialsFromEnv())
			if err != nil {
				log.Error(err, "Could not create src minio client for dot graph debugging")
				mc = nil
				break
			}
			outbucket = srcobject.Config.MinIO.GetBucket()
			outpath = cfg.GetDotPath(srcobject.Name)
			break
		}

		for range time.NewTicker(cfg.GetDotInterval()).C {

		DebugDotData:
			for cfg.DoDotDump() && mc != nil {
				if g == nil {
					g = graphviz.New()
				}
				var dotname, imgname string
				var dotdata, imgdata []byte

				dotdata = []byte(pipeline.DebugBinToDotData(gst.DebugGraphShowAll))

				if cfg.TimestampDotGraphs() {
					ts := time.Now().UTC().Format(time.RFC3339)
					dotname = path.Join(outpath, fmt.Sprintf("pipeline_%s.dot", ts))
					if render := cfg.GetDotRenderFormat(); render != "" {
						imgname = path.Join(outpath, fmt.Sprintf("pipeline_%s.%s", ts, strings.ToLower(render)))
					}
				} else {
					dotname = path.Join(outpath, "pipeline.dot")
					if render := cfg.GetDotRenderFormat(); render != "" {
						imgname = path.Join(outpath, fmt.Sprintf("pipeline.%s", strings.ToLower(render)))
					}
				}

				graph, err := graphviz.ParseBytes(dotdata)
				if err != nil {
					log.Error(err, "Failed to parse pipeline dot data")
					break DebugDotData
				}
				var buf bytes.Buffer
				switch strings.ToLower(cfg.GetDotRenderFormat()) {
				case "png":
					if err := g.Render(graph, graphviz.PNG, &buf); err != nil {
						log.Error(err, "Failed to convert dotdata to PNG")
						break DebugDotData
					}
					imgdata = buf.Bytes()
				case "svg":
					if err := g.Render(graph, graphviz.SVG, &buf); err != nil {
						log.Error(err, "Failed to convert dotdata to SVG")
						break DebugDotData
					}
					imgdata = buf.Bytes()
				case "jpg":
					if err := g.Render(graph, graphviz.JPG, &buf); err != nil {
						log.Error(err, "Failed to convert dotdata to JPG")
						break DebugDotData
					}
					imgdata = buf.Bytes()
				}

				_, err = mc.PutObject(context.Background(), outbucket, dotname, bytes.NewBuffer(dotdata), int64(len([]byte(dotdata))), minio.PutObjectOptions{
					ContentType: "application/octet-stream",
				})
				if err != nil {
					log.Error(err, "Failed to upload pipeline dot data")
					break DebugDotData
				}

				if imgdata != nil {
					if _, err := mc.PutObject(context.Background(), outbucket, imgname, bytes.NewBuffer(imgdata), int64(len(imgdata)), minio.PutObjectOptions{
						ContentType: "application/octet-stream",
					}); err != nil {
						log.Error(err, "Failed to upload rendered image of dot graph")
					}
				}

				break DebugDotData
			}

			posquery := gst.NewPositionQuery(gst.FormatTime)
			durquery := gst.NewDurationQuery(gst.FormatTime)
			pok := pipeline.Query(posquery)
			pipeline.Query(durquery) // we don't care if we don't have a return for this
			if !pok {
				log.Info("Failed to query the pipeline for the current position")
			}
			if pok {
				_, position := posquery.ParsePosition()
				_, duration := durquery.ParseDuration()
				log.Info(fmt.Sprintf("Current position %v/%v", time.Duration(position), time.Duration(duration)))
			}
		}
	}()

	mainLoop.Run()

	log.Info("Main loop has returned, ensuring pipeline has reached null state")
	for pipeline.GetState() == gst.StatePlaying {
	}

	log.Info("Pipeline finished", "State", pipeline.GetState())
}

func getPipelineCfgAndObjects() (cfg *pipelinesmeta.PipelineConfig, src *pipelinesmeta.Object, sinks []*pipelinesmeta.Object, err error) {
	cfg = &pipelinesmeta.PipelineConfig{}
	src = &pipelinesmeta.Object{}
	sinks = []*pipelinesmeta.Object{}
	if err = json.Unmarshal([]byte(os.Getenv(pipelinesmeta.JobPipelineConfigEnvVar)), cfg); err != nil {
		return
	}
	if err = json.Unmarshal([]byte(os.Getenv(pipelinesmeta.JobSrcObjectsEnvVar)), src); err != nil {
		return
	}
	if err = json.Unmarshal([]byte(os.Getenv(pipelinesmeta.JobSinkObjectsEnvVar)), &sinks); err != nil {
		return
	}
	return
}
