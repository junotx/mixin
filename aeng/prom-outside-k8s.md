
需要解决集群外部 Prometheus 访问集群内部服务指标的问题，可选方式：

1. 选择通过 Apiserver 代理集群内部的服务指标，参考[文档](https://v1-23.docs.kubernetes.io/zh/docs/tasks/access-application-cluster/access-cluster-services/#discovering-builtin-services)。

    ```shell
    # 这里是一个通过 Apiserver 代理访问 kube-state-metrics 指标的示例
    # Apiserver 端和 kube-state-metrics 端都要求认证时，此请求会失败
    curl  -H "Authorization: Bearer $(kubectl -n kubesphere-monitoring-system get secret $(kubectl -n kubesphere-monitoring-system get sa kube-state-metrics -o jsonpath='{.secrets[0].name}') -o go-template='{{.data.token}}' | base64 -d)" --insecure "https://192.168.1.4:6443/api/v1/namespaces/kubesphere-monitoring-system/pods/http:<kube-state-metrics pod name>:8443/proxy/metrics"
    ```

2. 选择其他专门的代理组件比如 [americanexpress/k8s-prometheus-proxy](https://github.com/americanexpress/k8s-prometheus-proxy) 。