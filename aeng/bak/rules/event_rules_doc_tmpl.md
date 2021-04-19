事件规则用来对K8S中的Event资源进行评估和告警



### 规则结构

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



### 内置规则更新

> 这里仅针对用于告警目的的内置事件规则进行操作

KubeSphere内置了一些必要的事件告警规则，对平台各类事件进行告警。各内置事件告警规则的定义请参考**附录**中的**内置事件告警规则**。

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

{{ if gt (len .Rules) 0 -}}

<table>
	<tr>
		<td>规则名称</td>
		<td>级别</td>
		<td>说明</td>
	</tr>
{{- range .Rules}}
	<tr>
		<td>{{.Name}}</td>
		<td>{{.Severity}}</td>
		<td>{{.Comment}}</td>
	</tr>
{{- end}}
</table>

{{end}}