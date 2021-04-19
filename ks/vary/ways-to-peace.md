
# 监控

- 查看Prometheus的配置
  - 通过ui: `http://prom_ip:prom_port/config`  
  - 请求api: 
    ```shell
    curl http://prom_ip:prom_port/api/v1/status/config | jq --raw-output '.data.yaml'
    ```  
  - 查看绑定的配置文件：
    ```shell
    kubectl -n kubesphere-monitoring-system get secrets prometheus-k8s -o go-template='{{range $k, $v := .data}}{{$v | base64decode}}{{end}}' | gzip -d -
    ```  

## 服务发现

- 限定prometheus选择的servicemonitors：
  ```shell
  kubectl -n kubesphere-monitoring-system patch prometheus k8s --type=merge -p '{"spec":{"serviceMonitorSelector":{"matchLabels":{"app.kubernetes.io/vendor":"kubesphere"}}}}'
  ```
  > `--type=merge`: https://github.com/kubernetes/kubernetes/issues/71024
- 查看宕机的目标服务或Prometheus无法连接的目标服务：  
  - 通过ui查看那些状态为`DOWN`的目标服务: `http://prom_ip:prom_port/targets`  
  - 请求API获取：  
    ```shell
    curl prom_ip:prom_port/api/v1/targets | jq '[.data.activeTargets[] | select(.health == "down") | del(.discoveredLabels)]'
    ```
- 启用Etcd监控：  
  ```shell
  # 根据etcd证书创建secret，后者后续将被绑定到Prometheus以请求etcd服务
  # 登录到一个集群节点执行命令
  kubectl -n kubesphere-monitoring-system create secret generic kube-etcd-client-certs  \
  --from-file=etcd-client-ca.crt=/etc/ssl/etcd/ssl/ca.pem  \
  --from-file=etcd-client.crt=/etc/ssl/etcd/ssl/node-$(hostname).pem  \
  --from-file=etcd-client.key=/etc/ssl/etcd/ssl/node-$(hostname)-key.pem
  
  # 启用Etcd的监控
  kubectl -n kubesphere-system patch cc ks-installer --type merge -p '{"spec":{"etcd":{"monitoring":true,"tlsEnable":true}}}'
  ```

# 告警

## 异常告警

> 这里的异常告警处理方法主要是针对组件和环境自身的BUG问题以及规则问题所导致的告警  

**TargetDown**  

`TargetDown`在Target服务(Prometheus从这些服务的实例中pull指标)副本不可用率高的情况下告警。  

KubeSphere3.0版本或更高版本(非kk安装)的HA集群中，sheduler/controller-manager/apiserver服务可能会出现这样的告警。检查这些服务的绑定地址(`--bind-address`启动参数指定)，其为`127.0.0.1`时，会导致Prometheus无法连接它们。它们以静态pod方式运行，可以调整`/etc/kubernetes/manifests`(仅针对kubeadm安装的环境)中的该配置为`0.0.0.0`。具体可参考https://kubesphere.com.cn/forum/d/4849-31

**CPUThrottlingHigh**  

`CPUThrottlingHigh`在pod的cpu处于节制状态时间占比高时将进行告警。  

为避免这类告警，可以提高该规则内所指定的阈值，或者针对发生该告警的组件，提高所涉及容器的资源限额尤其是CPU限额，才能消除这类告警。  

