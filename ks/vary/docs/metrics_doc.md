



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
指标名 | 说明 | 单位
--- | --- | ---
kubesphere_workspace_count | 平台工作空间总数 | 
kubesphere_user_count | 平台用户总数 | 
kubesphere_cluser_count | 平台集群总数 | 
kubesphere_app_template_count | 平台应用模板总数 | 


#### API

Path | 说明 | 主要查询参数  
--- | --- | ---  
`/kubesphere` | 获取平台级指标  

### Cluster

集群级指标表征了整个集群各种资源的状态，它们由Prometheus的采集指标和评估[Recording Rule](https://prometheus.io/docs/prometheus/latest/configuration/recording_rules/)的派生指标转换而来。

> - 参考[这里](https://v3-0.docs.kubesphere.io/zh/docs/cluster-administration/cluster-status-monitoring/)可以在Console上查看集群的资源使用情况，包括各种物理资源的用量信息。
> - 参考[这里](https://v3-0.docs.kubesphere.io/zh/docs/cluster-administration/application-resources-monitoring/)可以在Console上以集群视角查看各种应用资源的使用情况。

#### 指标
指标名 | 说明 | 单位
--- | --- | ---
cluster_cpu_utilisation | 集群CPU使用率 | 
cluster_cpu_usage | 集群CPU用量 | Core
cluster_cpu_total | 集群CPU总量 | Core
cluster_load1 | 集群1分钟CPU平均负载<sup><a href="#metric-comments-cluster_load1">1</a></sup> | 
cluster_load5 | 集群5分钟CPU平均负载 | 
cluster_load15 | 集群15分钟CPU平均负载 | 
cluster_memory_utilisation | 集群内存使用率 | 
cluster_memory_available | 集群可用内存 | Byte
cluster_memory_total | 集群内存总量 | Byte
cluster_memory_usage_wo_cache | 集群内存使用量<sup><a href="#metric-comments-cluster_memory_usage_wo_cache">2</a></sup> | Byte
cluster_net_utilisation | 集群网络数据传输速率 | Byte/s
cluster_net_bytes_transmitted | 集群网络数据发送速率 | Byte/s
cluster_net_bytes_received | 集群网络数据接受速率 | Byte/s
cluster_disk_read_iops | 集群磁盘每秒读次数 | 次/s
cluster_disk_write_iops | 集群磁盘每秒写次数 | 次/s
cluster_disk_read_throughput | 集群磁盘每秒读取数据量 | Byte/s
cluster_disk_write_throughput | 集群磁盘每秒写入数据量 | Byte/s
cluster_disk_size_usage | 集群磁盘使用量 | Byte
cluster_disk_size_utilisation | 集群磁盘使用率 | 
cluster_disk_size_capacity | 集群磁盘总容量 | Byte
cluster_disk_size_available | 集群磁盘可用大小 | Byte
cluster_disk_inode_total | 集群inode总数 | 
cluster_disk_inode_usage | 集群inode已使用数 | 
cluster_disk_inode_utilisation | 集群inode使用率 | 
cluster_node_online | 集群节点在线数 | 
cluster_node_offline | 集群节点下线数 | 
cluster_node_offline_ratio | 集群节点下线比例 | 
cluster_node_total | 集群节点总数 | 
cluster_pod_count | 集群中调度完成<sup><a href="#metric-comments-cluster_pod_count">3</a></sup>Pod数量 | 
cluster_pod_quota | 集群各节点Pod最大容纳量<sup><a href="#metric-comments-cluster_pod_quota">4</a></sup>总和 | 
cluster_pod_utilisation | 集群Pod最大容纳量使用率 | 
cluster_pod_running_count | 集群中处于Running阶段<sup><a href="#metric-comments-cluster_pod_running_count">5</a></sup>的Pod数量 | 
cluster_pod_succeeded_count | 集群中处于Succeeded阶段的Pod数量 | 
cluster_pod_abnormal_count | 集群中异常Pod<sup><a href="#metric-comments-cluster_pod_abnormal_count">6</a></sup>数量 | 
cluster_pod_abnormal_ratio | 集群中异常Pod比例 <sup><a href="#metric-comments-cluster_pod_abnormal_ratio">7</a></sup> | 
cluster_ingresses_extensions_count | 集群Ingress数 | 
cluster_cronjob_count | 集群CronJob数 | 
cluster_pvc_count | 集群PersistentVolumeClaim数 | 
cluster_daemonset_count | 集群DaemonSet数 | 
cluster_deployment_count | 集群Deployment数 | 
cluster_endpoint_count | 集群Endpoint数 | 
cluster_hpa_count | 集群Horizontal Pod Autoscaler数 | 
cluster_job_count | 集群Job数 | 
cluster_statefulset_count | 集群StatefulSet数 | 
cluster_replicaset_count | 集群ReplicaSet数 | 
cluster_service_count | 集群Service数 | 
cluster_secret_count | 集群Secret数 | 
cluster_namespace_count | 集群Namespace数 | 
> <sup><a id="metric-comments-cluster_load1">1</a></sup> 指单位时间内，单位CPU运行队列中处于可运行或不可中断状态的平均进程数。如果数值大于1，表示CPU不足以服务进程，有进程在等待。  
> <sup><a id="metric-comments-cluster_memory_usage_wo_cache">2</a></sup> 不包含buffer, cache。  
> <sup><a id="metric-comments-cluster_pod_count">3</a></sup> Pod已经被调度到节点上，即status.conditions.PodScheduled=true 。参考：https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#pod-conditions  
> <sup><a id="metric-comments-cluster_pod_quota">4</a></sup> 节点Pod最大容纳量一般默认110个Pod。参考：https://kubernetes.io/docs/reference/command-line-tools-reference/kubelet/#options  
> <sup><a id="metric-comments-cluster_pod_running_count">5</a></sup> Running阶段表示该Pod已经绑定到了一个节点上, Pod中所有的容器都已被创建。至少有一个容器正在运行，或者正处于启动或重启状态。参考：https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#pod-phase  
> <sup><a id="metric-comments-cluster_pod_abnormal_count">6</a></sup> 异常Pod：如果一个Pod的status.conditions.ContainersReady字段值为false，说明该Pod不可用。我们在判定Pod是否异常时，还需要考虑到Pod可能正处于ContainerCreating状态或者Succeeded已完成阶段。综合以上情况，异常Pod总数的算法可表示为：Abnormal Pods = Total Pods - ContainersReady Pods - ContainerCreating Pods - Succeeded Pods 。  
> <sup><a id="metric-comments-cluster_pod_abnormal_ratio">7</a></sup> 异常Pod比例：异常Pod数 / 非 Succeeded Pod 数。  


#### API

Path | 说明 | 主要查询参数  
--- | --- | ---  
`/cluster` | 获取集群级指标 | `metrics_filter`

> `metrics_filter`是`|`分隔的多个指标(cluster指标)，例如`cluster_cpu_usage|cluster_disk_size_usage`。


### Node

节点级指标表征了节点资源的状态，它们由Prometheus的采集指标和评估[Recording Rule](https://prometheus.io/docs/prometheus/latest/configuration/recording_rules/)的派生指标转换而来。

> - 参考[这里](https://v3-0.docs.kubesphere.io/zh/docs/cluster-administration/cluster-status-monitoring/#%E9%9B%86%E7%BE%A4%E8%8A%82%E7%82%B9%E7%8A%B6%E6%80%81)可以在Console上查看节点的资源使用情况。

#### 指标
指标名 | 说明 | 单位
--- | --- | ---
node_cpu_utilisation | 节点CPU使用率 | 
node_cpu_total | 节点CPU总量 | Core
node_cpu_usage | 节点CPU用量 | Core
node_load1 | 节点1分钟CPU平均负载 | 
node_load5 | 节点5分钟CPU平均负载 | 
node_load15 | 节点15分钟CPU平均负载 | 
node_memory_utilisation | 节点内存使用率 | 
node_memory_usage_wo_cache | 节点内存使用量<sup><a href="#metric-comments-node_memory_usage_wo_cache">1</a></sup> | Byte
node_memory_available | 节点可用内存 | Byte
node_memory_total | 节点内存总量 | Byte
node_net_utilisation | 节点网络数据传输速率 | Byte/s
node_net_bytes_transmitted | 节点网络数据发送速率 | Byte/s
node_net_bytes_received | 节点网络数据接受速率 | Byte/s
node_disk_read_iops | 节点磁盘每秒读次数 | 次/s
node_disk_write_iops | 节点磁盘每秒写次数 | 次/s
node_disk_read_throughput | 节点磁盘每秒读取数据量 | Byte/s
node_disk_write_throughput | 节点磁盘每秒写入数据量 | Byte/s
node_disk_size_capacity | 节点磁盘总容量 | Byte
node_disk_size_available | 节点磁盘可用大小 | Byte
node_disk_size_usage | 节点磁盘使用量 | Byte
node_disk_size_utilisation | 节点磁盘使用率 | 
node_disk_inode_total | 节点inode总数 | 
node_disk_inode_usage | 节点inode已使用数 | 
node_disk_inode_utilisation | 节点inode使用率 | 
node_pod_count | 节点调度完成Pod数量 | 
node_pod_quota | 节点Pod最大容纳量 | 
node_pod_utilisation | 节点Pod最大容纳量使用率 | 
node_pod_running_count | 节点中处于Running阶段的Pod数量 | 
node_pod_succeeded_count | 节点中处于Succeeded阶段的Pod数量 | 
node_pod_abnormal_count | 节点异常Pod数量 | 
node_pod_abnormal_ratio | 节点异常Pod比例 | 
> <sup><a id="metric-comments-node_memory_usage_wo_cache">1</a></sup> 不包含buffer、 cache  


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


指标名 | 说明 | 单位
--- | --- | ---
workspace_cpu_usage | 企业空间CPU用量 | Core
workspace_memory_usage | 企业空间内存使用量（包含缓存） | Byte
workspace_memory_usage_wo_cache | 企业空间内存使用量 | Byte
workspace_net_bytes_transmitted | 企业空间网络数据发送速率 | Byte/s
workspace_net_bytes_received | 企业空间网络数据接受速率 | Byte/s
workspace_pod_count | 企业空间内非终止阶段Pod数量 | 
workspace_pod_running_count | 企业空间内处于Running 阶段的Pod数量 | 
workspace_pod_succeeded_count | 企业空间内处于Succeeded阶段的Pod数量 | 
workspace_pod_abnormal_count | 企业空间异常Pod数量 | 
workspace_pod_abnormal_ratio | 企业空间异常Pod比例 | 
workspace_ingresses_extensions_count | 企业空间Ingress 数 | 
workspace_cronjob_count | 企业空间CronJob数 | 
workspace_pvc_count | 企业空间PersistentVolumeClaim数 | 
workspace_daemonset_count | 企业空间DaemonSet数 | 
workspace_deployment_count | 企业空间Deployment数 | 
workspace_endpoint_count | 企业空间Endpoint数 | 
workspace_hpa_count | 企业空间Horizontal Pod Autoscaler数 | 
workspace_job_count | 企业空间Job数 | 
workspace_statefulset_count | 企业空间StatefulSet数 | 
workspace_replicaset_count | 企业空间ReplicaSet数 | 
workspace_service_count | 企业空间Service数 | 
workspace_secret_count | 企业空间Secret数 | 
> <sup><a id="metric-comments-workspace_pod_count">1</a></sup> 非终止阶段的Pod指处于Pending、Running、Unknown阶段的 Pod，不包含被成功终止，或者因非 0 状态退出被系统终止的 Pod。参考：https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#pod-conditions  


- 以下指标在API被访问时经过实时统计后返回:


指标名 | 说明 | 单位
--- | --- | ---
workspace_namespace_count | 企业空间项目总数 | 
workspace_devops_project_count | 企业空间DevOps工程总数 | 
workspace_member_count | 企业空间成员数 | 
workspace_role_count | 企业空间角色数 | 


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
指标名 | 说明 | 单位
--- | --- | ---
namespace_cpu_usage | 项目CPU用量 | Core
namespace_memory_usage | 项目内存使用量（包含缓存） | Byte
namespace_memory_usage_wo_cache | 项目内存使用量 | Byte
namespace_net_bytes_transmitted | 项目网络数据发送速率 | Byte/s
namespace_net_bytes_received | 项目网络数据接受速率 | Byte/s
namespace_pod_count | 项目内非终止阶段Pod数量 | 
namespace_pod_running_count | 项目内处于Running阶段的Pod数量 | 
namespace_pod_succeeded_count | 项目内处于Succeeded阶段的Pod数量 | 
namespace_pod_abnormal_count | 项目异常Pod数量 | 
namespace_pod_abnormal_ratio | 项目异常Pod比例 | 
namespace_cronjob_count | 项目CronJob数 | 
namespace_pvc_count | 项目PersistentVolumeClaim数 | 
namespace_daemonset_count | 项目DaemonSet数 | 
namespace_deployment_count | 项目Deployment数 | 
namespace_endpoint_count | 项目Endpoint数 | 
namespace_hpa_count | 项目Horizontal Pod Autoscaler数 | 
namespace_job_count | 项目Job数 | 
namespace_statefulset_count | 项目StatefulSet数 | 
namespace_replicaset_count | 项目ReplicaSet数 | 
namespace_service_count | 项目Service数 | 
namespace_secret_count | 项目Secret数 | 
namespace_ingresses_extensions_count | 项目Ingress数 | 


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
指标名 | 说明 | 单位
--- | --- | ---


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
指标名 | 说明 | 单位
--- | --- | ---
pod_cpu_usage | 容器组CPU用量 | Core
pod_memory_usage | 容器组内存使用量（包含缓存） | Byte
pod_memory_usage_wo_cache | 容器组内存使用量 | Byte
pod_net_bytes_transmitted | 容器组网络数据发送速率 | Byte/s
pod_net_bytes_received | 容器组网络数据接受速率 | Byte/s


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
指标名 | 说明 | 单位
--- | --- | ---
container_cpu_usage | 容器CPU用量 | Core
container_memory_usage | 容器内存使用量（包含缓存） | Byte
container_memory_usage_wo_cache | 容器内存使用量 | Byte


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
指标名 | 说明 | 单位
--- | --- | ---
pvc_inodes_available | pvc的inode剩余数 | 
pvc_inodes_used | pvc的inode已使用数 | 
pvc_inodes_total | pvc的inode总数 | 
pvc_inodes_utilisation | pvc的inode利用率 | 
pvc_bytes_available | pvc剩余容量 | Byte
pvc_bytes_used | pvc使用量 | Byte
pvc_bytes_total | pvc总容量 | Byte
pvc_bytes_utilisation | pvc利用率 | 


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
指标名 | 说明 | 单位
--- | --- | ---
etcd_server_list | etcd集群节点列表<sup><a href="#metric-comments-etcd_server_list">1</a></sup> | 
etcd_server_total | etcd集群节点总数 | 
etcd_server_up_total | etcd集群在线节点数 | 
etcd_server_has_leader | etcd集群各节点是否有leader<sup><a href="#metric-comments-etcd_server_has_leader">2</a></sup> | 
etcd_server_leader_changes | etcd集群各节点观察到leader变化数（ 1h 内） | 
etcd_server_proposals_failed_rate | etcd集群各节点提案失败<sup><a href="#metric-comments-etcd_server_proposals_failed_rate">3</a></sup>频率平均数 | 次/s
etcd_server_proposals_applied_rate | etcd集群各节点提案应用频率平均数 | 次/s
etcd_server_proposals_committed_rate | etcd集群各节提案提交频率平均数 | 次/s
etcd_server_proposals_pending_count | etcd集群各节点排队提案数平均值 | 
etcd_mvcc_db_size | etcd 集群各节点数据库大小平均值 | Byte
etcd_network_client_grpc_received_bytes | etcd集群向gRPC客户端发送数据速率 | Byte/s
etcd_network_client_grpc_sent_bytes | etcd集群接受gRPC客户端数据速率 | Byte/s
etcd_grpc_call_rate | etcd集群gRPC请求速率 | 次/s
etcd_grpc_call_failed_rate | etcd集群gRPC请求失败速率 | 次/s
etcd_grpc_server_msg_received_rate | etcd集群gRPC流式消息接收速率 | 次/s
etcd_grpc_server_msg_sent_rate | etcd集群gRPC流式消息发送速率 | 次/s
etcd_disk_wal_fsync_duration | etcd集群各节点WAL日志同步时间平均值 | 秒
etcd_disk_wal_fsync_duration_quantile | etcd集群WAL日志同步时间平均值（按分位数统计）<sup><a href="#metric-comments-etcd_disk_wal_fsync_duration_quantile">4</a></sup> | 秒
etcd_disk_backend_commit_duration | etcd集群各节点库同步时间<sup><a href="#metric-comments-etcd_disk_backend_commit_duration">5</a></sup>平均值 | 秒
etcd_disk_backend_commit_duration_quantile | etcd集群各节点库同步时间平均值（按分位数统计） | 秒
apiserver_up_sum | APIServer<sup><a href="#metric-comments-apiserver_up_sum">6</a></sup>在线实例数 | 
apiserver_request_rate | APIServer每秒接受请求数 | 
apiserver_request_by_verb_rate | APIServer每秒接受请求数（按HTTP请求方法分类统计） | 
apiserver_request_latencies | APIServer请求平均迟延 | 秒
apiserver_request_by_verb_latencies | APIServer请求平均迟延（按HTTP请求方法分类统计） | 秒
scheduler_up_sum | 调度器<sup><a href="#metric-comments-scheduler_up_sum">7</a></sup>在线实例数 | 
scheduler_schedule_attempts | 调度器累计调度次数<sup><a href="#metric-comments-scheduler_schedule_attempts">8</a></sup> | 
scheduler_schedule_attempt_rate | 调度器调度频率 | 次/s
scheduler_e2e_scheduling_latency | 调度器调度延迟 | 秒
scheduler_e2e_scheduling_latency_quantile | 调度器调度延迟（按分位数统计） | 秒
controller_manager_up_sum | Controller Manager<sup><a href="#metric-comments-controller_manager_up_sum">9</a></sup>在线实例数 | 
> <sup><a id="metric-comments-etcd_server_list">1</a></sup> 如果某一节点返回值为 1 说明该etcd节点在线，0 说明节点下线。  
> <sup><a id="metric-comments-etcd_server_has_leader">2</a></sup> 如果某一节点返回值为0说明该节点没有leader，即该节点不可使用；如果集群中，所有节点都没有任何leader，则整个集群不可用。  
> <sup><a id="metric-comments-etcd_server_proposals_failed_rate">3</a></sup> 中英文对照说明：提案（consensus proposals）,失败提案（failed proposals），已提交提案（commited proposals），应用提案（applied proposals），排队提案（pending proposals）。  
> <sup><a id="metric-comments-etcd_disk_wal_fsync_duration_quantile">4</a></sup> 支持三种分位数统计：99th 百分位数、90th 百分位数、中位数。  
> <sup><a id="metric-comments-etcd_disk_backend_commit_duration">5</a></sup> 反映磁盘 I/O 延迟。如果数值过高，通常表示磁盘问题。  
> <sup><a id="metric-comments-apiserver_up_sum">6</a></sup> 指kube-apiserver。  
> <sup><a id="metric-comments-scheduler_up_sum">7</a></sup> 指kube-scheduler。  
> <sup><a id="metric-comments-scheduler_schedule_attempts">8</a></sup> 按调度结果分类统计：error（因调度器异常而无法调度的Pod数量）, scheduled（成功被调度的Pod数量）, unschedulable（无法被调度的Pod数量）。  
> <sup><a id="metric-comments-controller_manager_up_sum">9</a></sup> 指kube-controller-manager  


#### API

Path | 说明 | 主要查询参数  
--- | --- | ---  
`/components/{component}` | 获取某个组件的指标 | `metrics_filter`  

> - `{component}`是`etcd`,`apiserver`,`scheduler`之一。
> - `metrics_filter`是`|`分隔的多个指标名称(组件指标)。


## 附录一<a id="appendix-1"></a>: KubeSphere监控指标与Prometheus指标对照表

指标名 | PromQL模板
--- | ---
cluster_cpu_utilisation | :node_cpu_utilisation:avg1m
cluster_cpu_usage | round(:node_cpu_utilisation:avg1m * sum(node:node_num_cpu:sum), 0.001)
cluster_cpu_total | sum(node:node_num_cpu:sum)
cluster_load1 | sum(node_load1{job="node-exporter"}) / sum(node:node_num_cpu:sum)
cluster_load5 | sum(node_load5{job="node-exporter"}) / sum(node:node_num_cpu:sum)
cluster_load15 | sum(node_load15{job="node-exporter"}) / sum(node:node_num_cpu:sum)
cluster_memory_utilisation | :node_memory_utilisation:
cluster_memory_available | sum(node:node_memory_bytes_available:sum)
cluster_memory_total | sum(node:node_memory_bytes_total:sum)
cluster_memory_usage_wo_cache | sum(node:node_memory_bytes_total:sum) - sum(node:node_memory_bytes_available:sum)
cluster_net_utilisation | :node_net_utilisation:sum_irate
cluster_net_bytes_transmitted | sum(node:node_net_bytes_transmitted:sum_irate)
cluster_net_bytes_received | sum(node:node_net_bytes_received:sum_irate)
cluster_disk_read_iops | sum(node:data_volume_iops_reads:sum)
cluster_disk_write_iops | sum(node:data_volume_iops_writes:sum)
cluster_disk_read_throughput | sum(node:data_volume_throughput_bytes_read:sum)
cluster_disk_write_throughput | sum(node:data_volume_throughput_bytes_written:sum)
cluster_disk_size_usage | sum(max(node_filesystem_size_bytes{device=~"/dev/.*", device!~"/dev/loop\\d+", job="node-exporter"} - node_filesystem_avail_bytes{device=~"/dev/.*", device!~"/dev/loop\\d+", job="node-exporter"}) by (device, instance))
cluster_disk_size_utilisation | cluster:disk_utilization:ratio
cluster_disk_size_capacity | sum(max(node_filesystem_size_bytes{device=~"/dev/.*", device!~"/dev/loop\\d+", job="node-exporter"}) by (device, instance))
cluster_disk_size_available | sum(max(node_filesystem_avail_bytes{device=~"/dev/.*", device!~"/dev/loop\\d+", job="node-exporter"}) by (device, instance))
cluster_disk_inode_total | sum(node:node_inodes_total:)
cluster_disk_inode_usage | sum(node:node_inodes_total:) - sum(node:node_inodes_free:)
cluster_disk_inode_utilisation | cluster:disk_inode_utilization:ratio
cluster_node_online | sum(kube_node_status_condition{condition="Ready",status="true"})
cluster_node_offline | cluster:node_offline:sum
cluster_node_offline_ratio | cluster:node_offline:ratio
cluster_node_total | sum(kube_node_status_condition{condition="Ready"})
cluster_pod_count | cluster:pod:sum
cluster_pod_quota | sum(max(kube_node_status_capacity{resource="pods"}) by (node) unless on (node) (kube_node_status_condition{condition="Ready",status=~"unknown\|false"} > 0))
cluster_pod_utilisation | cluster:pod_utilization:ratio
cluster_pod_running_count | cluster:pod_running:count
cluster_pod_succeeded_count | count(kube_pod_info unless on (pod) (kube_pod_status_phase{phase=~"Failed\|Pending\|Unknown\|Running"} > 0) unless on (node) (kube_node_status_condition{condition="Ready",status=~"unknown\|false"} > 0))
cluster_pod_abnormal_count | cluster:pod_abnormal:sum
cluster_pod_abnormal_ratio | cluster:pod_abnormal:ratio
cluster_ingresses_extensions_count | sum(kube_ingress_labels)
cluster_cronjob_count | sum(kube_cronjob_labels)
cluster_pvc_count | sum(kube_persistentvolumeclaim_info)
cluster_daemonset_count | sum(kube_daemonset_labels)
cluster_deployment_count | sum(kube_deployment_labels)
cluster_endpoint_count | sum(kube_endpoint_labels)
cluster_hpa_count | sum(kube_hpa_labels)
cluster_job_count | sum(kube_job_labels)
cluster_statefulset_count | sum(kube_statefulset_labels)
cluster_replicaset_count | count(kube_replicaset_labels)
cluster_service_count | sum(kube_service_info)
cluster_secret_count | sum(kube_secret_info)
cluster_namespace_count | count(kube_namespace_labels)
node_cpu_utilisation | node:node_cpu_utilisation:avg1m{$1}
node_cpu_total | node:node_num_cpu:sum{$1}
node_cpu_usage | round(node:node_cpu_utilisation:avg1m{$1} * node:node_num_cpu:sum{$1}, 0.001)
node_load1 | node:load1:ratio{$1}
node_load5 | node:load5:ratio{$1}
node_load15 | node:load15:ratio{$1}
node_memory_utilisation | node:node_memory_utilisation:{$1}
node_memory_usage_wo_cache | node:node_memory_bytes_total:sum{$1} - node:node_memory_bytes_available:sum{$1}
node_memory_available | node:node_memory_bytes_available:sum{$1}
node_memory_total | node:node_memory_bytes_total:sum{$1}
node_net_utilisation | node:node_net_utilisation:sum_irate{$1}
node_net_bytes_transmitted | node:node_net_bytes_transmitted:sum_irate{$1}
node_net_bytes_received | node:node_net_bytes_received:sum_irate{$1}
node_disk_read_iops | node:data_volume_iops_reads:sum{$1}
node_disk_write_iops | node:data_volume_iops_writes:sum{$1}
node_disk_read_throughput | node:data_volume_throughput_bytes_read:sum{$1}
node_disk_write_throughput | node:data_volume_throughput_bytes_written:sum{$1}
node_disk_size_capacity | sum(max(node_filesystem_size_bytes{device=~"/dev/.*", device!~"/dev/loop\\d+", job="node-exporter"} * on (namespace, pod) group_left(node) node_namespace_pod:kube_pod_info:{$1}) by (device, node)) by (node)
node_disk_size_available | node:disk_space_available:{$1}
node_disk_size_usage | sum(max((node_filesystem_size_bytes{device=~"/dev/.*", device!~"/dev/loop\\d+", job="node-exporter"} - node_filesystem_avail_bytes{device=~"/dev/.*", device!~"/dev/loop\\d+", job="node-exporter"}) * on (namespace, pod) group_left(node) node_namespace_pod:kube_pod_info:{$1}) by (device, node)) by (node)
node_disk_size_utilisation | node:disk_space_utilization:ratio{$1}
node_disk_inode_total | node:node_inodes_total:{$1}
node_disk_inode_usage | node:node_inodes_total:{$1} - node:node_inodes_free:{$1}
node_disk_inode_utilisation | node:disk_inode_utilization:ratio{$1}
node_pod_count | node:pod_count:sum{$1}
node_pod_quota | max(kube_node_status_capacity{resource="pods",$1}) by (node) unless on (node) (kube_node_status_condition{condition="Ready",status=~"unknown\|false"} > 0)
node_pod_utilisation | node:pod_utilization:ratio{$1}
node_pod_running_count | node:pod_running:count{$1}
node_pod_succeeded_count | node:pod_succeeded:count{$1}
node_pod_abnormal_count | node:pod_abnormal:count{$1}
node_pod_abnormal_ratio | node:pod_abnormal:ratio{$1}
namespace_cpu_usage | round(namespace:container_cpu_usage_seconds_total:sum_rate{namespace!="", $1}, 0.001)
namespace_memory_usage | namespace:container_memory_usage_bytes:sum{namespace!="", $1}
namespace_memory_usage_wo_cache | namespace:container_memory_usage_bytes_wo_cache:sum{namespace!="", $1}
namespace_net_bytes_transmitted | sum by (namespace) (irate(container_network_transmit_bytes_total{namespace!="", pod!="", interface!~"^(cali.+\|tunl.+\|dummy.+\|kube.+\|flannel.+\|cni.+\|docker.+\|veth.+\|lo.*)", job="kubelet"}[5m]) * on (namespace) group_left(workspace) kube_namespace_labels{$1}) or on(namespace) max by(namespace) (kube_namespace_labels{$1} * 0)
namespace_net_bytes_received | sum by (namespace) (irate(container_network_receive_bytes_total{namespace!="", pod!="", interface!~"^(cali.+\|tunl.+\|dummy.+\|kube.+\|flannel.+\|cni.+\|docker.+\|veth.+\|lo.*)", job="kubelet"}[5m]) * on (namespace) group_left(workspace) kube_namespace_labels{$1}) or on(namespace) max by(namespace) (kube_namespace_labels{$1} * 0)
namespace_pod_count | sum by (namespace) (kube_pod_status_phase{phase!~"Failed\|Succeeded", namespace!=""} * on (namespace) group_left(workspace) kube_namespace_labels{$1}) or on(namespace) max by(namespace) (kube_namespace_labels{$1} * 0)
namespace_pod_running_count | sum by (namespace) (kube_pod_status_phase{phase="Running", namespace!=""} * on (namespace) group_left(workspace) kube_namespace_labels{$1}) or on(namespace) max by(namespace) (kube_namespace_labels{$1} * 0)
namespace_pod_succeeded_count | sum by (namespace) (kube_pod_status_phase{phase="Succeeded", namespace!=""} * on (namespace) group_left(workspace) kube_namespace_labels{$1}) or on(namespace) max by(namespace) (kube_namespace_labels{$1} * 0)
namespace_pod_abnormal_count | namespace:pod_abnormal:count{namespace!="", $1}
namespace_pod_abnormal_ratio | namespace:pod_abnormal:ratio{namespace!="", $1}
namespace_cronjob_count | sum by (namespace) (kube_cronjob_labels{namespace!=""} * on (namespace) group_left(workspace) kube_namespace_labels{$1})
namespace_pvc_count | sum by (namespace) (kube_persistentvolumeclaim_info{namespace!=""} * on (namespace) group_left(workspace) kube_namespace_labels{$1})
namespace_daemonset_count | sum by (namespace) (kube_daemonset_labels{namespace!=""} * on (namespace) group_left(workspace) kube_namespace_labels{$1})
namespace_deployment_count | sum by (namespace) (kube_deployment_labels{namespace!=""} * on (namespace) group_left(workspace) kube_namespace_labels{$1})
namespace_endpoint_count | sum by (namespace) (kube_endpoint_labels{namespace!=""} * on (namespace) group_left(workspace) kube_namespace_labels{$1})
namespace_hpa_count | sum by (namespace) (kube_hpa_labels{namespace!=""} * on (namespace) group_left(workspace) kube_namespace_labels{$1})
namespace_job_count | sum by (namespace) (kube_job_labels{namespace!=""} * on (namespace) group_left(workspace) kube_namespace_labels{$1})
namespace_statefulset_count | sum by (namespace) (kube_statefulset_labels{namespace!=""} * on (namespace) group_left(workspace) kube_namespace_labels{$1})
namespace_replicaset_count | count by (namespace) (kube_replicaset_labels{namespace!=""} * on (namespace) group_left(workspace) kube_namespace_labels{$1})
namespace_service_count | sum by (namespace) (kube_service_info{namespace!=""} * on (namespace) group_left(workspace) kube_namespace_labels{$1})
namespace_secret_count | sum by (namespace) (kube_secret_info{namespace!=""} * on (namespace) group_left(workspace) kube_namespace_labels{$1})
namespace_ingresses_extensions_count | sum by (namespace) (kube_ingress_labels{namespace!=""} * on (namespace) group_left(workspace) kube_namespace_labels{$1})
pod_cpu_usage | round(sum by (namespace, pod) (irate(container_cpu_usage_seconds_total{job="kubelet", pod!="", image!=""}[5m])) * on (namespace, pod) group_left(owner_kind, owner_name) kube_pod_owner{$1} * on (namespace, pod) group_left(node) kube_pod_info{$2}, 0.001)
pod_memory_usage | sum by (namespace, pod) (container_memory_usage_bytes{job="kubelet", pod!="", image!=""}) * on (namespace, pod) group_left(owner_kind, owner_name) kube_pod_owner{$1} * on (namespace, pod) group_left(node) kube_pod_info{$2}
pod_memory_usage_wo_cache | sum by (namespace, pod) (container_memory_working_set_bytes{job="kubelet", pod!="", image!=""}) * on (namespace, pod) group_left(owner_kind, owner_name) kube_pod_owner{$1} * on (namespace, pod) group_left(node) kube_pod_info{$2}
pod_net_bytes_transmitted | sum by (namespace, pod) (irate(container_network_transmit_bytes_total{pod!="", interface!~"^(cali.+\|tunl.+\|dummy.+\|kube.+\|flannel.+\|cni.+\|docker.+\|veth.+\|lo.*)", job="kubelet"}[5m])) * on (namespace, pod) group_left(owner_kind, owner_name) kube_pod_owner{$1} * on (namespace, pod) group_left(node) kube_pod_info{$2}
pod_net_bytes_received | sum by (namespace, pod) (irate(container_network_receive_bytes_total{pod!="", interface!~"^(cali.+\|tunl.+\|dummy.+\|kube.+\|flannel.+\|cni.+\|docker.+\|veth.+\|lo.*)", job="kubelet"}[5m])) * on (namespace, pod) group_left(owner_kind, owner_name) kube_pod_owner{$1} * on (namespace, pod) group_left(node) kube_pod_info{$2}
container_cpu_usage | round(sum by (namespace, pod, container) (irate(container_cpu_usage_seconds_total{job="kubelet", container!="POD", container!="", image!="", $1}[5m])), 0.001)
container_memory_usage | sum by (namespace, pod, container) (container_memory_usage_bytes{job="kubelet", container!="POD", container!="", image!="", $1})
container_memory_usage_wo_cache | sum by (namespace, pod, container) (container_memory_working_set_bytes{job="kubelet", container!="POD", container!="", image!="", $1})
pvc_inodes_available | max by (namespace, persistentvolumeclaim) (kubelet_volume_stats_inodes_free) * on (namespace, persistentvolumeclaim) group_left (storageclass) kube_persistentvolumeclaim_info{$1}
pvc_inodes_used | max by (namespace, persistentvolumeclaim) (kubelet_volume_stats_inodes_used) * on (namespace, persistentvolumeclaim) group_left (storageclass) kube_persistentvolumeclaim_info{$1}
pvc_inodes_total | max by (namespace, persistentvolumeclaim) (kubelet_volume_stats_inodes) * on (namespace, persistentvolumeclaim) group_left (storageclass) kube_persistentvolumeclaim_info{$1}
pvc_inodes_utilisation | max by (namespace, persistentvolumeclaim) (kubelet_volume_stats_inodes_used / kubelet_volume_stats_inodes) * on (namespace, persistentvolumeclaim) group_left (storageclass) kube_persistentvolumeclaim_info{$1}
pvc_bytes_available | max by (namespace, persistentvolumeclaim) (kubelet_volume_stats_available_bytes) * on (namespace, persistentvolumeclaim) group_left (storageclass) kube_persistentvolumeclaim_info{$1}
pvc_bytes_used | max by (namespace, persistentvolumeclaim) (kubelet_volume_stats_used_bytes) * on (namespace, persistentvolumeclaim) group_left (storageclass) kube_persistentvolumeclaim_info{$1}
pvc_bytes_total | max by (namespace, persistentvolumeclaim) (kubelet_volume_stats_capacity_bytes) * on (namespace, persistentvolumeclaim) group_left (storageclass) kube_persistentvolumeclaim_info{$1}
pvc_bytes_utilisation | max by (namespace, persistentvolumeclaim) (kubelet_volume_stats_used_bytes / kubelet_volume_stats_capacity_bytes) * on (namespace, persistentvolumeclaim) group_left (storageclass) kube_persistentvolumeclaim_info{$1}
etcd_server_list | label_replace(up{job="etcd"}, "node_ip", "$1", "instance", "(.*):.*")
etcd_server_total | count(up{job="etcd"})
etcd_server_up_total | etcd:up:sum
etcd_server_has_leader | label_replace(etcd_server_has_leader, "node_ip", "$1", "instance", "(.*):.*")
etcd_server_leader_changes | label_replace(etcd:etcd_server_leader_changes_seen:sum_changes, "node_ip", "$1", "node", "(.*)")
etcd_server_proposals_failed_rate | avg(etcd:etcd_server_proposals_failed:sum_irate)
etcd_server_proposals_applied_rate | avg(etcd:etcd_server_proposals_applied:sum_irate)
etcd_server_proposals_committed_rate | avg(etcd:etcd_server_proposals_committed:sum_irate)
etcd_server_proposals_pending_count | avg(etcd:etcd_server_proposals_pending:sum)
etcd_mvcc_db_size | avg(etcd:etcd_debugging_mvcc_db_total_size:sum)
etcd_network_client_grpc_received_bytes | sum(etcd:etcd_network_client_grpc_received_bytes:sum_irate)
etcd_network_client_grpc_sent_bytes | sum(etcd:etcd_network_client_grpc_sent_bytes:sum_irate)
etcd_grpc_call_rate | sum(etcd:grpc_server_started:sum_irate)
etcd_grpc_call_failed_rate | sum(etcd:grpc_server_handled:sum_irate)
etcd_grpc_server_msg_received_rate | sum(etcd:grpc_server_msg_received:sum_irate)
etcd_grpc_server_msg_sent_rate | sum(etcd:grpc_server_msg_sent:sum_irate)
etcd_disk_wal_fsync_duration | avg(etcd:etcd_disk_wal_fsync_duration:avg)
etcd_disk_wal_fsync_duration_quantile | avg(etcd:etcd_disk_wal_fsync_duration:histogram_quantile) by (quantile)
etcd_disk_backend_commit_duration | avg(etcd:etcd_disk_backend_commit_duration:avg)
etcd_disk_backend_commit_duration_quantile | avg(etcd:etcd_disk_backend_commit_duration:histogram_quantile) by (quantile)
apiserver_up_sum | apiserver:up:sum
apiserver_request_rate | apiserver:apiserver_request_total:sum_irate
apiserver_request_by_verb_rate | apiserver:apiserver_request_total:sum_verb_irate
apiserver_request_latencies | apiserver:apiserver_request_duration:avg
apiserver_request_by_verb_latencies | apiserver:apiserver_request_duration:avg_by_verb
scheduler_up_sum | scheduler:up:sum
scheduler_schedule_attempts | scheduler:scheduler_schedule_attempts:sum
scheduler_schedule_attempt_rate | scheduler:scheduler_schedule_attempts:sum_rate
scheduler_e2e_scheduling_latency | scheduler:scheduler_e2e_scheduling_duration:avg
scheduler_e2e_scheduling_latency_quantile | scheduler:scheduler_e2e_scheduling_duration:histogram_quantile

> - 第一列是在监控API调用时可用的指标名称，Prometheus中并未存储对应指标名称的时序数据。
> - 监控API调用时，将转换API参数为Prometheus标签选择器，代入到PromQL模板的`$1`,`$2`，然后使用生成的PromQL请求Prometheus服务。
> - Prometheus指标中，对于使用`:`分隔符命名的指标，通常是[Recording Rule](https://prometheus.io/docs/prometheus/latest/configuration/recording_rules/)派生的指标。

