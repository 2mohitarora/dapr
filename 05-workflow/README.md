# Workflow

# Build order-processor service container images
```
export KO_DOCKER_REPO=ko.local
export DOCKER_HOST="unix:///Users/mua0008/.colima/default/docker.sock"
ko build -B ./order-processor
```

# Delete existing deployments
```
kubectl delete deployment order-processor
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
kubectl get deployments -l app=order-processor -o wide
```

# Check the pods in deployment 
```
kubectl get pods -l app=order-processor
```

# Check application logs
```
kubectl logs -l app=order-processor --prefix
```