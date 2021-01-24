## GST Pipelines CRD Reference

### Packages:

-   [pipelines.gst.io/v1](#pipelines.gst.io%2fv1)

Types

-   [Job](#Job)
-   [JobSpec](#JobSpec)
-   [JobState](#JobState)
-   [SplitTransform](#SplitTransform)
-   [SplitTransformSpec](#SplitTransformSpec)
-   [Transform](#Transform)
-   [TransformSpec](#TransformSpec)

## pipelines.gst.io/v1

Package v1 file doc.go required for the doc generator to register this
as an API

Resource Types:

### Job

Job is the Schema for the jobs API

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
<td><code>metadata</code> <em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/#objectmeta-v1-meta">Kubernetes metav1.ObjectMeta</a></em></td>
<td>Refer to the Kubernetes API documentation for the fields of the <code>metadata</code> field.</td>
</tr>
<tr class="even">
<td><code>spec</code> <em><a href="#JobSpec">JobSpec</a></em></td>
<td><br />
<br />

<table>
<tbody>
<tr class="odd">
<td><code>pipelineRef</code> <em><a href="meta.md#PipelineReference">metav1.PipelineReference</a></em></td>
<td><p>A reference to the pipeline for this job’s configuration.</p></td>
</tr>
<tr class="even">
<td><code>src</code> <em><a href="meta.md#Object">metav1.Object</a></em></td>
<td><p>The source object for the pipeline.</p></td>
</tr>
<tr class="odd">
<td><code>sinks</code> <em><a href="meta.md#Object">metav1.Object</a></em></td>
<td><p>The output objects for the pipeline.</p></td>
</tr>
</tbody>
</table></td>
</tr>
<tr class="odd">
<td><code>status</code> <em><a href="#JobStatus">JobStatus</a></em></td>
<td></td>
</tr>
</tbody>
</table>

### JobSpec

(*Appears on:* [Job](#Job))

JobSpec defines the desired state of Job

<table>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>pipelineRef</code> <em><a href="meta.md#PipelineReference">metav1.PipelineReference</a></em></td>
<td><p>A reference to the pipeline for this job’s configuration.</p></td>
</tr>
<tr class="even">
<td><code>src</code> <em><a href="meta.md#Object">metav1.Object</a></em></td>
<td><p>The source object for the pipeline.</p></td>
</tr>
<tr class="odd">
<td><code>sinks</code> <em><a href="meta.md#Object">metav1.Object</a></em></td>
<td><p>The output objects for the pipeline.</p></td>
</tr>
</tbody>
</table>

JobState (`string` alias)

JobState represents the state of a pipeline job.

### SplitTransform

SplitTransform is the Schema for the splittransforms API

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
<td><code>metadata</code> <em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/#objectmeta-v1-meta">Kubernetes metav1.ObjectMeta</a></em></td>
<td>Refer to the Kubernetes API documentation for the fields of the <code>metadata</code> field.</td>
</tr>
<tr class="even">
<td><code>spec</code> <em><a href="#SplitTransformSpec">SplitTransformSpec</a></em></td>
<td><br />
<br />

<table>
<tbody>
<tr class="odd">
<td><code>globals</code> <em><a href="meta.md#SourceSinkConfig">metav1.SourceSinkConfig</a></em></td>
<td><p>Global configurations to apply when omitted from the src or sink configurations.</p></td>
</tr>
<tr class="even">
<td><code>src</code> <em><a href="meta.md#SourceSinkConfig">metav1.SourceSinkConfig</a></em></td>
<td><p>Configurations for src object to the pipeline.</p></td>
</tr>
<tr class="odd">
<td><code>video</code> <em><a href="meta.md#SourceSinkConfig">metav1.SourceSinkConfig</a></em></td>
<td><p>Configurations for video stream outputs. The linkto field in the pipeline config should be present with the value <code>video-out</code> to direct an element to this output.</p></td>
</tr>
<tr class="even">
<td><code>audio</code> <em><a href="meta.md#SourceSinkConfig">metav1.SourceSinkConfig</a></em></td>
<td><p>Configurations for audio stream outputs. The linkto field in the pipeline config should be present with the value <code>audio-out</code> to direct an element to this output.</p></td>
</tr>
<tr class="odd">
<td><code>pipeline</code> <em><a href="meta.md#PipelineConfig">metav1.PipelineConfig</a></em></td>
<td><p>The configuration for the processing pipeline</p></td>
</tr>
</tbody>
</table></td>
</tr>
<tr class="odd">
<td><code>status</code> <em><a href="#SplitTransformStatus">SplitTransformStatus</a></em></td>
<td></td>
</tr>
</tbody>
</table>

### SplitTransformSpec

(*Appears on:* [SplitTransform](#SplitTransform))

SplitTransformSpec defines the desired state of SplitTransform. Note
that due to current implementation, the various streams can be directed
to different buckets, but they have to be buckets accessible via the
same MinIO/S3 server(s).

<table>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>globals</code> <em><a href="meta.md#SourceSinkConfig">metav1.SourceSinkConfig</a></em></td>
<td><p>Global configurations to apply when omitted from the src or sink configurations.</p></td>
</tr>
<tr class="even">
<td><code>src</code> <em><a href="meta.md#SourceSinkConfig">metav1.SourceSinkConfig</a></em></td>
<td><p>Configurations for src object to the pipeline.</p></td>
</tr>
<tr class="odd">
<td><code>video</code> <em><a href="meta.md#SourceSinkConfig">metav1.SourceSinkConfig</a></em></td>
<td><p>Configurations for video stream outputs. The linkto field in the pipeline config should be present with the value <code>video-out</code> to direct an element to this output.</p></td>
</tr>
<tr class="even">
<td><code>audio</code> <em><a href="meta.md#SourceSinkConfig">metav1.SourceSinkConfig</a></em></td>
<td><p>Configurations for audio stream outputs. The linkto field in the pipeline config should be present with the value <code>audio-out</code> to direct an element to this output.</p></td>
</tr>
<tr class="odd">
<td><code>pipeline</code> <em><a href="meta.md#PipelineConfig">metav1.PipelineConfig</a></em></td>
<td><p>The configuration for the processing pipeline</p></td>
</tr>
</tbody>
</table>

### Transform

Transform is the Schema for the transforms API

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
<td><code>metadata</code> <em><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/#objectmeta-v1-meta">Kubernetes metav1.ObjectMeta</a></em></td>
<td>Refer to the Kubernetes API documentation for the fields of the <code>metadata</code> field.</td>
</tr>
<tr class="even">
<td><code>spec</code> <em><a href="#TransformSpec">TransformSpec</a></em></td>
<td><br />
<br />

<table>
<tbody>
<tr class="odd">
<td><code>globals</code> <em><a href="meta.md#SourceSinkConfig">metav1.SourceSinkConfig</a></em></td>
<td><p>Global configurations to apply when omitted from the src or sink configurations.</p></td>
</tr>
<tr class="even">
<td><code>src</code> <em><a href="meta.md#SourceSinkConfig">metav1.SourceSinkConfig</a></em></td>
<td><p>Configurations for src object to the pipeline.</p></td>
</tr>
<tr class="odd">
<td><code>sink</code> <em><a href="meta.md#SourceSinkConfig">metav1.SourceSinkConfig</a></em></td>
<td><p>Configurations for sink objects from the pipeline.</p></td>
</tr>
<tr class="even">
<td><code>pipeline</code> <em><a href="meta.md#PipelineConfig">metav1.PipelineConfig</a></em></td>
<td><p>The configuration for the processing pipeline</p></td>
</tr>
</tbody>
</table></td>
</tr>
<tr class="odd">
<td><code>status</code> <em><a href="#TransformStatus">TransformStatus</a></em></td>
<td></td>
</tr>
</tbody>
</table>

### TransformSpec

(*Appears on:* [Transform](#Transform))

TransformSpec defines the desired state of Transform

<table>
<thead>
<tr class="header">
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr class="odd">
<td><code>globals</code> <em><a href="meta.md#SourceSinkConfig">metav1.SourceSinkConfig</a></em></td>
<td><p>Global configurations to apply when omitted from the src or sink configurations.</p></td>
</tr>
<tr class="even">
<td><code>src</code> <em><a href="meta.md#SourceSinkConfig">metav1.SourceSinkConfig</a></em></td>
<td><p>Configurations for src object to the pipeline.</p></td>
</tr>
<tr class="odd">
<td><code>sink</code> <em><a href="meta.md#SourceSinkConfig">metav1.SourceSinkConfig</a></em></td>
<td><p>Configurations for sink objects from the pipeline.</p></td>
</tr>
<tr class="even">
<td><code>pipeline</code> <em><a href="meta.md#PipelineConfig">metav1.PipelineConfig</a></em></td>
<td><p>The configuration for the processing pipeline</p></td>
</tr>
</tbody>
</table>

------------------------------------------------------------------------

*Generated with `gen-crd-api-reference-docs` on git commit `d333833`.*
