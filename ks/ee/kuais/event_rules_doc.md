事件规则用来对K8S中的Event资源进行评估和告警



### 事件规则结构

事件规则在KubeSphere平台存储在`rules.events.kubesphere.io`所定义的资源中，其Spec结构如下：

```
spec
└──rules
   |
   |  ┌──annotations
   |  |——condition
   0——|——enable         (取值true或false来启用或禁用规则，默认false)
   |  |——labels
   |  |——name
   |  └──type            (取值alert或notification表示规则的类型)
   |
   |
   1,2...
```

> 对于type=notification的非告警规则在当前的KubeSphere版本中未有应用，暂可忽略。



### 内置事件告警规则更新

> 这里仅针对用于告警目的的内置事件规则进行操作

KubeSphere内置了一些必要的事件告警规则，对平台各类事件进行告警。各内置事件告警规则的定义请参考后文的**内置事件告警规则**。

基本规则位于`kubesphere-logging-system`项目下的`ks-events-cluster-rules-default`资源中，通过以下命令可修改其中的规则：

```shell
kubectl -n kubesphere-logging-system edit rules.events.kubesphere.io ks-events-cluster-rules-default
```

> 该命令会进入到资源的编辑界面，编辑用法与linux中编辑文件的`vim`命令类似。

请参考前文的事件规则结构，对需要调整的告警规则进行操作，比如规则禁用、更新、删除等，然后保存后(同`vim`命令的保存操作)即可自动同步更新至EventsRuler组件(该组件负责加载事件规则、触发事件告警)。



当只针对个别的告警规则进行删除操作时，可以参考使用以下删除单个告警规则的快捷命令：

```shell
# 这里将删除ks-events-cluster-rules-default资源中名称为ContainerBackoff、级别为warning的告警规则
# 若要删除其他规则，请调整命令中相应位置处的规则名称和规则级别
kubectl -n kubesphere-logging-system get rules.events.kubesphere.io ks-events-cluster-rules-default -ojson | jq 'delpaths([path(..|select(.type?=="alert" and .name?=="ContainerBackoff" and .labels.severity?=="warning"))])' | kubectl apply -f -
```

### 内置事件规则表

