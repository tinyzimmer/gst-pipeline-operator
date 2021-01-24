# gst-pipeline-operator

A Kubernetes operator for running audio/video processing pipelines

  - [API Reference](doc/pipelines.md)

## Quickstart

This project is very young and does not have a ton of features yet. For the POC it can watch MinIO buckets
for files to process, and feeds them through GStreamer pipelines defined in Custom Resources. 
The below guide walks you through setting up a cluster to process your own A/V in the same way.

### Setting up the Cluster

First you'll need a running Kubernetes cluster. For the purpose of the quickstart we'll use [`k3d`](https://github.com/rancher/k3d) to create the cluster.

```bash
k3d cluster create -p 9000:9000@loadbalancer  # Expose port 9000 for minio later
```

The operator currently only supports MinIO as a source or destination for pipeline objects.
The intention is to extend this to support other event sources/destinations (e.g. NFS, PVs, etc.).
For testing purposes, you can spin up a single node MinIO server inside the Kubernetes cluster using `helm`.

```bash
helm repo add minio https://helm.min.io/
helm repo update

# You may want to change some of these configurations for yourself. The examples later in the Quickstart
# will assume these values, and you should substitute for the ones you chose instead.
helm install minio minio/minio \
    --set service.type=LoadBalancer \
    --set accessKey=accesskey \
    --set secretKey=secretkey \
    --set buckets[0].name=gst-processing \
    --set buckets[0].policy=public
```

Finally you need to create a **Secret** holding the credentials for the operator and GStreamer pipelines to connect to MinIO.

```bash
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Secret
metadata:
  name: minio-credentials
  namespace: default
data:
  access-key-id: YWNjZXNza2V5        # The base64 encoded access key ID for MinIO
  secret-access-key: c2VjcmV0a2V5    # The base64 encoded secret access key for MinIO
EOF
```

### Install the Operator

With the above all done you can now install the `gst-pipeline-operator`. 
This repository contains a bundle manifest (with potential `helm` charts later) for installing all the required components.

```bash
kubectl apply -f https://raw.githubusercontent.com/tinyzimmer/gst-pipeline-operator/main/deploy/manifests/gst-pipeline-operator-full.yaml
```

The operator will be installed to the `gst-system` namespace. However, you can provision pipeline resources in any namespace you want. The following examples will all use the `default` namespace, which is also where we put the secret with the MinIO credentials.

There are currently two **Pipeline** CRs provided by the operator, with more to come later since the sky is the limit with GStreamer.

 - `Transform` - This pipeline takes objects dropped in the source bucket, feeds them through your pipeline, and drops the output into the defined sink bucket.
 - `SplitTransform` - This pipeline is the same as `Transform` except it can be used to separate video from audio in source files.

For the quickstart we'll do a simple `Transform` pipeline, but you can also find more examples [here](config/samples).

```yaml
---
apiVersion: pipelines.gst.io/v1
kind: Transform
metadata:
  name: mp4-converter
spec:
  # Globals will be merged into the `src` and `sink` configs during processing.
  # This is useful if all operations are happening against the same MinIO server
  # and buckets. You can also direct output to/from different servers and buckets
  # by declaring those values in their respective areas instead of here.
  globals:
    minio:
      endpoint: "minio.default.svc.cluster.local:9000"   # The endpoint for the MinIO server
      insecureNoTLS: true                                # Use HTTP
      region: us-east-1                                  # The region of the bucket
      bucket: gst-processing                             # The bucket to watch for files
      credentialsSecret:                                 # The secret containing READ credentials for MinIO
        name: minio-credentials
  src:
    minio:
      key: drop/   # Watch all files placed in the drop/ prefix
  sink:
    minio:
      key: "mp4/{{ .SrcName }}.mp4"   # Generate an output name from this template. If the src file was called
  pipeline:                           # drop/video.mkv then this would evaluate to mp4/video.mp4.
    debug:
      dot:
        path: debug/  # Optionally dump DOT graphs to the debug/ prefix for each pipeline
        render: png   # Optionally render those DOT graphs to PNG in addition to the DOT format.

    # The pipeline definition. This a "yamlized" version of the better known gst-launch-1.0 syntax.
    elements:
      - name: decodebin
        alias: dbin          # The same as applying a `name` in gst-launch-1.0

      - goto: dbin           # Take a compatible sink pad from the decodebin
      - name: queue
      - name: audioconvert
      - name: audioresample
      - name: voaacenc
      - linkto: mux          # Link the output of this chain to the `mux` element

      - goto: dbin           # Go back to the decodebin and take the next compatible sink pad
      - name: queue
      - name: videoconvert
      - name: x264enc

      - name: mp4mux         # Joins the output from the audio stream and the video stream into an mp4
        alias: mux           # container.
    
      # The last element evaluated in the pipeline (either in order or via goto/linkto) has its output
      # sent to the MinIO output object.
```

You can now log-in to the MinIO server at `127.0.0.1:9000` and place files in `gst-processing:drop/` to have them handled.
The operator provides another CRD, `jobs.pipelines` for tracking the state of processing jobs.
Down the road, a UI is envisioned for getting a bird's eye view of all the processing happening in the cluster.