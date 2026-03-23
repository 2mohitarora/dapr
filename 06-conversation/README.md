# Conversation Service

# Build convsvc service container images
```
export KO_DOCKER_REPO=ko.local
export DOCKER_HOST="unix:///Users/mua0008/.colima/default/docker.sock"
ko build -B ./convsvc --platform=linux/arm64
```

# Delete existing deployments
```
kubectl delete deployment convsvc
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
kubectl get deployments -l app=convsvc -o wide
```

# Check the pods in deployment 
```
kubectl get pods -l app=convsvc
```

# Check application logs
```
kubectl logs -l app=convsvc --prefix
```
