
# 内置告警策略调整参考

## 默认调整

这里对KubeSphere的内置告警规则进行了适应性调整，指标类告警规则保留了平台资源和配额、节点资源类告警规则，kube-apiserver、kubelet、kube-scheduler、kube-controller-manager、prometheus等平台组件的告警规则，以及k8s应用类的告警规则。事件告警规则仅保留启用了集群关键事件的告警规则。

请参考以下步骤更新到集群。

1. 内置指标告警规则  
   如果k8s版本大于等于v1.16，使用如下命令更新：
    ```yaml
    kubectl apply -f https://raw.githubusercontent.com/junotx/mixin/main/ks/ee/kuais/rules/prometheus-rules-v1.16+.yaml
    ```
    否则，请使用下列命令：
    ```yaml
    kubectl apply -f https://raw.githubusercontent.com/junotx/mixin/main/ks/ee/kuais/rules/prometheus-rules.yaml
    ```

2. 内置事件告警规则

    ```yaml
    kubectl apply -f https://raw.githubusercontent.com/junotx/mixin/main/prom/rules/kuais/ks-events-cluster-rules-default.yaml
    ```

## 自定义调整

请参考[这里](metric_rules_doc.md#配置说明-1)的文档自定义调整内置指标告警策略，参考[这里](event_rules_doc.md#内置事件告警规则更新)的文档自定义调整内置事件告警策略。