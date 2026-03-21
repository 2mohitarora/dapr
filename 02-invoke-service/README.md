# Service Invocation

# Install redis on k8s
```
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
helm install redis bitnami/redis
```

# Build front end service and genid service container images
```
export KO_DOCKER_REPO=ko.local
export DOCKER_HOST="unix:///Users/mua0008/.colima/default/docker.sock"
ko build -B ./frontendsvc
ko build -B ./genidsvc
```

# Delete existing deployments
```
kubectl delete deployment frontendsvc
kubectl delete deployment genidsvc
```

# Deploy manifests
```
kubectl apply -f ./manifest
```

# Deploy Additional manifests to showcase remote service invoke and resiliency
```
kubectl apply -f ./manifest-ext
```

# Discuss resiliency in detail

# Check Dapr UI 

# Check dapr components
```
dapr components -k
```

# Check k8s deployment 
```
kubectl get deployments -l app=frontendsvc -o wide
kubectl get deployments -l app=genidsvc -o wide
```

# Check the pods in deployment 
```
kubectl get pods -l app=frontendsvc
kubectl get pods -l app=genidsvc
```

# Test application using port forwarding 
```
kubectl port-forward deployment/frontendsvc 8081:8080
```

# Test the application
```
curl -i -d '{ "items": ["automobile"]}'  -H "Content-type: application/json" "http://localhost:8081/orders/new"

curl -i  -H "Content-type: application/json" "http://localhost:8081/orders/order/order-36a99c85-71dd-49f6-94b3-4c1807b850a8"
```

# Check application logs
```
kubectl logs -l app=frontendsvc --prefix
```