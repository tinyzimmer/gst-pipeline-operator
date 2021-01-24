## kVDI CRD Reference

### Packages:

-   [meta.gst.io/v1](#meta.gst.io%2fv1)

Types

-   [DebugConfig](#DebugConfig)
-   [DotConfig](#DotConfig)
-   [ElementConfig](#ElementConfig)
-   [GstElementConfig](#GstElementConfig)
-   [GstLaunchConfig](#GstLaunchConfig)
-   [MinIOConfig](#MinIOConfig)
-   [Object](#Object)
-   [PipelineConfig](#PipelineConfig)
-   [PipelineKind](#PipelineKind)
-   [PipelineReference](#PipelineReference)
-   [PipelineState](#PipelineState)
-   [SourceSinkConfig](#SourceSinkConfig)
-   [StreamType](#StreamType)

## meta.gst.io/v1

Package v1 contains API Schema definitions for the pipelines v1 API
group

Resource Types:

### DebugConfig

(*Appears on:* [PipelineConfig](#PipelineConfig))

DebugConfig represents debug configurations for a GStreamer pipeline.

<table>
<colgroup>
<col style="width: 50%" />
<col style="width: 50%" />
</colgroup>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>logLevel</code><br />
<em>int</em></td>
<td><p>The level of log output to produce from the gstreamer process. This value gets set to the GST_DEBUG variable. Defaults to INFO level (4). Higher numbers mean more output.</p></td>
</tr>
<tr class="even">
<td><code>dot</code><br />
<em><a href="#DotConfig">DotConfig</a></em></td>
<td><p>Dot specifies to dump a dot file of the pipeline layout for debugging.</p></td>
</tr>
</tbody>
</table>

### DotConfig

(*Appears on:* [DebugConfig](#DebugConfig))

DotConfig represents a configuration for the dot output of a pipeline.

<table>
<colgroup>
<col style="width: 50%" />
<col style="width: 50%" />
</colgroup>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>path</code><br />
<em>string</em></td>
<td><p>The path to save files. The configuration other than the path is assumed to be that of the source of the pipeline. For example, for a MinIO source, this should be a prefix in the same bucket as the source (but not overlapping with the watch prefix otherwise an infinite loop will happen). The files will be saved in directories matching the source object’s name with the _debug suffix.</p></td>
</tr>
<tr class="even">
<td><code>render</code><br />
<em>string</em></td>
<td><p>Specify to also render the pipeline graph to images in the given format. Accepted formats are png, svg, or jpg.</p></td>
</tr>
<tr class="odd">
<td><code>timestamped</code><br />
<em>bool</em></td>
<td><p>Whether to save timestamped versions of the pipeline layout. This will produce a new graph for every interval specified by Interval. The default is to only keep the latest graph.</p></td>
</tr>
<tr class="even">
<td><code>interval</code><br />
<em>int</em></td>
<td><p>The interval in seconds to save pipeline graphs. Defaults to every 3 seconds.</p></td>
</tr>
</tbody>
</table>

### ElementConfig

(*Appears on:* [GstElementConfig](#GstElementConfig))

ElementConfig represents the configuration of a single element in a
transform pipeline.

<table>
<colgroup>
<col style="width: 50%" />
<col style="width: 50%" />
</colgroup>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>name</code><br />
<em>string</em></td>
<td><p>The name of the element. See the GStreamer plugin documentation for a comprehensive list of all the plugins available. Custom pipeline images can also be used that are prebaked with additional plugins.</p></td>
</tr>
<tr class="even">
<td><code>alias</code><br />
<em>string</em></td>
<td><p>Applies an alias to this element in the pipeline configuration. This allows you to specify an element block with this value as the name and have it act as a “goto” or “linkto” while building the pipeline. Note that the aliases “video-out” and “audio-out” are reserved for internal use.</p></td>
</tr>
<tr class="odd">
<td><code>goto</code><br />
<em>string</em></td>
<td><p>The alias to an element to treat as this configuration. Useful for directing the output of elements with multiple src pads, such as decodebin.</p></td>
</tr>
<tr class="even">
<td><code>linkto</code><br />
<em>string</em></td>
<td><p>The alias to an element to link the previous element’s sink pad to. Useful for directing the branches of a multi-stream pipeline to a muxer. A linkto almost always needs to be followed by a goto, except when the element being linked to is next in the pipeline, in which case you can omit the linkto entirely.</p></td>
</tr>
<tr class="odd">
<td><code>properties</code><br />
<em>map[string]string</em></td>
<td><p>Optional properties to apply to this element. To not piss off the CRD generator values are declared as a string, but almost anything that can be passed to gst-launch-1.0 will work. Caps will be parsed from their string representation.</p></td>
</tr>
</tbody>
</table>

### GstElementConfig

GstElementConfig is an extension of the ElementConfig struct providing
private fields for internal tracking while building a dynamic pipeline.

<table>
<colgroup>
<col style="width: 50%" />
<col style="width: 50%" />
</colgroup>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>ElementConfig</code><br />
<em><a href="#ElementConfig">ElementConfig</a></em></td>
<td></td>
</tr>
<tr class="even">
<td><code>pipelineName</code><br />
<em>string</em></td>
<td></td>
</tr>
<tr class="odd">
<td><code>peers</code><br />
<em><a href="#GstElementConfig">[]*github.com/tinyzimmer/gst-pipeline-operator/apis/meta/v1.GstElementConfig</a></em></td>
<td></td>
</tr>
</tbody>
</table>

GstLaunchConfig
(`[]*github.com/tinyzimmer/gst-pipeline-operator/apis/meta/v1.GstElementConfig`
alias)

GstLaunchConfig is a slice of ElementConfigs that contain internal
fields used for dynamic linking.

### MinIOConfig

(*Appears on:* [SourceSinkConfig](#SourceSinkConfig))

MinIOConfig defines a source or sink location for pipelines.

<table>
<colgroup>
<col style="width: 50%" />
<col style="width: 50%" />
</colgroup>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>endpoint</code><br />
<em>string</em></td>
<td><p>The MinIO endpoint <em>without</em> the leading <code>http(s)://</code>.</p></td>
</tr>
<tr class="even">
<td><code>insecureNoTLS</code><br />
<em>bool</em></td>
<td><p>Do not use TLS when communicating with the MinIO API.</p></td>
</tr>
<tr class="odd">
<td><code>endpointCA</code><br />
<em>string</em></td>
<td><p>A base64-endcoded PEM certificate chain to use when verifying the certificate supplied by the MinIO server.</p></td>
</tr>
<tr class="even">
<td><code>insecureSkipVerify</code><br />
<em>bool</em></td>
<td><p>Skip verification of the certificate supplied by the MinIO server.</p></td>
</tr>
<tr class="odd">
<td><code>region</code><br />
<em>string</em></td>
<td><p>The region to connect to in MinIO.</p></td>
</tr>
<tr class="even">
<td><code>bucket</code><br />
<em>string</em></td>
<td><p>In the context of a src config, the bucket to watch for objects to pass through the pipeline. In the context of a sink config, the bucket to save processed objects.</p></td>
</tr>
<tr class="odd">
<td><code>key</code><br />
<em>string</em></td>
<td><p>In the context of a src config, a directory prefix to match for objects to be sent through the pipeline. An empty value means ALL objects in the bucket, or the equivalent of <code>/</code>. In the context of a sink config, a go-template to use for the destination name. The template allows sprig functions and is passed the value “SrcName” representing the base of the key of the object that triggered the pipeline, and “SrcExt” with the extension. An empty value represents using the same key as the source which would only work for a objects being processed to different buckets and prefixes. When splitting streams the prefixes “audio<em>” and “video</em>” respectively will be added to the resulting filenames.</p></td>
</tr>
<tr class="even">
<td><code>credentialsSecret</code><br />
<em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/#localobjectreference-v1-core">Kubernetes core/v1.LocalObjectReference</a></em></td>
<td><p>The secret that contains the credentials for connecting to MinIO. The secret must contain two keys. The <code>access-key-id</code> key must contain the contents of the Access Key ID. The <code>secret-access-key</code> key must contain the contents of the Secret Access Key.</p></td>
</tr>
</tbody>
</table>

### Object

Object represents either a source or destination object for a job.

<table>
<colgroup>
<col style="width: 50%" />
<col style="width: 50%" />
</colgroup>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>name</code><br />
<em>string</em></td>
<td><p>The actual name for the object being read or written to. In the context of a source object this is pulled from a watch event. In the context of a destination this is computed by the controller from the user supplied configuration.</p></td>
</tr>
<tr class="even">
<td><code>config</code><br />
<em><a href="#SourceSinkConfig">SourceSinkConfig</a></em></td>
<td><p>The endpoint and bucket configurations for the object.</p></td>
</tr>
<tr class="odd">
<td><code>streamType</code><br />
<em><a href="#StreamType">StreamType</a></em></td>
<td><p>The type of the stream for this object. Only applies to sinks. For a split transform pipeline there will be an Object for each stream. Otherwise there will be a single object with a StreamTypeAll.</p></td>
</tr>
</tbody>
</table>

### PipelineConfig

PipelineConfig represents a series of elements through which to pass the
contents of processed objects.

<table>
<colgroup>
<col style="width: 50%" />
<col style="width: 50%" />
</colgroup>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>image</code><br />
<em>string</em></td>
<td><p>The image to use to run a/v processing pipelines.</p></td>
</tr>
<tr class="even">
<td><code>debug</code><br />
<em><a href="#DebugConfig">DebugConfig</a></em></td>
<td><p>Debug configurations for the pipeline</p></td>
</tr>
<tr class="odd">
<td><code>elements</code><br />
<em><a href="#ElementConfig">[]*github.com/tinyzimmer/gst-pipeline-operator/apis/meta/v1.ElementConfig</a></em></td>
<td><p>A list of element configurations in the order they will be used in the pipeline. Using these is mutually exclusive with a decodebin configuration. This only really works for linear pipelines. That is to say, not the syntax used by <code>gst-launch-1.0</code> that allows naming elements and referencing them later in the pipeline. For complex handling of multiple streams decodebin will still be better to work with for now, despite its shortcomings.</p></td>
</tr>
<tr class="even">
<td><code>resources</code><br />
<em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/#resourcerequirements-v1-core">Kubernetes core/v1.ResourceRequirements</a></em></td>
<td><p>Resource restraints to place on jobs created for this pipeline.</p></td>
</tr>
</tbody>
</table>

PipelineKind (`string` alias)

(*Appears on:* [PipelineReference](#PipelineReference))

PipelineKind represents a type of pipeline.

### PipelineReference

PipelineReference is used to refer to the pipeline that holds the
configuration for a given job.

<table>
<colgroup>
<col style="width: 50%" />
<col style="width: 50%" />
</colgroup>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>name</code><br />
<em>string</em></td>
<td><p>Name is the name of the Pipeline CR</p></td>
</tr>
<tr class="even">
<td><code>kind</code><br />
<em><a href="#PipelineKind">PipelineKind</a></em></td>
<td><p>Kind is the type of the Pipeline CR</p></td>
</tr>
</tbody>
</table>

PipelineState (`string` alias)

PipelineState represents the state of a Pipeline CR.

### SourceSinkConfig

(*Appears on:* [Object](#Object))

SourceSinkConfig is used to declare configurations related to the
retrieval or saving of pipeline objects.

<table>
<colgroup>
<col style="width: 50%" />
<col style="width: 50%" />
</colgroup>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>minio</code><br />
<em><a href="#MinIOConfig">MinIOConfig</a></em></td>
<td><p>Configurations for a MinIO source or sink</p></td>
</tr>
</tbody>
</table>

StreamType (`string` alias)

(*Appears on:* [Object](#Object))

StreamType represents a type of stream found in a source input, or
designated for an output.

------------------------------------------------------------------------

*Generated with `gen-crd-api-reference-docs` on git commit `8269c7e`.*