该规则的原理可以参考[这里的描述](https://mp.weixin.qq.com/s?__biz=MzU1MzY4NzQ1OA==&mid=2247488086&idx=2&sn=d4ee70dbeec45961569b9b3b29339df9&chksm=fbee529bcc99db8d28f5123518012b2ccf8c300699932846d8448eb8afff94311ebb050b7446&mpshare=1&scene=1&srcid=0121XpIHrFxHMq3ARxFh5mEt&sharer_sharetime)。  

**AggregatedAPIDown**  

`AggregatedAPIDown`对不可用的AggregatedApi进行告警。比如它的一个告警消息的描述可能是`An aggregated API v1beta1.metrics.k8s.io/default is down`，此时可以查看对应的资源状态`kubectl -n default get apiservice v1beta1.metrics.k8s.io`。  

该规则在k8s1.18之前版本因为[`aggregator_unavailable_apiservice`的问题](https://github.com/kubernetes/kube-aggregator/issues/23)可能会有异常的告警，该问题被[kubernetes/kubernetes#87778](https://github.com/kubernetes/kubernetes/pull/87778)所修复。  

**KubeHpaReplicasMismatch**  

`KubeHpaReplicasMismatch`在HPA副本数不匹配时告警。以下是其异常告警的一个示例。

收到其告警的消息详情：`HPA istio-system/jaeger-collector has not matched the desired number of replicas for longer than 15 minutes. `。  

获取该HPA的期望副本为0：`kubectl -n istio-system get hpa jaeger-collector -o go-template='{{.status|printf "%v\n"}}'`

`kubectl -n kube-system logs -f $(kubectl -n kube-system get po --selector component=kube-controller-manager -o name | head -n 1) | grep jaeger-collector`可以发现副本计算异常如下：  

```console
I0729 03:09:18.032052       1 event.go:281] Event(v1.ObjectReference{Kind:"HorizontalPodAutoscaler", Namespace:"istio-system", Name:"jaeger-collector", UID:"7e346710-0eac-4745-9326-2204af4e6d11", APIVersion:"autoscaling/v2beta2", ResourceVersion:"28029119", FieldPath:""}): type: 'Warning' reason: 'FailedGetResourceMetric' missing request for memory
E0729 03:09:18.034070       1 horizontal.go:214] failed to compute desired number of replicas based on listed metrics for Deployment/istio-system/jaeger-collector: invalid metrics (2 invalid out of 2), first error is: failed to get memory utilization: missing request for memory
I0729 03:09:18.034132       1 event.go:281] Event(v1.ObjectReference{Kind:"HorizontalPodAutoscaler", Namespace:"istio-system", Name:"jaeger-collector", UID:"7e346710-0eac-4745-9326-2204af4e6d11", APIVersion:"autoscaling/v2beta2", ResourceVersion:"28029119", FieldPath:""}): type: 'Warning' reason: 'FailedGetResourceMetric' missing request for cpu
I0729 03:09:18.034158       1 event.go:281] Event(v1.ObjectReference{Kind:"HorizontalPodAutoscaler", Namespace:"istio-system", Name:"jaeger-collector", UID:"7e346710-0eac-4745-9326-2204af4e6d11", APIVersion:"autoscaling/v2beta2", ResourceVersion:"28029119", FieldPath:""}): type: 'Warning' reason: 'FailedComputeMetricsReplicas' invalid metrics (2 invalid out of 2), first error is: failed to get memory utilization: missing request for memory
```

查看istio-system中负载jaeger-collector的配置，没有`resources.request`，这导致配置hpa控制器计算异常后将期望副本数默认设置为0。[kubernetes/kubernetes#102728](https://github.com/kubernetes/kubernetes/pull/102728)将修复此问题。

**KubeCronJobRunning**  

`KubeCronJobRunning`针对执行任务耗时太久的CronJob进行告警。如果任务执行正常，依然收到这样的告警，可以参考以下的说明。  

比如收到一个这样的告警消息`CronJob istio-system/jaeger-es-index-cleaner is taking more than 1h to complete`，而查看其相关job都成功完成`kubectl -n istio-system get job --selector app.kubernetes.io/name=jaeger-es-index-cleaner`。  

该告警和指标`kube_cronjob_next_schedule_time`相关，后者由kube-state-metrics组件根据cronjob的上次调度时间和`schedule`表达式来[计算](https://github.com/kubernetes/kube-state-metrics/blob/v2.1.0/internal/store/cronjob.go#L245)。由于负责cronjob调度的kube-controller-manager组件绑定了节点时区，而kube-state-metrics组件默认使用了UTC时区，所以当节点时区为非UTC时区时，此时若没为`schedule`表达式指定时区，那么该指标的样本值与实际的下次调度时间就出现了差异。  

针对这种异常可以为cron表达式配置时区，通常通过前缀来指定，CRON_TZ=<time zone> 0 1 * * * ，可以参考[这里](https://en.wikipedia.org/wiki/Cron#Time_zone_handling)。或者通过以下方式为kube-state-metrics组件绑定主机时区

```shell
kubectl -n kubesphere-monitoring-system patch deploy kube-state-metrics --type='json' -p='[{"op": "add","path":"/spec/template/spec/volumes","value":[{"name":"host-time","hostPath":{"path":"/etc/localtime"}}]},{"op":"add","path":"/spec/template/spec/containers/0/volumeMounts","value":[{"name":"host-time","mountPath":"/etc/localtime"}]}]'
```

KubeSphere3.1.1之后的版本将自动为kube-state-metrics组件绑定主机时区，参考[kubesphere/ks-installer#1618](https://github.com/kubesphere/ks-installer/pull/1618)。  


**etcdHighNumberOfFailedGRPCRequests**  

`etcdHighNumberOfFailedGRPCRequests`在etcd集群的grpc请求失败率高时进行告警。如果etcd集群状态正常情况下，仍然收到如下的告警消息(关键词: `Watch failed`)  
```console
etcd cluster "etcd": 100%of requests fro Watch failed ont etcd instance 192.168.0.2:2379
```  
可以将对应规则表达式中的`grpc_code!="OK"`更新为`grpc_code=~"Unknown|FailedPrecondition|ResourceExhausted|Internal|Unavailable|DataLoss|DeadlineExceeded"`来解决。  

> 参考：[etcd-io/etcd#13127](https://github.com/etcd-io/etcd/pull/13127)  


## 规则加载失败

更新告警规则时，thanos ruler负载的sidecar容器configmap-reloader将请求主容器重载入规则。下述端点无法到达的异常将导致载入失败：  

```console
2021/03/30 10:14:37 error: Post http://localhost:10902/-/reload: dial tc [::1]:10902: connect: cannot assign requested address
2021/03/31 01:33:15 configmap updated
2021/03/31 01:33:15 successfully triggered reload
...
2021/03/31 03:23:57 configmap updated
2021/03/31 03:23:57 successfully triggered reload
2021/03/31 03:31:47 configmap updated
2021/03/31 03:31:47 error: Post http://localhost:10902/-/reload: dial tc [::1]:10902: connect: cannot assign requested address
```

configmap-reloader提供了webhook-retries参数进行重试，但是volume-dir参数不支持glob配置，所以通过ThanosRuler的crd暂时无法覆盖configmap-reloader容器配置。  

> [configmap-reload usage](https://github.com/jimmidyson/configmap-reload#usage) 和[ThanosRuler CRD](https://github.com/prometheus-operator/prometheus-operator/blob/v0.47.0/Documentation/api.md#thanosrulerspec) 中的containers配置。

## 日志之ES只读

通常ES在磁盘空间将满时，将启动保护机制，关闭索引的写功能。此时ES的报错日志可能如下：  

```console
Caused by: org.elasticsearch.ElasticsearchStatusException: Elasticsearch exception [type=cluster_block_exception, reason=blocked by: [FORBIDDEN/12/index read-only / allow delete (api)];]
```

可以通过扩展ES集群存储或删除过期索引恢复其可写状态。

<br/>

另外一种情况是存储系统只读。此时ES报错日志如下：  

```console
java.nio.file.FileSystemException: /usr/share/elasticsearch/data/nodes/0/indices/NHCkqsluRoycoyWVGgVNGw: Read-only file system
```

查看集群状态不再green，并发现有大量的shard下线：  

```shell
curl es-host:9200/_cat/health
curl es-host:9200/_cat/shards | grep UNASSIGNED
```

需要先解决磁盘写入问题，然后让下线的shard重新上线。集群级操作可使用API: `/_cluster/reroute`  

# k8s

[bash auto-completion on linux](https://kubernetes.io/docs/tasks/tools/included/optional-kubectl-configs-bash-linux/)    
[bash auto-completion on macOS](https://kubernetes.io/docs/tasks/tools/included/optional-kubectl-configs-bash-mac/)  

kubectl pretty-print: https://github.com/kubernetes/kubernetes/tree/master/staging/src/k8s.io/cli-runtime/pkg/printers  

