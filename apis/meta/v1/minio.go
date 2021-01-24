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

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"path"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/minio/minio-go/v7/pkg/credentials"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// MinIOConfig defines a source or sink location for pipelines.
type MinIOConfig struct {
	// The MinIO endpoint *without* the leading `http(s)://`.
	Endpoint string `json:"endpoint,omitempty"`
	// Do not use TLS when communicating with the MinIO API.
	InsecureNoTLS bool `json:"insecureNoTLS,omitempty"`
	// A base64-endcoded PEM certificate chain to use when verifying the certificate
	// supplied by the MinIO server.
	EndpointCA string `json:"endpointCA,omitempty"`
	// Skip verification of the certificate supplied by the MinIO server.
	InsecureSkipVerify bool `json:"insecureSkipVerify,omitempty"`
	// The region to connect to in MinIO.
	Region string `json:"region,omitempty"`
	// In the context of a src config, the bucket to watch for objects to pass through
	// the pipeline. In the context of a sink config, the bucket to save processed objects.
	Bucket string `json:"bucket,omitempty"`
	// In the context of a src config, a directory prefix to match for objects to be sent
	// through the pipeline. An empty value means ALL objects in the bucket, or the equivalent of
	// `/`. In the context of a sink config, a go-template to use for the destination name. The
	// template allows sprig functions and is passed the value "SrcName" representing the base of the key
	// of the object that triggered the pipeline, and "SrcExt" with the extension. An empty value represents
	// using the same key as the source which would only work for objects being processed to different
	// buckets and prefixes.
	Prefix string `json:"key,omitempty"`
	// The secret that contains the credentials for connecting to MinIO. The secret must contain
	// two keys. The `access-key-id` key must contain the contents of the Access Key ID. The
	// `secret-access-key` key must contain the contents of the Secret Access Key.
	CredentialsSecret *corev1.LocalObjectReference `json:"credentialsSecret,omitempty"`
}

// GetEndpoint returns the API endpoint for this configuration.
func (m *MinIOConfig) GetEndpoint() string { return m.Endpoint }

// GetSecure returns whether to use HTTPS for API communication.
func (m *MinIOConfig) GetSecure() bool { return !m.InsecureNoTLS }

// GetSkipVerify returns where to skip TLS verification of the server certificate.
func (m *MinIOConfig) GetSkipVerify() bool { return !m.InsecureSkipVerify }

// GetBucket returns the bucket for this configuration.
func (m *MinIOConfig) GetBucket() string { return m.Bucket }

// GetPrefix returns the prefix for this configuration.
func (m *MinIOConfig) GetPrefix() string { return m.Prefix }

// GetRegion returns the region to connect the client to.
func (m *MinIOConfig) GetRegion() string {
	if m.Region == "" {
		return DefaultRegion
	}
	return m.Region
}

// GetRootPEM returns the raw PEM of the root certificate included in the configuration.
func (m *MinIOConfig) GetRootPEM() ([]byte, error) {
	if m.EndpointCA == "" {
		return nil, nil
	}
	return base64.StdEncoding.DecodeString(m.EndpointCA)
}

// GetRootCAs returns an x509.CertPool for any provided CA certificates.
func (m *MinIOConfig) GetRootCAs() (*x509.CertPool, error) {
	certPEM, err := m.GetRootPEM()
	if err != nil {
		return nil, err
	}
	if certPEM == nil {
		return nil, nil
	}
	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(certPEM); !ok {
		return nil, errors.New("Failed to append CA PEM certificates to cert pool")
	}
	return certPool, nil
}

// GetCredentialsSecret returns the name of the credentials secret.
func (m *MinIOConfig) GetCredentialsSecret() (string, error) {
	if m.CredentialsSecret == nil {
		return "", errors.New("No secret reference included in the CR for endpoint credentials")
	}
	return m.CredentialsSecret.Name, nil
}

// GetCredentials attemps to retrieve the access key ID and secret access key for this config.
func (m *MinIOConfig) GetCredentials(client client.Client, namespace string) (accessKeyID, secretAccessKey string, err error) {
	secretName, err := m.GetCredentialsSecret()
	if err != nil {
		return "", "", err
	}
	secret := &corev1.Secret{}
	if err := client.Get(context.TODO(), types.NamespacedName{Name: secretName, Namespace: namespace}, secret); err != nil {
		return "", "", err
	}
	accessKeyIDRaw, ok := secret.Data[AccessKeyIDKey]
	if !ok {
		return "", "", fmt.Errorf("No %s in secret %s/%s", AccessKeyIDKey, namespace, secretName)
	}
	secretAccessKeyRaw, ok := secret.Data[SecretAccessKeyKey]
	if !ok {
		return "", "", fmt.Errorf("No %s in secret %s/%s", SecretAccessKeyKey, namespace, secretName)
	}
	return string(accessKeyIDRaw), string(secretAccessKeyRaw), nil
}

// GetStaticCredentials attempts to return API credentials for MinIO using the given
// client looking in the given namespace.
func (m *MinIOConfig) GetStaticCredentials(client client.Client, namespace string) (*credentials.Credentials, error) {
	accessKeyID, secretAccessKey, err := m.GetCredentials(client, namespace)
	if err != nil {
		return nil, err
	}
	return credentials.NewStaticV4(accessKeyID, secretAccessKey, ""), nil
}

// GetDestinationKey computes what the destination object's name should be based on the
// given source object name. If a template is present and it fails to execute, it is logged
// and the default behavior is returned.
func (m *MinIOConfig) GetDestinationKey(objectKey string) string {
	tmpl := m.GetPrefix()
	for tmpl != "" {
		ext := path.Ext(objectKey)
		name := path.Base(strings.TrimSuffix(objectKey, ext))
		var buf bytes.Buffer
		var t *template.Template
		t, err := template.New("").Funcs(sprig.TxtFuncMap()).Parse(tmpl)
		if err != nil {
			fmt.Println(err)
			break
		}
		t.Execute(&buf, map[string]string{
			"SrcName": name,
			"SrcExt":  ext,
		})
		return buf.String()
	}
	return path.Join(strings.TrimSuffix(m.GetPrefix(), "/"), path.Base(objectKey))
}
