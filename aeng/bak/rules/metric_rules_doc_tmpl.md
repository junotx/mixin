KubeSphere平台从3.1版本开始重新调整了告警系统的设计，可以兼容Prometheus风格的告警规则。这里所提及的告警仅针对各类资源的指标进行告警。

KubeSphere定义了自定义告警策略的API，并提供了界面化的交互式UI，方便用户对自定义告警策略的增删改查操作。区别于自定义告警策略，3.1版本之前Prometheus所内置的告警规则保留为内置告警策略，为兼容平台外部Prometheus，Console仅提供了内置告警策略的查看功能。

为支持多租户场景，自定义告警策略分为了集群和项目两个层级。这两个层级都针对常用的指标告警场景提供了便捷的模板化配置，也开放了自定义PromQL来满足复杂的业务。



## 告警策略存储方式

无论是自定义告警策略，还是内置告警策略(这里仅指平台内置Prometheus的策略)，它们都首先存储在prometheus-operator所定义的`prometheusrules.monitoring.coreos.com`资源中。该资源的Spec结构请参考下图。这些资源的更新将由prometheus-operator同步至告警系统中。

```
spec
└──groups
   |
   |  ┌──name     (group name)
   0──|
   |  └──rules
   |      |
   |      |  ┌──expr
   |      0——|──labels        (recording rule)
   |      |  └──record
   |      |
   |      1,2...
   |
   |
   |  ┌──name      (group name)
   1——|
   |  └──rules
   |      |
   |      |  ┌──alert
   |      |  |──annotations
   |      0——|——expr          (alerting rule)
   |      |  |——for
   |      |  └──labels
   |      |
   |      1,2...
   |
   2,3...
```

> 这里请只参考告警规则，即alerting rules。  
> 一个规则组中的规则通常只包括recording rules或只包括alerting rules。

平台默认配置下，可以通过命令`kubectl -n kubesphere-monitoring-system get prometheusrules -l prometheus=k8s,role=alert-rules`获取所有内置告警策略存储的资源，通过命令`kubectl get prometheusrules -l thanosruler=thanos-ruler,role=thanos-alerting-rules -A`获取自定义告警策略存储的资源。

> 请勿手动修改自定义告警策略的CRD资源，而应通过Console或API调用来更新策略。

以下是单个告警策略的存储结构说明：

```yaml
alert: <string>
expr: <string>
for: <duration>
labels:
  [<label_name>: <label_value>...]
annotations:
  [<annotation_name>: <annotation_value>...]
```

- `alert`: 策略名称/规则名称/告警名称。
- `expr`: 规则表达式，一个合法的PromQL表达式。
- `for`: 告警持续时间。达到该持续时间的告警消息才被下发。
- `labels`: 标签集。通常会有一个名称为`severity`，值为`warning`/`error`/`critical`的标签来标识告警的严重程度。这些labels将被加入到告警消息的labels中。
- `annotations`: 注解集。用来丰富通知消息的内容。通常会有一个名称为`summary`的注解说明告警消息的摘要信息，和一个名称为`message`的注解说明告警消息的详细信息。



## 告警原理说明

这里以`TargetDown`这个内置告警策略为例，进行告警原理的说明。

该策略的目的是，针对Prometheus的抓取目标服务异常情况进行告警，当某个目标服务的副本不可用率大于10%，且持续超过10分钟时，发送告警消息。

```yaml
alert: TargetDown
annotations:
  message: >-
    {{`{{ printf "%.4g" $value }}`}}% of the {{`{{ $labels.job }}`}}/{{`{{
    $labels.service }}`}} targets in {{`{{ $labels.namespace }}`}} namespace
    are down.
expr: >-
  100 * (count(up == 0) BY (job, namespace, service) / count(up) BY
  (job, namespace, service)) > 10
for: 10m
labels:
  severity: warning
```

告警系统在发现该策略后，将通过`expr`表达式来周期性地查询指标系统，结果集将是副本不可用率大于10%的那些目标服务。如果在`for`所指定的时间范围内，每次查询的结果集之中都包含目标服务A，那么，以`TargetDown`命名且包含服务A属性的告警消息，就将被发送到下游通知系统。这之后，如果查询结果集中继续包含A服务，相应的告警消息将继续发送，反之则在下次查询结果集中出现服务A时进行重新计时，直到再次满足`for`所指定的时间范围。

