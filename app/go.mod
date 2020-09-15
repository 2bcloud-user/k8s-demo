module k8s-demo

go 1.15

require (
	github.com/imdario/mergo v0.3.11 // indirect
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e // indirect
	k8s.io/api v0.0.0-20200726131424-9540e4cac147
	k8s.io/apimachinery v0.0.0-20200726131235-945d4ebf362b
	k8s.io/client-go v0.0.0-20200726131703-36233866f1c7
	k8s.io/klog v1.0.0 // indirect
	k8s.io/klog/v2 v2.2.0
	k8s.io/utils v0.0.0-20200720150651-0bdb4ca86cbc
)

replace (
	k8s.io/api => k8s.io/api v0.0.0-20200726131424-9540e4cac147
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20200726131235-945d4ebf362b
)
