# State Management

# Install redis on k8s
```
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
helm install redis bitnami/redis
```

# Build front end service container image
```
export KO_DOCKER_REPO=localhost:5050
export DOCKER_HOST="unix:///Users/mua0008/.orbstack/run/docker.sock"
ko build -B ./frontendsvc --platform=linux/arm64
```

# Deploy manifests
```
kubectl apply -f ./manifest
```

# Check Dapr UI 

# Check dapr components
```
dapr components -k
```

# Check k8s deployment 
```
kubectl get deployments -l app=frontendsvc -o wide
```

# Check the pods in deployment 
```
kubectl get pods -l app=frontendsvc
--2 containers running in the application pod (one for our application and the other for the Dapr sidecar)
kubectl describe pods -l app=frontendsvc
```

# Test application using port forwarding 
```
kubectl port-forward deployment/frontendsvc 8081:8080
```

# Test the application
```
curl -i -d '{ "items": ["automobile"]}'  -H "Content-type: application/json" "http://localhost:8081/orders/new"

curl -i  -H "Content-type: application/json" "http://localhost:8081/orders/order/order-29e17d22"
```

# Check application logs
```
kubectl logs -l app=frontendsvc --prefix
```