{{define "list_metrics_info"}}
指标名 | 说明 | 单位
--- | --- | ---
{{- $idx := 1 -}}
{{- range .}}
{{.Name}} |
    {{- if eq .Comm "" -}}
        {{- printf " %s " .Desc -}}
    {{- else -}}
        {{- $cite := printf "<sup><a href=\"#metric-comments-%s\">%d</a></sup>" .Name $idx -}}
        {{- printf " %s " (replace "^" $cite .Desc) -}}
        {{- $idx = add1 $idx -}}
    {{- end -}}
| {{.Unit}}
{{- end -}}
{{- /*备注*/ -}}
{{- $idx = 1}}
{{range .}}
    {{- if ne .Comm "" -}}
> <sup><a id="metric-comments-{{.Name}}">{{$idx}}</a></sup> {{printf "%s  \n" .Comm}}
    {{- $idx = add1 $idx -}}
    {{- end -}}
{{- end -}}
{{end}}

{{define "list_metrics_expr"}}
{{- range .}}
    {{- if ne .Expr "" }}
{{.Name}} | {{replace "|" "\\|" .Expr}}
    {{- end -}}
{{- end -}}
{{end}}

## 监控指标和API文档

KubeSphere资源监控系统基于Prometheus开源监控平台，为KubeSphere平台提供了全面的、多维度的监控指标信息，主要涵盖了平台与集群状态、物理资源、应用资源、存储以及KubeSphere核心组件等多方面的监控信息，具体分为以下几个层级：

- KubeSphere(平台)  
- Cluster(集群)  
- Workspace(企业空间)  
- Namespace(项目)  
- Workload(工作负载)  
- Pod(容器组)  
- Container(容器)  
- Component(KubeSphere核心组件)  

下文涉及监控API部分基于以下API组和版本：

- Group: `monitoring.kubesphere.io`  
- Version: `v1alpha3`  

