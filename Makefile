# Current Operator version
VERSION ?= latest
# Default bundle image tag
BUNDLE_IMG ?= $(REPO)/controller-bundle:$(VERSION)
# Options for 'bundle-build'
ifneq ($(origin CHANNELS), undefined)
BUNDLE_CHANNELS := --channels=$(CHANNELS)
endif
ifneq ($(origin DEFAULT_CHANNEL), undefined)
BUNDLE_DEFAULT_CHANNEL := --default-channel=$(DEFAULT_CHANNEL)
endif
BUNDLE_METADATA_OPTS ?= $(BUNDLE_CHANNELS) $(BUNDLE_DEFAULT_CHANNEL)

REPO ?= ghcr.io/tinyzimmer/gst-pipeline-operator

# Image URL to use all building/pushing image targets
IMG ?= $(REPO)/controller:$(VERSION)
CRD_OPTIONS ?= "crd:crdVersions=v1"

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: manager

# Run tests
ENVTEST_ASSETS_DIR = $(shell pwd)/testbin
test: generate fmt vet manifests
	mkdir -p $(ENVTEST_ASSETS_DIR)
	test -f $(ENVTEST_ASSETS_DIR)/setup-envtest.sh || curl -sSLo $(ENVTEST_ASSETS_DIR)/setup-envtest.sh https://raw.githubusercontent.com/kubernetes-sigs/controller-runtime/v0.6.3/hack/setup-envtest.sh
	source $(ENVTEST_ASSETS_DIR)/setup-envtest.sh; fetch_envtest_tools $(ENVTEST_ASSETS_DIR); setup_envtest_env $(ENVTEST_ASSETS_DIR); go test ./... -coverprofile cover.out

LDFLAGS ?= -X github.com/tinyzimmer/gst-pipeline-operator/pkg/version.Version=$(VERSION) \
				-X github.com/tinyzimmer/gst-pipeline-operator/pkg/version.GitCommit=$(shell git rev-parse HEAD)

# Build manager binary
manager: generate fmt vet
	go build -o bin/manager -ldflags="$(LDFLAGS)" main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet manifests
	go run ./main.go

# Install CRDs into a cluster
install: manifests kustomize
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

# Uninstall CRDs from a cluster
uninstall: manifests kustomize
	$(KUSTOMIZE) build config/crd | kubectl delete -f -

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests kustomize
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build config/default | kubectl apply -f -

deploy-manifest: manifests kustomize
	mkdir -p deploy/manifests
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build config/default > deploy/manifests/gst-pipeline-operator-full.yaml

# Generate manifests e.g. CRD, RBAC etc.
manifests: controller-gen
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Generate code
generate: controller-gen
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

# Build the docker image
docker-build:
	docker build . -t ${IMG} --build-arg LDFLAGS="$(LDFLAGS)"

# Push the docker image
docker-push:
	docker push ${IMG}

# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen 2> /dev/null))
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.4.0 ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen 2> /dev/null)
endif

kustomize:
ifeq (, $(shell which kustomize 2> /dev/null))
	@{ \
	set -e ;\
	KUSTOMIZE_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$KUSTOMIZE_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/kustomize/kustomize/v3@v3.5.4 ;\
	rm -rf $$KUSTOMIZE_GEN_TMP_DIR ;\
	}
KUSTOMIZE=$(GOBIN)/kustomize
else
KUSTOMIZE=$(shell which kustomize 2> /dev/null)
endif

# Generate bundle manifests and metadata, then validate generated files.
.PHONY: bundle
bundle: manifests kustomize
	operator-sdk generate kustomize manifests -q
	cd config/manager && $(KUSTOMIZE) edit set image controller=$(IMG)
	$(KUSTOMIZE) build config/manifests | operator-sdk generate bundle -q --overwrite --version $(VERSION) $(BUNDLE_METADATA_OPTS)
	operator-sdk bundle validate ./bundle

# Build the bundle image.
.PHONY: bundle-build
bundle-build:
	docker build -f bundle.Dockerfile -t $(BUNDLE_IMG) .

## Custom Targets

GST_IMAGE ?= $(REPO)/gstreamer:$(VERSION)
docker-gst-build:
	docker build -f gst/Dockerfile -t $(GST_IMAGE) .

docker-gst-push:
	docker push $(GST_IMAGE)

docker-build-all: docker-build docker-gst-build

K3D ?= $(GOBIN)/k3d
$(K3D):
	curl -s https://raw.githubusercontent.com/rancher/k3d/main/install.sh | K3D_INSTALL_DIR=$(GOBIN) bash -s -- --no-sudo

CLUSTER_NAME ?= gst
CLUSTER_KUBECONFIG ?= $(CURDIR)/kubeconfig.yaml
local-cluster:
	$(K3D) cluster create $(CLUSTER_NAME) \
		--update-default-kubeconfig=false \
		-p 9000:9000@loadbalancer
	$(K3D) kubeconfig get $(CLUSTER_NAME) > $(CLUSTER_KUBECONFIG)

local-import: docker-build-all
	$(K3D) image import --cluster=$(CLUSTER_NAME) $(IMG) $(GST_IMAGE)

local-deploy: local-import deploy-manifest
	KUBECONFIG=$(CLUSTER_KUBECONFIG) kubectl apply -f deploy/manifests/gst-pipeline-operator-full.yaml

TEST_HELM ?= KUBECONFIG=$(CLUSTER_KUBECONFIG) helm
local-minio:
	$(TEST_HELM) repo add minio https://helm.min.io/
	$(TEST_HELM) repo update
	$(TEST_HELM) install \
		--set service.type=LoadBalancer \
		--set accessKey=accesskey \
		--set secretKey=secretkey \
		--set buckets[0].name=gst-processing \
		--set buckets[0].policy=public \
		minio minio/minio

local-samples:
	KUBECONFIG=$(CLUSTER_KUBECONFIG) kubectl apply -f config/samples/minio_credentials.yaml
	KUBECONFIG=$(CLUSTER_KUBECONFIG) kubectl apply -f config/samples/pipelines_v1_transform.yaml
	KUBECONFIG=$(CLUSTER_KUBECONFIG) kubectl apply -f config/samples/pipelines_v1_splittransform.yaml

local-full: local-cluster local-minio local-deploy

delete-local-cluster:
	$(K3D) cluster delete $(CLUSTER_NAME)
