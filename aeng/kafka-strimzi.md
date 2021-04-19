
[strimzi-kafka-operator](https://github.com/strimzi/strimzi-kafka-operator.git)提供了一种在Kubernetes平台上创建、配置和管理Kafka集群的方式。

# 配置参考

[KafkaSpec](https://strimzi.io/docs/operators/0.25.0/using.html#type-KafkaSpec-reference):
- kafka:
    - config: 未被operator纳管的一些kafka配置项
    - storage: 支持[Ephemeral](https://strimzi.io/docs/operators/0.25.0/using.html#ref-ephemeral-storage-str), [Persistent](https://strimzi.io/docs/operators/0.25.0/using.html#ref-persistent-storage-str), [JBOD storage](https://strimzi.io/docs/operators/0.25.0/using.html#ref-jbod-storage-str)三种存储类型
    - livenessProbe
    - readinessProbe: 就绪探针。broker就绪与否是通过[kafka-agent](https://github.com/strimzi/strimzi-kafka-operator/tree/main/kafka-agent)监测[broker状态](https://docs.itrsgroup.com/docs/geneos/5.3.0/Integrations/Kafka/kafka_monitoring_tr.html)来实现的(`brokerState=3`时认为已就绪)
    - jvmOptions
    - jmxOptions
    - metricsConfig: 配置[jmx_exporter](https://github.com/prometheus/jmx_exporter.git)需要导出的指标。配置[样例](https://github.com/strimzi/strimzi-kafka-operator/blob/main/examples/metrics/kafka-metrics.yaml)
- zookeeper：分布式协调服务。维护/协调kafka集群的broker, topic, partition等的一致性
- entityOperator: entity-operator组件。包含了topic和user的operator，使可以通过创建crd资源的方式创建topic和user
- cruiseControl：kafka集群的巡检控制组件。kafka集群自动化调整系统，实现动态的负载均衡，自检和自修复。strimzi中Cruise Control支持的特性参考[说明](https://strimzi.io/docs/operators/0.25.0/using.html#cruise-control-concepts-str)
- jmxTrans: [jmxtrans](https://github.com/jmxtrans/jmxtrans.git)组件。收集jmx指标
- kafkaExporter: [kafka-exporter](https://github.com/danielqsj/kafka_exporter.git)组件。提供kafka额外的监控指标，比如[topic/partition的消费延迟](https://strimzi.io/docs/operators/0.25.0/deploying.html#con-metrics-kafka-exporter-lag-str)。配置[样例](https://strimzi.io/docs/operators/0.25.0/deploying.html#proc-kafka-exporter-configuring-str)。

> JBOD(just a bunch of disk), [RAID](https://zhuanlan.zhihu.com/p/51170719)(Redundant Array of Independent Disks)

# 文档参考

- [部署前准备](https://strimzi.io/docs/operators/0.25.0/deploying.html#deploy-tasks-prereqs_str)
- [部署](https://strimzi.io/docs/operators/0.25.0/deploying.html#deploy-tasks_str)
    - [部署operator](https://strimzi.io/docs/operators/0.25.0/deploying.html#cluster-operator-str)
    - [部署kafka](https://strimzi.io/docs/operators/0.25.0/deploying.html#kafka-cluster-str)
- [监控](https://strimzi.io/docs/operators/0.25.0/deploying.html#assembly-metrics-str)
- [升级](https://strimzi.io/docs/operators/0.25.0/deploying.html#assembly-upgrade-str)
- [版本支持](https://strimzi.io/downloads/)

# 快速安装

文档参考自https://strimzi.io/docs/operators/0.25.0/quickstart.html#assembly-evaluation-str  

请首先创建好`kafka`和`my-kafka-project`这两个命名空间，前者用来安装operator，后者用来安装kafka集群。

```shell
kubectl create ns kafka

kubectl create ns my-kafka-project
```

## 安装operator

请确保这里的操作用户拥有k8s集群管理的权限。具体的安装步骤如下：

1. 获取安装文件

   ```shell
   STRIMZI_VERSION="0.25.0"
   
   curl -L https://github.com/strimzi/strimzi-kafka-operator/releases/download/${STRIMZI_VERSION}/strimzi-${STRIMZI_VERSION}.tar.gz | tar -C ./ -xz
   
   cd strimzi-${STRIMZI_VERSION}
   ```

2. 调整operator的默认配置

  - 修改RoleBinding和ClusterRoleBinding资源中所绑定的ServiceAccount，将其namespace替换为`kafka`

    ```shell
    sed -i 's/namespace: .*/namespace: kafka/' install/cluster-operator/*RoleBinding*.yaml
    ```

  - 配置operator需要watch的命名空间。operator将watch其中的相关资源来进行后续kafka集群的部署。

    通过以下命令：

    ```shell
    vim install/cluster-operator/060-Deployment-strimzi-cluster-operator.yaml
    ```

    调整`STRIMZI_NAMESPACE`为`my-kafka-project`

    ```yaml
    # ...
    env:
    - name: STRIMZI_NAMESPACE
      value: my-kafka-project
    # ...
    ```

    > watch[多个命名空间](https://github.com/strimzi/strimzi-kafka-operator/blob/0.25.0/documentation/modules/deploying/proc-deploy-cluster-operator-watch-multiple-namespaces.adoc)请以`,`分隔，watch[所有命名空间](https://github.com/strimzi/strimzi-kafka-operator/blob/0.25.0/documentation/modules/deploying/proc-deploy-cluster-operator-watch-whole-cluster.adoc)可直接赋值为`*`。

3. 安装部署operator

   ```shell
   kubectl create -f install/cluster-operator/ -n kafka
   ```

4. 授权operator。operator需要拥有操作命名空间`my-kafka-project`中相关资源的权限。

   ```shell
   kubectl create -f install/cluster-operator/020-RoleBinding-strimzi-cluster-operator.yaml -n my-kafka-project
   
   kubectl create -f install/cluster-operator/031-RoleBinding-strimzi-cluster-operator-entity-operator-delegation.yaml -n my-kafka-project
   ```



## 安装Kafka集群

operator安装好后，接下来在命名空间`my-kafka-project`中创建kafka集群。这里使用非特权用户即可。

operator支持创建kafka集群的同时创建zookeeper集群，具体请参考以下操作：

1. 这里创建一个名为`my-cluster`的集群，集群包括一个kafka broker节点和一个zookeeper节点

   ```shell
   cat << EOF | kubectl create -n my-kafka-project -f -
   apiVersion: kafka.strimzi.io/v1beta2
   kind: Kafka
   metadata:
     name: my-cluster
   spec:
     kafka:
       replicas: 1
       listeners:
         - name: plain
           port: 9092
           type: internal
           tls: false
         - name: tls
           port: 9093
           type: internal
           tls: true
           authentication:
             type: tls
         - name: external
           port: 9094
           type: nodeport
           tls: false
       storage:
         type: jbod
         volumes:
         - id: 0
           type: persistent-claim
           # class: my-storage-class
           size: 10Gi
           deleteClaim: false
       config:
         offsets.topic.replication.factor: 1
         transaction.state.log.replication.factor: 1
         transaction.state.log.min.isr: 1
     zookeeper:
       replicas: 1
       storage:
         type: persistent-claim
         # class: my-storage-class
         size: 5Gi
         deleteClaim: false
     entityOperator:
       topicOperator: {}
       userOperator: {}
   EOF
   ```

   > 这里使用k8s上所配置的默认持久化存储类型

2. 等待集群创建成功

   ```shell
   kubectl wait kafka/my-cluster --for=condition=Ready --timeout=300s -n my-kafka-project
   ```

3. 为kafka集群创建一个topic

   ```shell
   cat << EOF | kubectl create -n my-kafka-project -f -
   apiVersion: kafka.strimzi.io/v1beta2
   kind: KafkaTopic
   metadata:
     name: my-topic
     labels:
       strimzi.io/cluster: "my-cluster"
   spec:
     partitions: 3
     replicas: 1
   EOF
   ```



## 收发消息测试

针对创建好的`my-cluster`集群进行测试和验证。这里需要在本地机器上运行kafka的producer和consumer。

1. 在本地机器上获取kafka执行包

   ```shell
   curl -L https://dlcdn.apache.org/kafka/3.0.0/kafka_2.13-3.0.0.tgz | tar -C ./ -xz
   ```

2. 获取kafka服务的地址和端口

   获取kafka服务的端口：

   ```shell
   kubectl get service my-cluster-kafka-external-bootstrap -n my-kafka-project -o=jsonpath='{.spec.ports[0].nodePort}{"\n"}'
   ```

   服务的地址使用节点的地址：

   ```shell
   kubectl get nodes --output=jsonpath='{range .items[*]}{.status.addresses[?(@.type=="InternalIP")].address}{"\n"}{end}'
   ```

3. 启动一个producer。然后发送一些消息。

   ```shell
   bin/kafka-console-producer.sh --broker-list <node-address>:<node-port> --topic my-topic
   ```

4. 启动一个consumer。验证接收到的消息。

   ```shell
   bin/kafka-console-consumer.sh --bootstrap-server <node-address>:<node-port> --topic my-topic --from-beginning
   ```


# 监控


> 参考 https://github.com/strimzi/strimzi-kafka-operator/blob/0.25.0/examples/metrics/kafka-metrics.yaml

1. 创建jvm_exporter的导出指标配置  

  ```shell
  cat <<\EOF | kubectl -n my-kafka-project apply -f -
  kind: ConfigMap
  apiVersion: v1
  metadata:
    name: kafka-metrics
    labels:
      app: strimzi
  data:
    kafka-metrics-config.yml: |
      # See https://github.com/prometheus/jmx_exporter for more info about JMX Prometheus Exporter metrics
      lowercaseOutputName: true
      rules:
      # Special cases and very specific rules
      - pattern: kafka.server<type=(.+), name=(.+), clientId=(.+), topic=(.+), partition=(.*)><>Value
        name: kafka_server_$1_$2
        type: GAUGE
        labels:
         clientId: "$3"
         topic: "$4"
         partition: "$5"
      - pattern: kafka.server<type=(.+), name=(.+), clientId=(.+), brokerHost=(.+), brokerPort=(.+)><>Value
        name: kafka_server_$1_$2
        type: GAUGE
        labels:
         clientId: "$3"
         broker: "$4:$5"
      - pattern: kafka.server<type=(.+), cipher=(.+), protocol=(.+), listener=(.+), networkProcessor=(.+)><>connections
        name: kafka_server_$1_connections_tls_info
        type: GAUGE
        labels:
          listener: "$2"
          networkProcessor: "$3"
          protocol: "$4"
          cipher: "$5"
      - pattern: kafka.server<type=(.+), clientSoftwareName=(.+), clientSoftwareVersion=(.+), listener=(.+), networkProcessor=(.+)><>connections
        name: kafka_server_$1_connections_software
        type: GAUGE
        labels:
          clientSoftwareName: "$2"
          clientSoftwareVersion: "$3"
          listener: "$4"
          networkProcessor: "$5"
      - pattern: "kafka.server<type=(.+), listener=(.+), networkProcessor=(.+)><>(.+):"
        name: kafka_server_$1_$4
        type: GAUGE
        labels:
         listener: "$2"
         networkProcessor: "$3"
      - pattern: kafka.server<type=(.+), listener=(.+), networkProcessor=(.+)><>(.+)
        name: kafka_server_$1_$4
        type: GAUGE
        labels:
         listener: "$2"
         networkProcessor: "$3"
      # Some percent metrics use MeanRate attribute
      # Ex) kafka.server<type=(KafkaRequestHandlerPool), name=(RequestHandlerAvgIdlePercent)><>MeanRate
      - pattern: kafka.(\w+)<type=(.+), name=(.+)Percent\w*><>MeanRate
        name: kafka_$1_$2_$3_percent
        type: GAUGE
      # Generic gauges for percents
      - pattern: kafka.(\w+)<type=(.+), name=(.+)Percent\w*><>Value
        name: kafka_$1_$2_$3_percent
        type: GAUGE
      - pattern: kafka.(\w+)<type=(.+), name=(.+)Percent\w*, (.+)=(.+)><>Value
        name: kafka_$1_$2_$3_percent
        type: GAUGE
        labels:
          "$4": "$5"
      # Generic per-second counters with 0-2 key/value pairs
      - pattern: kafka.(\w+)<type=(.+), name=(.+)PerSec\w*, (.+)=(.+), (.+)=(.+)><>Count
        name: kafka_$1_$2_$3_total
        type: COUNTER
        labels:
          "$4": "$5"
          "$6": "$7"
      - pattern: kafka.(\w+)<type=(.+), name=(.+)PerSec\w*, (.+)=(.+)><>Count
        name: kafka_$1_$2_$3_total
        type: COUNTER
        labels:
          "$4": "$5"
      - pattern: kafka.(\w+)<type=(.+), name=(.+)PerSec\w*><>Count
        name: kafka_$1_$2_$3_total
        type: COUNTER
      # Generic gauges with 0-2 key/value pairs
      - pattern: kafka.(\w+)<type=(.+), name=(.+), (.+)=(.+), (.+)=(.+)><>Value
        name: kafka_$1_$2_$3
        type: GAUGE
        labels:
          "$4": "$5"
          "$6": "$7"
      - pattern: kafka.(\w+)<type=(.+), name=(.+), (.+)=(.+)><>Value
        name: kafka_$1_$2_$3
        type: GAUGE
        labels:
          "$4": "$5"
      - pattern: kafka.(\w+)<type=(.+), name=(.+)><>Value
        name: kafka_$1_$2_$3
        type: GAUGE
      # Emulate Prometheus 'Summary' metrics for the exported 'Histogram's.
      # Note that these are missing the '_sum' metric!
      - pattern: kafka.(\w+)<type=(.+), name=(.+), (.+)=(.+), (.+)=(.+)><>Count
        name: kafka_$1_$2_$3_count
        type: COUNTER
        labels:
          "$4": "$5"
          "$6": "$7"
      - pattern: kafka.(\w+)<type=(.+), name=(.+), (.+)=(.*), (.+)=(.+)><>(\d+)thPercentile
        name: kafka_$1_$2_$3
        type: GAUGE
        labels:
          "$4": "$5"
          "$6": "$7"
          quantile: "0.$8"
      - pattern: kafka.(\w+)<type=(.+), name=(.+), (.+)=(.+)><>Count
        name: kafka_$1_$2_$3_count
        type: COUNTER
        labels:
          "$4": "$5"
      - pattern: kafka.(\w+)<type=(.+), name=(.+), (.+)=(.*)><>(\d+)thPercentile
        name: kafka_$1_$2_$3
        type: GAUGE
        labels:
          "$4": "$5"
          quantile: "0.$6"
      - pattern: kafka.(\w+)<type=(.+), name=(.+)><>Count
        name: kafka_$1_$2_$3_count
        type: COUNTER
      - pattern: kafka.(\w+)<type=(.+), name=(.+)><>(\d+)thPercentile
        name: kafka_$1_$2_$3
        type: GAUGE
        labels:
          quantile: "0.$4"
    zookeeper-metrics-config.yml: |
      # See https://github.com/prometheus/jmx_exporter for more info about JMX Prometheus Exporter metrics
      lowercaseOutputName: true
      rules:
      # replicated Zookeeper
      - pattern: "org.apache.ZooKeeperService<name0=ReplicatedServer_id(\\d+)><>(\\w+)"
        name: "zookeeper_$2"
        type: GAUGE
      - pattern: "org.apache.ZooKeeperService<name0=ReplicatedServer_id(\\d+), name1=replica.(\\d+)><>(\\w+)"
        name: "zookeeper_$3"
        type: GAUGE
        labels:
          replicaId: "$2"
      - pattern: "org.apache.ZooKeeperService<name0=ReplicatedServer_id(\\d+), name1=replica.(\\d+), name2=(\\w+)><>(Packets\\w+)"
        name: "zookeeper_$4"
        type: COUNTER
        labels:
          replicaId: "$2"
          memberType: "$3"
      - pattern: "org.apache.ZooKeeperService<name0=ReplicatedServer_id(\\d+), name1=replica.(\\d+), name2=(\\w+)><>(\\w+)"
        name: "zookeeper_$4"
        type: GAUGE
        labels:
          replicaId: "$2"
          memberType: "$3"
      - pattern: "org.apache.ZooKeeperService<name0=ReplicatedServer_id(\\d+), name1=replica.(\\d+), name2=(\\w+), name3=(\\w+)><>(\\w+)"
        name: "zookeeper_$4_$5"
        type: GAUGE
        labels:
          replicaId: "$2"
          memberType: "$3"
  EOF
  ```

2. 配置jvm_exporter和kafka_exporter  
  ```shell
  kubectl -n my-kafka-project patch kafka/my-cluster --type=merge -p $'
  spec: 
    kafka: 
      metricsConfig:
        type: jmxPrometheusExporter
        valueFrom:
          configMapKeyRef:
            name: kafka-metrics
            key: kafka-metrics-config.yml
    kafkaExporter:
      topicRegex: ".*"
      groupRegex: ".*"
  '
  ```
3. 配置Prometheus抓取指标
  ```shell
  cat <<EOF | kubectl apply -f -
  apiVersion: monitoring.coreos.com/v1
  kind: PodMonitor
  metadata:
    name: my-cluster-kafka-monitor
    namespace: my-kafka-project
    labels: 
      monitoring: kafka
  spec:
    jobLabel: kafka-monitor
    podMetricsEndpoints:
    - interval: 30s
      path: /metrics
      port: tcp-prometheus
    selector:
      matchLabels:
        strimzi.io/cluster: my-cluster
        strimzi.io/kind: Kafka
  EOF
  ```


# 测试

调整kafka和zookeeper的副本和资源配置  

```shell
kubectl -n my-kafka-project patch kafka/my-cluster --type=merge -p $'
spec: 
  kafka: 
    replicas: 3
    resources: 
      limits:
        cpu: '1'
        memory: 2Gi
  zookeeper: 
    replicas: 3
    resources: 
      limits:
        cpu: 300m
        memory: 500Mi
'
```

快速创建100个topic  

```shell
cat << \OEOF | bash -s
for ((a=1; a <= 100; a++))
do
cat <<EOF | kubectl apply -n my-kafka-project -f -
apiVersion: kafka.strimzi.io/v1beta2
kind: KafkaTopic
metadata:
  name: my-topic-$a
  labels:
    strimzi.io/cluster: "my-cluster"
spec:
  partitions: 40
  replicas: 2
EOF
sleep 10
done
OEOF
```