<table>
	<tr>
		<td>规则名称</td>
		<td>级别</td>
		<td>说明</td>
	</tr>
	<tr>
		<td>ContainerFailed</td>
		<td>warning</td>
		<td>容器失败</td>
	</tr>
	<tr>
		<td>ContainerPreempting</td>
		<td>warning</td>
		<td>容器抢占中</td>
	</tr>
	<tr>
		<td>ContainerBackoff</td>
		<td>warning</td>
		<td>容器回退</td>
	</tr>
	<tr>
		<td>ContainerUnhealthy</td>
		<td>warning</td>
		<td>容器状态不良</td>
	</tr>
	<tr>
		<td>ContainerProbeWarning</td>
		<td>warning</td>
		<td>容器探测警告</td>
	</tr>
	<tr>
		<td>PodKillingExceededGracePeriod</td>
		<td>warning</td>
		<td>pod终止超时</td>
	</tr>
	<tr>
		<td>PodKillFailed</td>
		<td>warning</td>
		<td>pod终止失败</td>
	</tr>
	<tr>
		<td>PodContainerCreateFailed</td>
		<td>warning</td>
		<td>pod容器创建失败</td>
	</tr>
	<tr>
		<td>PodFailed</td>
		<td>warning</td>
		<td>pod失败</td>
	</tr>
	<tr>
		<td>PodNetworkNotReady</td>
		<td>warning</td>
		<td>Pod网络异常</td>
	</tr>
	<tr>
		<td>ImagePullPolicyError</td>
		<td>warning</td>
		<td>镜像拉取策略错误</td>
	</tr>
	<tr>
		<td>ImageInspectFailed</td>
		<td>warning</td>
		<td>镜像检查失败</td>
	</tr>
	<tr>
		<td>KubeletSetupFailed</td>
		<td>warning</td>
		<td>kubelet安装失败</td>
	</tr>
	<tr>
		<td>VolumeAttachFailed</td>
		<td>warning</td>
		<td>存储卷装载失败</td>
	</tr>
	<tr>
		<td>VolumeMountFailed</td>
		<td>warning</td>
		<td>存储卷挂载失败</td>
	</tr>
	<tr>
		<td>VolumeResizeFailed</td>
		<td>warning</td>
		<td>存储卷扩缩容失败</td>
	</tr>
	<tr>
		<td>FileSystemResizeFailed</td>
		<td>warning</td>
		<td>文件系统扩缩容失败</td>
	</tr>
	<tr>
		<td>VolumeMapFailed</td>
		<td>warning</td>
		<td>存储卷映射失败</td>
	</tr>
	<tr>
		<td>VolumeAlreadyMounted</td>
		<td>warning</td>
		<td>存储卷已被挂载</td>
	</tr>
	<tr>
		<td>NodeRebooted</td>
		<td>warning</td>
		<td>节点重启</td>
	</tr>
	<tr>
		<td>ContainerGCFailed</td>
		<td>warning</td>
		<td>容器GC失败</td>
	</tr>
	<tr>
		<td>ImageGCFailed</td>
		<td>warning</td>
		<td>镜像GC失败</td>
	</tr>
	<tr>
		<td>NodeAllocatableEnforcementFailed</td>
		<td>warning</td>
		<td>节点可分配资源更新失败</td>
	</tr>
	<tr>
		<td>SandboxCreateFailed</td>
		<td>warning</td>
		<td>Sandbox创建失败</td>
	</tr>
	<tr>
		<td>SandboxStatusFailed</td>
		<td>warning</td>
		<td>获取Sandbox状态错误</td>
	</tr>
	<tr>
		<td>DiskCapacityInvalid</td>
		<td>warning</td>
		<td>磁盘容量配置不合法</td>
	</tr>
	<tr>
		<td>DiskSpaceFreeFailed</td>
		<td>warning</td>
		<td>磁盘空间释放失败</td>
	</tr>
	<tr>
		<td>PodStatusSyncFailed</td>
		<td>warning</td>
		<td>Pod状态同步失败</td>
	</tr>
	<tr>
		<td>ConfigurationValidationFailed</td>
		<td>warning</td>
		<td>配置验证失败</td>
	</tr>
	<tr>
		<td>LifecycleHookPostStartFailed</td>
		<td>warning</td>
		<td>容器启动后的生命周期钩子运行失败</td>
	</tr>
	<tr>
		<td>LifecycleHookPreStopFailed</td>
		<td>warning</td>
		<td>容器停止前的生命周期钩子运行失败</td>
	</tr>
	<tr>
		<td>HPASelectorError</td>
		<td>warning</td>
		<td>HPA选择器错误</td>
	</tr>
	<tr>
		<td>HPAMetricError</td>
		<td>warning</td>
		<td>HPA对象指标错误</td>
	</tr>
	<tr>
		<td>HPAConvertFailed</td>
		<td>warning</td>
		<td>HPA转换失败</td>
	</tr>
	<tr>
		<td>HPAGetScaleFailed</td>
		<td>warning</td>
		<td>HPA规模获取失败</td>
	</tr>
	<tr>
		<td>HPAComputeReplicasFailed</td>
		<td>warning</td>
		<td>HPA副本计算失败</td>
	</tr>
	<tr>
		<td>HPARescaleFailed</td>
		<td>warning</td>
		<td>HPA规模调整失败</td>
	</tr>
	<tr>
		<td>NodeSystemOOM</td>
		<td>warning</td>
		<td>节点内存溢出</td>
	</tr>
	<tr>
		<td>VolumeBindingFailed</td>
		<td>warning</td>
		<td>存储卷绑定失败</td>
	</tr>
	<tr>
		<td>VolumeMismatch</td>
		<td>warning</td>
		<td>存储卷不匹配</td>
	</tr>
	<tr>
		<td>VolumeRecycleFailed</td>
		<td>warning</td>
		<td>存储卷回收失败</td>
	</tr>
	<tr>
		<td>VolumeRecyclerPodError</td>
		<td>warning</td>
		<td>存储卷回收器错误</td>
	</tr>
	<tr>
		<td>VolumeDeleteFailed</td>
		<td>warning</td>
		<td>存储卷删除失败</td>
	</tr>
	<tr>
		<td>VolumeProvisionFailed</td>
		<td>warning</td>
		<td>存储申请失败</td>
	</tr>
	<tr>
		<td>VolumeProvisionCleanupFailed</td>
		<td>warning</td>
		<td>清理存储失败</td>
	</tr>
	<tr>
		<td>VolumeExternalExpandingError</td>
		<td>warning</td>
		<td>存储外部扩展错误</td>
	</tr>
	<tr>
		<td>PodScheduleFailed</td>
		<td>warning</td>
		<td>pod调度失败</td>
	</tr>
	<tr>
		<td>PodCreateFailed</td>
		<td>warning</td>
		<td>pod创建失败</td>
	</tr>
	<tr>
		<td>PodDeleteFailed</td>
		<td>warning</td>
		<td>pod删除失败</td>
	</tr>
	<tr>
		<td>ReplicaSetCreateError</td>
		<td>warning</td>
		<td>副本集创建错误</td>
	</tr>
	<tr>
		<td>DeploymentRollbackFailed</td>
		<td>warning</td>
		<td>部署回滚失败</td>
	</tr>
	<tr>
		<td>DeploySelectorAll</td>
		<td>warning</td>
		<td>deploy选择了所有pod</td>
	</tr>
	<tr>
		<td>DaemonSelectorAll</td>
		<td>warning</td>
		<td>daemonset选择了所有pod</td>
	</tr>
	<tr>
		<td>DaemonPodFailed</td>
		<td>warning</td>
		<td>daemonset的pod失败</td>
	</tr>
	<tr>
		<td>LoadBalancerSyncFailed</td>
		<td>warning</td>
		<td>负载据衡器不可用</td>
	</tr>
	<tr>
		<td>LoadBalancerUnAvailable</td>
		<td>warning</td>
		<td>负载据衡器不可用</td>
	</tr>
	<tr>
		<td>LoadBalancerUpdateFailed</td>
		<td>warning</td>
		<td>更新负载据衡器失败</td>
	</tr>
	<tr>
		<td>LoadBalancerDeleteFailed</td>
		<td>warning</td>
		<td>负载据衡器删除失败</td>
	</tr>
	<tr>
		<td>JobGetFailed</td>
		<td>warning</td>
		<td>任务获取失败</td>
	</tr>
	<tr>
		<td>JobCreateFailed</td>
		<td>warning</td>
		<td>任务创建失败</td>
	</tr>
	<tr>
		<td>JobDeleteFailed</td>
		<td>warning</td>
		<td>任务删除失败</td>
	</tr>
	<tr>
		<td>JobUnexpected</td>
		<td>warning</td>
		<td>任务非预期</td>
	</tr>
	<tr>
		<td>JobScheduleFailed</td>
		<td>warning</td>
		<td>任务调度失败</td>
	</tr>
</table>

