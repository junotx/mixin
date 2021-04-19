module github.com/junotx/mixin/se

go 1.15

require (
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Masterminds/sprig v2.22.0+incompatible
	github.com/alecthomas/chroma v0.8.2
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/huandu/xstrings v1.3.1 // indirect
	github.com/kubesphere/kube-events v0.2.0
	github.com/mitchellh/copystructure v1.0.0 // indirect
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.49.0
	github.com/prometheus-operator/prometheus-operator/pkg/client v0.49.0
	golang.org/x/net v0.0.0-20210610132358-84b48f89b13b // indirect
	golang.org/x/oauth2 v0.0.0-20210514164344-f6687ab2804c // indirect
	k8s.io/apimachinery v0.21.3
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/klog/v2 v2.9.0 // indirect
)

replace (
	k8s.io/client-go => k8s.io/client-go v0.21.2
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.9.0
)
