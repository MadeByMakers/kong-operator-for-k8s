module github.com/MadeByMakers/kong-operator-for-k8s

go 1.16

require (
	github.com/go-logr/logr v0.4.0
	github.com/google/uuid v1.1.2 // indirect
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.17.0
	github.com/prometheus/common v0.26.0 // indirect
	github.com/redhat-cop/operator-utils v1.3.2
	k8s.io/api v0.22.1
	k8s.io/apiextensions-apiserver v0.22.1 // indirect
	k8s.io/apimachinery v0.22.1
	k8s.io/client-go v0.22.1
	sigs.k8s.io/controller-runtime v0.10.0
)
