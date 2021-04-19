

# otel-collector

> [OpenTelemetry](https://github.com/open-telemetry) 基于跟踪、指标、日志这三大可观测性数据，定义了统一的 API 标准，并提供了 SDK 和工具库用于数据的创建、导出、采集和传输。  

[OpenTelemetry Collector](https://github.com/open-telemetry/opentelemetry-collector.git) 用于各类可观测数据的统一采集和传输。其以多管道的形式来实现，用管道的类型来定义所处理的数据类型。并提供了插件式的扩展方案，以下列出部分插件：

receiver | processor | exporter
--- | --- | ---  
hostmetrics | batch | file
jaeger | filter | prometheus
prometheus | memorylimiter | prometheusremotewrite
otlp | *groupbyattrs* | otlp
*filelog* | *k8s* | *elasticsearch*
*k8scluster* | *routing* | *loki*
| | | *awsprometheusremotewrite*

注意事项：

- [配置中`$`需转义](https://github.com/open-telemetry/opentelemetry-collector/blob/main/receiver/prometheusreceiver/README.md#getting-started)


其他参考：

- [glossary](https://opentelemetry.io/docs/concepts/glossary/)
- [opentelemetry-collector-builder](https://github.com/open-telemetry/opentelemetry-collector-builder.git): 用来构建自己的otel-collector分发版，例如[aws-otel-collector](https://github.com/aws-observability/aws-otel-collector.git)  
- [OpenTelemetry Collector WAL Design(draft)](https://docs.google.com/document/d/1Y4vNthCGdYI61ezeAzL5dXWgiZ73y9eSjIDitk3zXsU/edit#)

