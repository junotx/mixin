package rules

import (
	"context"
	"encoding/csv"
	"os"
	"strings"
	"testing"
	"text/template"

	"github.com/Masterminds/sprig"
	eventsv1alpha1 "github.com/kubesphere/kube-events/pkg/apis/v1alpha1"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	monitoringversioned "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

var builtinAlertRulesInfo = map[string]string{
	// general.rules
	"TargetDown": "Target服务的副本不可用率高",
	"Watchdog":   "",

	// kubernetes-system-scheduler
	"KubeSchedulerDown": "kube-scheduler不可用",

	// kubernetes-system-controller-manager
	"KubeControllerManagerDown": "kube-controller-manager不可用",

	// kubernetes-system-kubelet
	"KubeNodeNotReady":             "k8s节点长时间未就绪",
	"KubeNodeUnreachable":          "k8s节点不可达",
	"KubeletTooManyPods":           "节点的pod使用率高",
	"KubeNodeReadinessFlapping":    "节点就绪状态频繁变化",
	"KubeletPlegDurationHigh":      "kubelet的PLEG操作耗时长",
	"KubeletPodStartUpLatencyHigh": "kubelet启动pod时间长",
	"KubeletDown":                  "kubelet不可用",

	// kubernetes-system-apiserver
	"KubeAPILatencyHigh":              "KubeAPI资源请求延迟时间长",
	"KubeAPIErrorsHigh":               "KubeAPI资源请求异常率高",
	"KubeClientCertificateExpiration": "k8s客户端证书将过期",
	"AggregatedAPIErrors":             "AggregatedAPI异常，异常值高表示相关服务的可用性频繁切换",
	"AggregatedAPIDown":               "AggregatedAPI不可用",
	"KubeAPIDown":                     "KubeAPI不可用",

	// kube-apiserver-slos
	"KubeAPIErrorBudgetBurn": "kube-apiserver组件异常多",

	// kubernetes-apps
	"KubePodCrashLooping":               "容器组频繁重启",
	"KubePodNotReady":                   "容器组长时间未就绪",
	"KubeDeploymentGenerationMismatch":  "Deployment版本号不匹配",
	"KubeDeploymentReplicasMismatch":    "Deployment副本数不匹配",
	"KubeStatefulSetReplicasMismatch":   "StatefulSet副本数不匹配",
	"KubeStatefulSetGenerationMismatch": "StatefulSet版本号不匹配",
	"KubeStatefulSetUpdateNotRolledOut": "StatefulSet更新未被回滚",
	"KubeDaemonSetRolloutStuck":         "DaemonSet回滚阻塞",
	"KubeContainerWaiting":              "容器长时间处于等待状态",
	"KubeDaemonSetNotScheduled":         "DaemonSet的pod未调度",
	"KubeDaemonSetMisScheduled":         "DaemonSet的pod调度位置不对",
	"KubeCronJobRunning":                "CronJob完成任务耗时久",
	"KubeJobCompletion":                 "Job耗时久",
	"KubeJobFailed":                     "Job执行失败",
	"KubeHpaReplicasMismatch":           "HPA副本数不匹配",
	"KubeHpaMaxedOut":                   "HPA长时间处于最大副本状态",

	// kubernetes-resources
	"KubeCPUOvercommit":         "k8s集群CPU资源请求超额，将无法容忍节点故障",
	"KubeMemoryOvercommit":      "k8s集群内存资源请求超额，将无法容忍节点故障",
	"KubeCPUQuotaOvercommit":    "namespace的cpu资源请求超额",
	"KubeMemoryQuotaOvercommit": "namespace的内存资源请求超额",
	"KubeQuotaExceeded":         "namespace的资源用量高",
	"CPUThrottlingHigh":         "cpu处于节制状态时间占比高",

	// kubernetes-storage
	"KubePersistentVolumeFillingUp": "持久化存储卷空间即将用尽",
	"KubePersistentVolumeErrors":    "持久化存储卷状态异常",

	// kube-state-metrics
	"KubeStateMetricsListErrors":  "kube-state-metrics执行k8s资源的list操作异常，可能无法导出对应资源的指标数据",
	"KubeStateMetricsWatchErrors": "kube-state-metrics执行k8s资源的watch操作异常，可能无法导出对应资源的指标数据",

	// node-exporter
	"NodeFilesystemSpaceFillingUp":       "节点存储空间即将用尽",
	"NodeFilesystemAlmostOutOfSpace":     "节点存储空间几乎用尽",
	"NodeFilesystemFilesFillingUp":       "节点inodes即将用尽",
	"NodeFilesystemAlmostOutOfFiles":     "节点inodes几乎用尽",
	"NodeNetworkReceiveErrs":             "节点接收网络数据异常多",
	"NodeNetworkTransmitErrs":            "节点发送网络数据异常多",
	"NodeHighNumberConntrackEntriesUsed": "节点conntrack使用量接近限制",
	"NodeClockSkewDetected":              "节点时钟倾斜",

	// node-network
	"NodeNetworkInterfaceFlapping": "节点网络接口状态频繁变化",

	// prometheus
	"PrometheusBadConfig":                             "prometheus加载配置文件失败",
	"PrometheusNotificationQueueRunningFull":          "prometheus的告警通知队列将满",
	"PrometheusErrorSendingAlertsToSomeAlertmanagers": "prometheus发送告警到部分alertmanager实例出错",
	"PrometheusErrorSendingAlertsToAnyAlertmanager":   "prometheus发送告警到所有alertmanager实例出错",
	"PrometheusNotConnectedToAlertmanagers":           "prometheus未连接任何alertmanager",
	"PrometheusTSDBReloadsFailing":                    "prometheus加载磁盘块数据失败",
	"PrometheusTSDBCompactionsFailing":                "prometheus执行compact操作失败",
	"PrometheusNotIngestingSamples":                   "prometheus未摄入数据",
	"PrometheusDuplicateTimestamps":                   "prometheus摄入数据的时间戳重复，重复时间戳的数据将被丢弃",
	"PrometheusOutOfOrderTimestamps":                  "prometheus摄入数据的时间戳出现乱序，相应的数据将被丢弃",
	"PrometheusRemoteStorageFailures":                 "prometheus写远程数据失败",
	"PrometheusRemoteWriteBehind":                     "prometheus写远程数据滞后时间长",
	"PrometheusRemoteWriteDesiredShards":              "prometheus写远程需要更多shards。prometheus写远程时会启用多个shards并行写，当计算的最优shards数大于配置shards数时，会触发该告警",
	"PrometheusRuleFailures":                          "prometheus规则评估异常",
	"PrometheusMissingRuleEvaluations":                "prometheus错过规则评估，一般是由于规则评估过慢",

	// alertmanager
	"AlertmanagerConfigInconsistent":  "alertmanager配置不同步",
	"AlertmanagerFailedReload":        "alertmanager加载配置失败",
	"AlertmanagerMembersInconsistent": "alertmanager节点状态不一致，找不到集群内其他节点",

	// prometheus-operator
	"PrometheusOperatorReconcileErrors":  "prometheus-operator reconcile操作异常",
	"PrometheusOperatorNodeLookupErrors": "prometheus-operator reconcile prometheus异常",

	// other component-specific alert rules

	// etcd
	"etcdMembersDown":                    "etcd节点不可用",
	"etcdInsufficientMembers":            "etcd可用节点不足",
	"etcdNoLeader":                       "etcd没有leader节点",
	"etcdHighNumberOfLeaderChanges":      "etcd的leader节点频繁变更",
	"etcdHighNumberOfFailedGRPCRequests": "etcd的grpc请求失败率高",
	"etcdGRPCRequestsSlow":               "etcd处理GRPC请求慢",
	"etcdMemberCommunicationSlow":        "etcd节点间通信慢",
	"etcdHighNumberOfFailedProposals":    "etcd的proposal失败率高",
	"etcdHighFsyncDurations":             "etcd的fsync操作高延迟",
	"etcdHighCommitDurations":            "etcd的commit操作高延迟",
	"etcdHighNumberOfFailedHTTPRequests": "etcd的http请求失败率高",
	"etcdHTTPRequestsSlow":               "etcd处理http请求慢",

	// es
	"ElasticsearchFilesystemTooHigh":  "es存储用量高",
	"ElasticsearchTooFewNodesRunning": "es运行节点少",
	"ElasticsearchHeapTooHigh":        "es堆用量高",

	// coredns
	"CoreDNSDown":                           "coredns不可用",
	"CoreDNSLatencyHigh":                    "coredns响应延迟高",
	"CoreDNSErrorsHigh":                     "coredns响应异常率高",
	"CoreDNSForwardLatencyHigh":             "coredns转发延迟高",
	"CoreDNSForwardErrorsHigh":              "coredns转发异常率高",
	"CoreDNSForwardHealthcheckFailureCount": "coredns的上游服务健康检查失败",
	"CoreDNSForwardHealthcheckBrokenCount":  "coredns的上游服务健康检查全部失败",
}

type MetricAlertRule struct {
	Name     string
	Expr     string
	Severity string
	For      string
	Comment  string
	Supple   string
}

func parsePrometheusRule(pr *monitoringv1.PrometheusRule) (groups []string, metricRules map[string][]MetricAlertRule) {
	metricRules = make(map[string][]MetricAlertRule)

	type RuleSupple struct {
		Name     string
		Severity string
		For      string
		Supple   string
	}
	var ruleSupples = make(map[string][]RuleSupple)

	f, err := os.Open("metric_rules_duplicate.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var r = csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		panic(err)
	}
	for i, record := range records {
		if i == 0 || record[0] == "" {
			continue
		}
		ruleSupples[record[0]] = append(ruleSupples[record[0]], RuleSupple{
			Name:     record[0],
			Severity: record[1],
			For:      record[2],
			Supple:   record[3],
		})
	}

	if pr == nil {
		return
	}
	for _, g := range pr.Spec.Groups {
		groups = append(groups, g.Name)
		for _, r := range g.Rules {
			if r.Alert == "" {
				continue
			}
			rule := MetricAlertRule{
				Name:     r.Alert,
				Expr:     r.Expr.String(),
				Severity: r.Labels["severity"],
				For:      r.For,
				Comment:  builtinAlertRulesInfo[r.Alert],
			}
			if supples, ok := ruleSupples[rule.Name]; ok {
				for _, supple := range supples {
					if supple.Severity == rule.Severity && supple.For == rule.For {
						rule.Supple = supple.Supple
					}
				}
			}
			metricRules[g.Name] = append(metricRules[g.Name], rule)
		}
	}
	return
}

func TestGenMetricRulesDoc(t *testing.T) {

	var (
		templateFile = "metric_rules_doc_tmpl.md"
		outFile      = templateFile[0:strings.LastIndexByte(templateFile, '_')] + ".md"

		kubeconfigPath          = "D:/ks/conf/ks3-config"
		prometheusRuleNamespace = "kubesphere-monitoring-system"
		prometheusRuleName      = "prometheus-k8s-rules"

		etcdRuleNamespace = "kubesphere-monitoring-system"
		etcdRuleName      = "prometheus-k8s-etcd-rules"
	)

	conf, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		panic(err)
	}
	clientset := monitoringversioned.NewForConfigOrDie(conf)

	promRuleResource, err := clientset.MonitoringV1().PrometheusRules(prometheusRuleNamespace).
		Get(context.TODO(), prometheusRuleName, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	groups, rules := parsePrometheusRule(promRuleResource)

	etcdRuleResource, err := clientset.MonitoringV1().PrometheusRules(etcdRuleNamespace).
		Get(context.TODO(), etcdRuleName, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	etcdGroups, etcdRules := parsePrometheusRule(etcdRuleResource)

	groups = append(groups, etcdGroups...)
	for g, rs := range etcdRules {
		rules[g] = rs
	}

	tpl := template.Must(template.New(templateFile).Funcs(sprig.TxtFuncMap()).ParseFiles(templateFile))
	f, err := os.Create(outFile)
	if err != nil {
		panic(err)
	}
	err = tpl.Execute(f, struct {
		Groups []string
		Rules  map[string][]MetricAlertRule
	}{groups, rules})
	if err != nil {
		panic(err)
	}
}

type EventAlertRule struct {
	Name     string
	Expr     string
	Severity string
	Comment  string
}

func parseEventRule(er *eventsv1alpha1.Rule) []EventAlertRule {
	var rules []EventAlertRule
	if er == nil {
		return rules
	}
	for _, r := range er.Spec.Rules {
		if r.Type != "alert" {
			continue
		}
		rules = append(rules, EventAlertRule{
			Name:     r.Name,
			Expr:     r.Condition,
			Severity: r.Labels["severity"],
			Comment:  r.Annotations["summaryCn"],
		})
	}
	return rules
}

func TestGenEventRulesDoc(t *testing.T) {
	var (
		templateFile = "event_rules_doc_tmpl.md"
		outFile      = templateFile[0:strings.LastIndexByte(templateFile, '_')] + ".md"

		kubeconfigPath     = "D:/ks/conf/ks3-config"
		eventRuleNamespace = "kubesphere-logging-system"
		eventRuleName      = "ks-events-cluster-rules-default"
	)

	conf, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		panic(err)
	}
	clientset := dynamic.NewForConfigOrDie(conf)
	unstructed, err := clientset.Resource(
		schema.GroupVersionResource{Group: "events.kubesphere.io", Version: "v1alpha1", Resource: "rules"}).
		Namespace(eventRuleNamespace).Get(context.TODO(), eventRuleName, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	var eventRule eventsv1alpha1.Rule
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructed.Object, &eventRule)
	if err != nil {
		panic(err)
	}

	tpl := template.Must(template.New(templateFile).Funcs(sprig.TxtFuncMap()).ParseFiles(templateFile))
	f, err := os.Create(outFile)
	if err != nil {
		panic(err)
	}
	err = tpl.Execute(f, struct {
		Rules []EventAlertRule
	}{parseEventRule(&eventRule)})
	if err != nil {
		panic(err)
	}
}

func wrapForce(str string, wrapLength int, newLineStr string) string {
	var ret strings.Builder
	var l int
	_ = strings.FieldsFunc(str, func(r rune) bool {
		if l >= wrapLength {
			ret.WriteRune(r)
			ret.WriteString(newLineStr)
			l = 0
		} else {
			rl, _ := ret.WriteRune(r)
			l += rl
		}
		return true
	})

	return ret.String()
}
