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

package util

import (
	"crypto/tls"
	"errors"
	"net/http"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"sigs.k8s.io/controller-runtime/pkg/client"

	pipelinesmeta "github.com/tinyzimmer/gst-pipeline-operator/apis/meta/v1"
	"github.com/tinyzimmer/gst-pipeline-operator/pkg/types"
)

// MinIOCredentialsGetter is a credential getter for minio clients. Various
// credentials sources are implemented in this package.
type MinIOCredentialsGetter interface {
	GetCredentials() (*credentials.Credentials, error)
}

// MinIOSinkCredentialsFromEnv returns a credentials getter that retrieves the credentials
// from the environment variables configured by the controller for the sink.
func MinIOSinkCredentialsFromEnv() MinIOCredentialsGetter { return &sinkCredentialsFromEnv{} }

type sinkCredentialsFromEnv struct{}

func (s *sinkCredentialsFromEnv) GetCredentials() (*credentials.Credentials, error) {
	return credentials.NewStaticV4(os.Getenv(pipelinesmeta.MinIOSinkAccessKeyIDEnvVar), os.Getenv(pipelinesmeta.MinIOSinkSecretAccessKeyEnvVar), ""), nil
}

// MinIOSrcCredentialsFromEnv returns a credentials getter that retrieves the credentials
// from the environment variables configured by the controller for the src.
func MinIOSrcCredentialsFromEnv() MinIOCredentialsGetter { return &srcCredentialsFromEnv{} }

type srcCredentialsFromEnv struct{}

func (s *srcCredentialsFromEnv) GetCredentials() (*credentials.Credentials, error) {
	return credentials.NewStaticV4(os.Getenv(pipelinesmeta.MinIOSrcAccessKeyIDEnvVar), os.Getenv(pipelinesmeta.MinIOSrcSecretAccessKeyEnvVar), ""), nil
}

// MinIOWatchCredentialsFromCR returns a credentials getter that uses the given client
// and CR to produce credentials to the bucket being watched for transformations.
func MinIOWatchCredentialsFromCR(client client.Client, cr types.Pipeline) MinIOCredentialsGetter {
	return &pipelineWatchCredentials{
		client: client,
		cr:     cr,
	}
}

type pipelineWatchCredentials struct {
	client client.Client
	cr     types.Pipeline
}

func (p *pipelineWatchCredentials) GetCredentials() (*credentials.Credentials, error) {
	srcConfig := p.cr.GetSrcConfig()
	if srcConfig == nil || srcConfig.MinIO == nil {
		return nil, errors.New("There is no MinIO configuration for this source")
	}
	return srcConfig.MinIO.GetStaticCredentials(p.client, p.cr.GetNamespace())
}

// GetMinIOClient is a utility function for returning a MinIO client to the given
// configuration.
func GetMinIOClient(cfg *pipelinesmeta.MinIOConfig, credsGetter MinIOCredentialsGetter) (*minio.Client, error) {
	transport := http.DefaultTransport.(*http.Transport).Clone()

	if cfg.GetSecure() {
		certPool, err := cfg.GetRootCAs()
		if err != nil {
			return nil, err
		}
		if transport.TLSClientConfig == nil {
			transport.TLSClientConfig = &tls.Config{}
		}
		transport.TLSClientConfig.RootCAs = certPool
		transport.TLSClientConfig.InsecureSkipVerify = cfg.GetSkipVerify()
	}

	creds, err := credsGetter.GetCredentials()
	if err != nil {
		return nil, err
	}

	return minio.New(cfg.GetEndpoint(), &minio.Options{
		Creds:     creds,
		Secure:    cfg.GetSecure(),
		Region:    cfg.GetRegion(),
		Transport: transport,
	})
}
