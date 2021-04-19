[jmx_exporter](https://github.com/prometheus/jmx_exporter.git) 用来导出 JMX Beans (指标对象)：

- 通常作为一个 Java Agent 运行并导出当前 JVM 中的应用指标。
- 也可以作为一个独立的 HTTP 服务，通过 RMI 获取远端的 JMX 目标对象即指标，这种方式配置相当复杂，不能暴露内存和 CPU 用量等进程指标。  

> JMX 默认在本地主机的 IP 地址上提供 RMI 服务，此时 jmx_exporter 的 jmxUrl 需要以 ip 方式配置，例如 `service:jmx:rmi:///jndi/rmi://127.0.0.1:1234/jmxrmi`。(是否可通过 `-Djava.rmi.server.hostname` 设置 hostname 尚需验证)  

jmx_exporter 需要为每一个 java 进程启动一个 exporter ，[gentlezuo/jmx_exporter_multi](https://github.com/gentlezuo/jmx_exporter_multi.git) 实现了抓取多个应用的 JMX 对象。  

[banzaicloud/prometheus-jmx-exporter-operator](https://github.com/banzaicloud/prometheus-jmx-exporter-operator.git) 实现了在已运行的 pod 中载入 jmx_exporter 来暴露 pod 中 java 进程的指标。  

> 实际载入 jmx_exporter 由 [banzaicloud/jmx-exporter-loader](https://github.com/banzaicloud/jmx-exporter-loader.git) 完成。    
> jmx_exporter 默认在业务应用运行前加载执行，[banzaicloud/jmx_exporter](https://github.com/banzaicloud/jmx_exporter.git) 通过[一些处理](https://github.com/banzaicloud/jmx_exporter/commit/e83a7f123a983402aac2d831a716da4f4cd1ed5d)实现了[运行时载入](https://github.com/banzaicloud/jmx_exporter/commit/e83a7f123a983402aac2d831a716da4f4cd1ed5d)。  