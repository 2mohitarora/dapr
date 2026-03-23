# Service Invocation

# Install redis on k8s
```
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
helm install redis bitnami/redis
```

# Build front end service and genid service container images
```
export KO_DOCKER_REPO=localhost:5050
export DOCKER_HOST="unix:///Users/mua0008/.orbstack/run/docker.sock"
ko build -B ./frontendsvc --platform=linux/arm64
ko build -B ./genidsvc --platform=linux/arm64
```

# Delete existing deployments
```
kubectl delete deployment frontendsvc
kubectl delete deployment genidsvc
```

# Deploy manifests
```
kubectl apply -f ./manifest/redis-store.yaml
kubectl apply -f ./manifest/frontend.yaml
kubectl apply -f ./manifest/genid.yaml
kubectl apply -f ./manifest/resiliency.yaml
```

# Deploy Additional manifests to showcase remote service invoke and resiliency
```
kubectl apply -f ./manifest-ext/remote-genid-endpoint.yaml
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

# Test the application
```
curl -i -d '{ "items": ["automobile"]}'  -H "Content-type: application/json" "http://192.168.97.254/orders/new"

curl -i  -H "Content-type: application/json" "http://192.168.97.254/orders/order/order-e46d35bf-e088-4587-9818-d6af68ce81d2"
```

# Check application logs
```
kubectl logs -l app=frontendsvc --prefix
```