# k8s-demo application for K8s API development

* Preparation/Local run:
```
cd app/
go mod download

# * Execution:

go run kubeClient.go
```

* Build a Docker Image:
```
docker build -t k8s-demo .
```

* Apply the helm chart
```
helm3 upgrade --install myrelease ./k8s-demo --set image.repository=k8s-demo
```