定制的规则和监控：

1. 平台基础指标告警规则  
   限制了规则生效的namespace范围，以下规则将不处理：

   - 无需ns限制的特定组件告警
   - 项目资源用量的超限告警
   - PVC存储临界告警

   配置根据k8s版本不同：

   - k8s version < v1.16.0  
     ```shell
     export metricNsFilter='namespace=~"kube-.*|kubesphere-.*|trident|yunshan-deepflow|ump|orion|f5-bigip-ctlr|operators"' && curl -s https://raw.githubusercontent.com/junotx/mixin/main/ks/ee/gd/prometheus-rules.yaml | envsubst '${metricNsFilter}' | kubectl apply --force -f -
     ```
   - k8s version >= v1.16.0: https://raw.githubusercontent.com/junotx/mixin/main/ks/ee/gd/prometheus-rules-v1.16+.yaml
     ```shell
     export metricNsFilter='namespace=~"kube-.*|kubesphere-.*|trident|yunshan-deepflow|ump|orion|f5-bigip-ctlr|operators"' && curl -s https://raw.githubusercontent.com/junotx/mixin/main/ks/ee/gd/prometheus-rules-v1.16+.yaml | envsubst '${metricNsFilter}' | kubectl apply --force -f -
     ```
   > 获取k8s版本：`kubectl version -o json | jq '.serverVersion.gitVersion' | sed s/\"//g`
   
2. etcd告警规则  
   使用下列命令配置：  
   
   ```shell
   kubectl apply --force -f https://raw.githubusercontent.com/junotx/mixin/main/ks/ee/gd/prometheus-rulesEtcd.yaml
   ```
   
3. ks额外告警规则  
   定制如下：
   
   - 容器cpu和内存用量临界的告警（限制了生效的namespace范围）
   - ks-apiserver请求响应效率的告警
   - thanos ruler组件的相关告警
   
   使用下列命令配置：
   
   ```shell
   export metricNsFilter='namespace=~"kube-.*|kubesphere-.*|trident|yunshan-deepflow|ump|orion|f5-bigip-ctlr|operators"' && curl -s https://raw.githubusercontent.com/junotx/mixin/main/ks/ee/gd/kubesphere-rules.yaml | envsubst '${metricNsFilter}' | kubectl apply --force -f -
   ```

4. coredns告警规则  
   配置根据coredns版本不同：

   - coredns version < 1.7.0: 
   ```shell
   kubectl apply --force -f https://raw.githubusercontent.com/junotx/mixin/main/ks/ee/gd/coredns-rules.yaml
   ```
   - coredns version >= 1.7.0: 
   ```shell
   kubectl apply --force -f https://raw.githubusercontent.com/junotx/mixin/main/ks/ee/gd/coredns-rules-1.7.0+.yaml
   ```
   
   > 获取coredns版本：`kubectl -n kube-system get deployment coredns -o jsonpath='{.spec.template.spec.containers[0].image}' | awk -F ':' '{print $2}'`

5. 事件告警规则  
   定制如下：
   - 限制规则生效的namespace范围（无需ns限制的节点事件告警除外）
   - 禁用部分与指标告警可能重叠的规则以及一些操作即现型事件的规则

   使用下列命令配置：
   ```shell
   export metricNsFilter='namespace=~"kube-.*|kubesphere-.*|trident|yunshan-deepflow|ump|orion|f5-bigip-ctlr|operators"' && curl -s https://raw.githubusercontent.com/junotx/mixin/main/ks/ee/gd/ks-events-cluster-rules-default.yaml | envsubst '${metricNsFilter}' | kubectl apply --force -f -
   ```

6. ks-apiserver和thanos ruler监控启用：

   ```shell
   kubectl apply --force -f https://raw.githubusercontent.com/junotx/mixin/main/ks/ee/gd/service-monitors-ks.yaml
   ```