> 参考[这里](https://v3-0.docs.kubesphere.io/zh/docs/api-reference/api-docs/#api-%E5%8F%82%E8%80%83)获取详细的API文档

另外，下文列表中所列述指标名称，并非原始Prometheus指标名称，其中大部分指标与Prometheus指标的具体对应关系请参考<a href="#appendix-1">附录一</a>。

### KubeSphere

平台级指标表征了平台相关资源的状态，在API被访问时它们经过实时统计后返回。

#### 指标

{{- template "list_metrics_info" (index . "kubesphere")}}

#### API

Path | 说明 | 主要查询参数  
--- | --- | ---  
`/kubesphere` | 获取平台级指标  

### Cluster

集群级指标表征了整个集群各种资源的状态，它们由Prometheus的采集指标和评估[Recording Rule](https://prometheus.io/docs/prometheus/latest/configuration/recording_rules/)的派生指标转换而来。

> - 参考[这里](https://v3-0.docs.kubesphere.io/zh/docs/cluster-administration/cluster-status-monitoring/)可以在Console上查看集群的资源使用情况，包括各种物理资源的用量信息。
> - 参考[这里](https://v3-0.docs.kubesphere.io/zh/docs/cluster-administration/application-resources-monitoring/)可以在Console上以集群视角查看各种应用资源的使用情况。

#### 指标

{{- template "list_metrics_info" (index . "cluster")}}

#### API

Path | 说明 | 主要查询参数  
--- | --- | ---  
`/cluster` | 获取集群级指标 | `metrics_filter`

> `metrics_filter`是`|`分隔的多个指标(cluster指标)，例如`cluster_cpu_usage|cluster_disk_size_usage`。


### Node

节点级指标表征了节点资源的状态，它们由Prometheus的采集指标和评估[Recording Rule](https://prometheus.io/docs/prometheus/latest/configuration/recording_rules/)的派生指标转换而来。

> - 参考[这里](https://v3-0.docs.kubesphere.io/zh/docs/cluster-administration/cluster-status-monitoring/#%E9%9B%86%E7%BE%A4%E8%8A%82%E7%82%B9%E7%8A%B6%E6%80%81)可以在Console上查看节点的资源使用情况。

#### 指标

{{- template "list_metrics_info" (index . "node")}}

#### API

Path | 说明 | 主要查询参数  
--- | --- | ---
`/nodes` | 获取节点指标 | `metrics_filter`,`resources_filter`  
`/nodes/{node}` | 获取某个节点的指标 | `metrics_filter`  

> - `{node}`是节点名称。  
> - `metrics_filter`是`|`分隔的多个指标名称(节点指标)。  
> - `resources_filter`是`|`分隔的多个节点名称。  


### Workspace

企业空间级指标表征了企业空间中的资源状态，它们由Prometheus的采集指标和评估[Recording Rule](https://prometheus.io/docs/prometheus/latest/configuration/recording_rules/)的派生指标转换而来。  

> - Console上企业空间的概览页展示了该企业空间的资源状态信息。

#### 指标

- 以下指标由Prometheus指标转换而来:

{{ template "list_metrics_info" (index . "workspace")}}

- 以下指标在API被访问时经过实时统计后返回:

{{ template "list_metrics_info" (index . "workspace_stat")}}

#### API

Path | 说明 | 主要查询参数  
--- | --- | ---  
`/workspaces` | 获取工作空间指标 | `metrics_filter`,`resources_filter`
`/workspaces/{workspace}` | 获取某个工作空间的指标 | `metrics_filter`,`type`  

> - `{workspace}`是企业空间名称。
> - `metrics_filter`是`|`分隔的多个指标名称(企业空间指标) 。
> - `resources_filter`是`|`分隔的多个企业空间名称。
> - `type`指定为`statistics`时返回实时统计的那些指标。


### Namespace

项目级指标表征了项目中的资源状态，它们由Prometheus的采集指标和评估[Recording Rule](https://prometheus.io/docs/prometheus/latest/configuration/recording_rules/)的派生指标转换而来。

> - Console上项目的概览页展示了该项目的资源状态信息。  

#### 指标

{{- template "list_metrics_info" (index . "namespace")}}

#### API

Path | 说明 | 主要查询参数  
--- | --- | ---  
`/namespaces` | 获取项目指标 | `metrics_filter`,`resources_filter`
`/namespaces/{namespace}` | 获取某个项目的指标 | `metrics_filter`  

> `{namespace}`是项目名称。
> `metrics_filter`是`|`分隔的多个指标名称(项目指标)。
> `resources_filter`是`|`分隔的多个项目名称。

### Workload

工作负载级指标表征了部署、有状态副本集、守护进程集等负载的资源状态信息，它们由Prometheus的采集指标和评估[Recording Rule](https://prometheus.io/docs/prometheus/latest/configuration/recording_rules/)的派生指标转换而来。  

> - 参考[这里](https://v3-0.docs.kubesphere.io/zh/docs/project-user-guide/application-workloads/deployments/#%E7%9B%91%E6%8E%A7)可以在Console上查看某个部署的资源状态。
> - 参考[这里](https://v3-0.docs.kubesphere.io/zh/docs/project-user-guide/application-workloads/statefulsets/#%E7%9B%91%E6%8E%A7)可以在Console上查看某个有状态副本集的资源状态。
> - 参考[这里](https://v3-0.docs.kubesphere.io/zh/docs/project-user-guide/application-workloads/daemonsets/#%E7%9B%91%E6%8E%A7)可以在Console上查看某个守护进程集的资源状态。

#### 指标

{{- template "list_metrics_info" (index . "workLoad")}}

#### API

Path | 说明 | 主要查询参数  
--- | --- | ---  
`/namespaces/{namespace}/workloads` | 获取某个项目下的工作负载指标 | `metrics_filter`,`resources_filter`  
`/namespaces/{namespace}/workloads/{kind}` | 获取某个项目下某一类别工作负载的指标 | `metrics_filter`  

> - `{namespace}`是pod所属项目的名称。
> - `{kind}`是工作负载类别，`deployment`,`statefulset`,`daemonset`之一。
> - `metrics_filter`是`|`分隔的多个指标名称(工作负载指标)。
> - `resources_filter`是`|`分隔的多个工作负载名称。

### Pod

容器组级指标表征了容器组的资源状态信息，它们主要由Prometheus的采集指标转换而来。

#### 指标

{{- template "list_metrics_info" (index . "pod")}}

#### API

Path | 说明 | 主要查询参数  
--- | --- | ---  
`/namespaces/{namespace}/pods` | 获取某个项目下的pod指标 | `metrics_filter`,`resources_filter`  
`/namespaces/{namespace}/pods/{pod}` | 获取某个项目下某个pod的指标 | `metrics_filter`  
`/namespaces/{namespace}/workloads/{kind}/{workload}/pods` | 获取某个项目下某个工作负载的pod指标 | `metrics_filter`,`resources_filter`
`/nodes/{node}/pods` | 获取某个节点下的pod指标 | `metrics_filter`,`resources_filter`

> - `{namespace}`是pod所属项目的名称。
> - `{pod}`是pod名称。
> - `{kind}`是工作负载类别，`deployment`,`statefulset`,`daemonset`之一。
> - `{workload}`是工作负载名称。
> - `{node}`是节点名称。
> - `metrics_filter`是`|`分隔的多个指标名称(pod指标)。
> - `resources_filter`是`|`分隔的多个pod名称。

### Container

容器级指标表征了容器的资源状态信息，它们主要由Prometheus的采集指标转换而来。

#### 指标

{{- template "list_metrics_info" (index . "container")}}

#### API

Path | 说明 | 主要查询参数  
--- | --- | ---  
`/namespaces/{namespace}/pods/{pod}/containers` | 获取某个项目下某个pod的容器指标 | `metrics_filter`,`resources_filter`
`/namespaces/{namespace}/pods/{pod}/containers/{container}` | 获取某个项目下某个pod里的某个容器的指标 | `metrics_filter`

> - `{namespace}`是pod所属项目的名称
> - `{pod}`是pod名称
> - `{container}`是容器名称
> - `{kind}`是工作负载类别，`deployment`,`statefulset`,`daemonset`之一
> - `{workload}`是工作负载名称
> - `{node}`是节点名称
> - `metrics_filter`是`|`分隔的多个指标名称(容器指标)
> - `resources_filter`是`|`分隔的多个容器名称

### PVC

PVC指标表征了PVC的inode和容量使用情况，它们主要由Prometheus的采集指标转换而来。

> - 参考[这里](https://v3-0.docs.kubesphere.io/zh/docs/project-user-guide/storage/volumes/#%E7%9B%91%E6%8E%A7%E5%AD%98%E5%82%A8%E5%8D%B7)可以在Console上查看存储卷的状态信息。

#### 指标

{{- template "list_metrics_info" (index . "pvc")}}

#### API

Path | 说明 | 主要查询参数  
--- | --- | ---  
`/storageclasses/{storageclass}/persistentvolumeclaims` | 获取某个存储类型的pvc指标 | `metrics_filter`,`resources_filter`
`/namespaces/{namespace}/persistentvolumeclaims` | 获取某个项目下的pvc指标 | `metrics_filter`,`resources_filter`
`/namespaces/{namespace}/persistentvolumeclaims/{pvc}` | 获取某个项目下某个pvc的指标 | `metrics_filter`  

> - `{storageclass}`是存储类型的名称。
> - `{namespace}`是pvc所属项目的名称。
> - `{pvc}`是pvc名称。
> - `metrics_filter`是`|`分隔的多个指标名称(pvc指标)。
> - `resources_filter`是`|`分隔的多个pvc名称。

### Component

Component指标表征了一些核心组件的状态信息，它们由Prometheus的采集指标和评估[Recording Rule](https://prometheus.io/docs/prometheus/latest/configuration/recording_rules/)的派生指标转换而来。  

> 参考[这里](https://v3-0.docs.kubesphere.io/zh/docs/cluster-administration/cluster-status-monitoring/#etcd-%E7%9B%91%E6%8E%A7)可以在Console上查看etcd组件的状态信息。
> 参考[这里](https://v3-0.docs.kubesphere.io/zh/docs/cluster-administration/cluster-status-monitoring/#apiserver-%E7%9B%91%E6%8E%A7)可以在Console上查看kube-apiserver组件的状态信息。
> 参考[这里](https://v3-0.docs.kubesphere.io/zh/docs/cluster-administration/cluster-status-monitoring/#%E8%B0%83%E5%BA%A6%E5%99%A8%E7%9B%91%E6%8E%A7)可以在Console上查看scheduler组件的状态信息。

#### 指标

{{- template "list_metrics_info" (index . "component")}}

#### API

Path | 说明 | 主要查询参数  
--- | --- | ---  
`/components/{component}` | 获取某个组件的指标 | `metrics_filter`  

> - `{component}`是`etcd`,`apiserver`,`scheduler`之一。
> - `metrics_filter`是`|`分隔的多个指标名称(组件指标)。


## 附录一<a id="appendix-1"></a>: KubeSphere监控指标与Prometheus指标对照表

指标名 | PromQL模板
--- | ---
{{- template "list_metrics_expr" (index . "cluster")}}
{{- template "list_metrics_expr" (index . "node")}}
{{- template "list_metrics_expr" (index . "workSpace")}}
{{- template "list_metrics_expr" (index . "namespace")}}
{{- template "list_metrics_expr" (index . "workLoad")}}
{{- template "list_metrics_expr" (index . "pod")}}
{{- template "list_metrics_expr" (index . "container")}}
{{- template "list_metrics_expr" (index . "pvc")}}
{{- template "list_metrics_expr" (index . "component")}}

> - 第一列是在监控API调用时可用的指标名称，Prometheus中并未存储对应指标名称的时序数据。
> - 监控API调用时，将转换API参数为Prometheus标签选择器，代入到PromQL模板的`$1`,`$2`，然后使用生成的PromQL请求Prometheus服务。
> - Prometheus指标中，对于使用`:`分隔符命名的指标，通常是[Recording Rule](https://prometheus.io/docs/prometheus/latest/configuration/recording_rules/)派生的指标。