告警消息主要包括`alertname`，`labels`和`annotations`三个属性。`alertname`来自于告警策略名称，`labels`包含了表达式查询结果中的`labels`和告警策略中的`labels`，`annotations`来自于告警策略的`annotations`。

> 告警策略中的`annotaions`支持配置模板，具体请参考[这里](https://prometheus.io/docs/prometheus/latest/configuration/alerting_rules/#templating)，模板执行后的结果会放在告警消息中的`annotations`。
>
> 持续时间未设置或被设置为0时，告警系统通过expr查询到结果后，将不等待就发送告警消息到下游。



## 自定义告警策略

### 配置说明

具体的配置步骤请参考[集群告警策略配置文档](https://v3-1.docs.kubesphere.io/zh/docs/cluster-administration/cluster-wide-alerting-and-notification/alerting-policy/)和[项目告警策略配置文档](https://v3-1.docs.kubesphere.io/zh/docs/project-user-guide/alerting/alerting-policy/)。

这里针对自定义告警策略API所定义的告警策略数据传输结构与Console上的界面要素的关系进行说明。前者的结构如下：

```yaml
name: <rule_name>
query: <query_string>
duration: <duration>
labels:
  [<label_name>: <label_value>...]
annotations:
  [<annotation_name>: <annotation_value>...]
```

> 这里的`name`, `query`, `duration`，分别与前文告警策略存储结构中的`alert`, `expr`, `for`一一对应。

策略*名称*：自定义策略的名称要求是一个合法的k8s资源名称，这与Prometheus有所区别。

*持续时间*：对应到`duration`属性。

*告警级别*：将作为标签添加到`labels`中，标签名`severity`，标签值支持`warning`、`error`、`critical`，依次是*一般告警*、*重要告警*、*危险告警*，告警的严重程度或紧急程度依次递增。

*规则模板*和*自定义规则*：使用规则模板配置时，将根据输入来自动组装PromQL表达式，填充到`query`中。而自定义规则则是直接配置`query`表达式。

*通知内容*: 用来丰富告警消息的内容。其中的*标题*作为名称为`summary`的注解添加到`annotations`，*消息*则对应到名称为`message`的注解。

对`query`表达式和`duration`的评估结果，决定了告警消息是否产生和是否下发。据此，告警策略的*告警状态*分为了以下三种：

- *未触发*：表示二者条件都不满足，此时未产生告警消息。
- *待触发*：表示满足`query`但不满足`duration`，可以认为此时已产生告警消息，但暂未下发。
- *触发中*，表示二者条件都满足，此时已开始(或正准备)发送告警消息到下游通知系统。


### 模板规则配置参考

#### 集群级别

Console上集群层级的告警策略提供了配置模板，可以针对节点的CPU、内存、本地磁盘、网络等各种资源类指标，进行快速的告警规则配置。下表列出了针对这些指标的建议阈值配置，提供给用户配置时参考。

| 指标名称           | 操作符      | 建议阈值 | 单位   |
| ------------------ | ----------- | -------- | ------ |
| 容器组异常率       | `>` 或 `>=` | 3        | %      |
| 容器组利用率       | `>` 或 `>=` | 80       | %      |
| CPU利用率          | `>` 或 `>=` | 80       | %      |
| CPU 1分钟平均负载  | `>` 或 `>=` | -        | Core   |
| CPU 5分钟平均负载  | `>` 或 `>=` | -        | Core   |
| 可用内存           | `<` 或 `<=` | -        | GB     |
| 内存利用率         | `>` 或 `>=` | 80       | %      |
| 本地磁盘可用空间   | `<` 或 `<=` | -        | GB     |
| 本地磁盘空间利用率 | `>` 或 `>=` | 80       | %      |
| inode利用率        | `>` 或 `>=` | 80       | %      |
| 本地磁盘读取IOPS   | `>` 或 `>=` | -        | 次数/s |
| 本地磁盘写入IOPS   | `>` 或 `>=` | -        | 次数/s |
| 本地磁盘读取吞吐量 | `>` 或 `>=` | -        | KB/s   |
| 本地磁盘写入吞吐量 | `>` 或 `>=` | -        | KB/s   |
| 网络发送数据速率   | `>` 或 `>=` | -        | Mbps   |
| 网络接收数据速率   | `>` 或 `>=` | -        | Mbps   |

> - 单位已由Console指定，配置时无需设定。
> - 未给出建议阈值的指标，用户请根据平台规模和业务需要自行配置。


#### 项目级别

Console为项目层级的告警规则配置，提供了针对部署、有状态副本集、守护进程集等工作负载，CPU用量、内存用量、网路数据收发速率、副本不可用率等指标在内的模板化告警规则配置

| 指标名称           | 操作符      | 建议阈值 | 单位 |
| ------------------ | ----------- | -------- | ---- |
| CPU用量            | `>` 或 `>=` | -        | Core |
| 内存用量           | `>` 或 `>=` | -        | Mi   |
| 内存用量(包含缓存) | `>` 或 `>=` | -        | Mi   |
| 网络发送数据速率   | `>` 或 `>=` | -        | Kbps |
| 网络接收数据速率   | `>` 或 `>=` | -        | Kbps |
| 副本不可用率       | `>` 或 `>=` | -        | %    |

> - 单位已由Console指定，配置时无需设定。
> - 这里未给出建议阈值，请根据实际业务需求进行配置。



## 内置告警策略

KubeSphere内置了一些必要的告警策略，对平台物理资源、应用资源、关键性组件的各类指标进行告警。这些内置告警策略将由Prometheus组件来评估和告警，它们的含义请参考**附录一：内置告警规则表**。

通过集群管理的告警策略页可以查询和查看内置告警策略。通常不建议对这些内置告警策略进行调整，若有需求，请参考后续的配置说明。

> 通过命令`kubectl -n kubesphere-monitoring-system get prometheusrules -l prometheus=k8s,role=alert-rules`可以获取存储内置告警策略的资源。



### 配置说明

内置告警策略的绝大部分位于`kubesphere-monitoring-system`项目下的`prometheus-k8s-rules`资源中，该资源的结构请参考前述的**告警策略存储方式**。通过以下命令可修改其中的策略规则：

```shell
kubectl -n kubesphere-monitoring-system edit prometheusrules.monitoring.coreos.com prometheus-k8s-rules
```

> 该命令会进入到资源的编辑界面，编辑用法与linux中编辑文件的`vim`命令类似。

请参考前文的告警策略结构，对需要调整的告警策略进行操作，比如更新、删除等，然后保存后(同`vim`命令的保存操作)即可自动同步更新至Prometheus组件。

当只针对个别的内置告警策略进行删除操作时，请参考使用以下删除单个告警策略的快捷命令：

```shell
# 这里将删除prometheus-k8s-rules资源中名称为KubePodCrashLooping、级别为warning的告警规则
# 若要删除其他规则，请调整命令中相应位置处的规则名称和规则级别
kubectl -n kubesphere-monitoring-system get prometheusrules.monitoring.coreos.com prometheus-k8s-rules -ojson | jq 'delpaths([path(..|select(.alert?=="KubePodCrashLooping" and .labels.severity?=="warning"))])' | kubectl apply -f -
```

> etcd相关的内置告警策略位于`kubesphere-monitoring-system`项目下的`prometheus-k8s-etcd-rules`资源中。

## 附录一：内置告警策略表

<table>
	<tr>
		<td>组</td>
		<td>规则名称</td>
		<td>持续时间</td>
		<td>级别</td>
		<td>说明</td>
	</tr>
{{- $allRules := .Rules -}}
{{- range .Groups -}}
	{{- $groupName := . -}}
	{{- $rules := (index $allRules $groupName) -}}
	{{- if gt (len $rules) 0 -}}
		{{- range $idx, $rule := $rules}}
	<tr>
			{{- if eq $idx 0}}
		<td rowspan={{len $rules}}>{{$groupName}}</td>
			{{- end}}
		<td>{{$rule.Name}}</td>
		<td>{{$rule.For}}</td>
		<td>{{$rule.Severity}}</td>
		<td>{{$rule.Comment}}{{if ne $rule.Supple ""}}({{$rule.Supple}}){{end}}</td>
	</tr>
		{{- end -}}
	{{- end -}}
{{- end}}
</table>