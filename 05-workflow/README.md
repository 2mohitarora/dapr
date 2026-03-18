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

# Check workflow status
```
kubectl port-forward svc/redis-master 6379:6379

export REDIS_PASSWORD=$(kubectl get secret --namespace default redis -o jsonpath="{.data.redis-password}" | base64 --decode)

dapr workflow list -k --app-id order-processor \
  --connection-string "redis://:$REDIS_PASSWORD@localhost:6379" \
  -o wide

dapr workflow history 3bf1d79f-6bf6-489d-8ae9-1f3f1182d8f6 -k --app-id order-processor -o wide

dapr workflow run OrderProcessingWorkflow -k --app-id order-processor -i order-001 -x "{\"ItemName\":\"Kubernetes Cluster\",\"Quantity\":1,\"TotalAmount\":5000}"

```
