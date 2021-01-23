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
	"io/ioutil"
	"os"
	"strconv"

	"github.com/tinyzimmer/go-glib/glib"
	"github.com/tinyzimmer/go-gst/gst"

	pipelinesmeta "github.com/tinyzimmer/gst-pipeline-operator/apis/meta/v1"
)

func makeSrcElement(objCfg *pipelinesmeta.Object) (*gst.Element, error) {
	elem, err := gst.NewElement("miniosrc")
	if err != nil {
		return nil, err
	}

	cfg := objCfg.Config.MinIO // TODO
	log.Info("Creating src element", "Config", *cfg, "Key", objCfg.Name)

	elem.SetProperty("endpoint", cfg.GetEndpoint())
	elem.SetProperty("use-tls", cfg.GetSecure())
	elem.SetProperty("region", cfg.GetRegion())
	elem.SetProperty("bucket", cfg.GetBucket())
	elem.SetProperty("key", objCfg.Name)
	elem.SetProperty("access-key-id", os.Getenv(pipelinesmeta.MinIOSrcAccessKeyIDEnvVar))
	elem.SetProperty("secret-access-key", os.Getenv(pipelinesmeta.MinIOSrcSecretAccessKeyEnvVar))

	rootCA, err := cfg.GetRootPEM()
	if err != nil {
		return nil, err
	}
	if rootCA != nil {
		if err := ioutil.WriteFile("/tmp/ca-src.crt", rootCA, 0644); err != nil {
			return nil, err
		}
		elem.SetProperty("ca-cert-file", "/tmp/ca-src.crt")
	}

	return elem, nil
}

func makeSinkElement(objCfg *pipelinesmeta.Object) (*gst.Element, *pipelinesmeta.GstElementConfig, error) {
	elem, err := gst.NewElement("miniosink")
	if err != nil {
		return nil, nil, err
	}

	cfg := objCfg.Config.MinIO // TODO
	log.Info("Creating sink element", "Config", *cfg, "Key", objCfg.Name)

	elem.SetProperty("endpoint", cfg.GetEndpoint())
	elem.SetProperty("use-tls", cfg.GetSecure())
	elem.SetProperty("region", cfg.GetRegion())
	elem.SetProperty("bucket", cfg.GetBucket())
	elem.SetProperty("key", objCfg.Name)
	elem.SetProperty("access-key-id", os.Getenv(pipelinesmeta.MinIOSinkAccessKeyIDEnvVar))
	elem.SetProperty("secret-access-key", os.Getenv(pipelinesmeta.MinIOSinkSecretAccessKeyEnvVar))

	rootCA, err := cfg.GetRootPEM()
	if err != nil {
		return nil, nil, err
	}
	if rootCA != nil {
		if err := ioutil.WriteFile("/tmp/ca-sink.crt", rootCA, 0644); err != nil {
			return nil, nil, err
		}
		elem.SetProperty("ca-cert-file", "/tmp/ca-sink.crt")
	}

	elemcfg := &pipelinesmeta.GstElementConfig{}
	elemcfg.SetPipelineName(elem.GetName())

	return elem, elemcfg, nil
}

func makeElement(cfg *pipelinesmeta.GstElementConfig) (*gst.Element, error) {
	elem, err := gst.NewElement(cfg.Name)
	if err != nil {
		return nil, err
	}
	if cfg.Properties == nil {
		return elem, nil
	}
	for propName, propValue := range cfg.Properties {
		propType, err := elem.GetPropertyType(propName)
		if err != nil {
			return nil, err
		}
		gval, err := goStringToGValueByType(propType, propValue)
		if err != nil {
			return nil, err
		}
		if err := elem.SetPropertyValue(propName, gval); err != nil {
			return nil, err
		}
	}
	return elem, nil
}

func goStringToGValueByType(t glib.Type, s string) (*glib.Value, error) {
	value, err := glib.ValueInit(t)
	if err != nil {
		return nil, err
	}
	switch t {

	case glib.TYPE_CHAR:
		i, err := strconv.ParseInt(s, 10, 8)
		if err != nil {
			return nil, err
		}
		value.SetSChar(int8(i))

	case glib.TYPE_UCHAR:
		i, err := strconv.ParseInt(s, 10, 8)
		if err != nil {
			return nil, err
		}
		value.SetUChar(uint8(i))

	case glib.TYPE_BOOLEAN:
		switch s {
		case "true":
			value.SetBool(true)
		case "false":
			value.SetBool(false)
		default:
			return nil, fmt.Errorf("Unrecognized bool value: %s", s)
		}

	case glib.TYPE_INT:
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}
		value.SetInt(int(i))

	case glib.TYPE_ENUM:
		// Need ENUM helpers in glib
		return nil, errors.New("enums not implemented")

	case glib.TYPE_INT64:
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}
		value.SetInt64(int64(i))

	case glib.TYPE_UINT:
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}
		value.SetUInt(uint(i))

	case glib.TYPE_FLAGS:
		// Need helpers
		return nil, errors.New("flags not implemented")

	case glib.TYPE_UINT64:
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}
		value.SetUInt64(uint64(i))

	case glib.TYPE_FLOAT:
		f, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return nil, err
		}
		value.SetFloat(float32(f))

	case glib.TYPE_DOUBLE:
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, err
		}
		value.SetDouble(float64(f))

	case glib.TYPE_STRING:
		value.SetString(s)

	case glib.TYPE_BOXED:
		// Assumes caps for now
		caps := gst.NewCapsFromString(s)
		if caps == nil {
			return nil, fmt.Errorf("Could not parse provided caps to string: %s", s)
		}
		return caps.ToGValue(), nil

		// need to handle other GST types like fractions

	}
	return value, nil
}
