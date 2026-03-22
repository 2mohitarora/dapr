# Workflow

# Build order-processor service container images
```
export KO_DOCKER_REPO=localhost:5050
export DOCKER_HOST="unix:///Users/mua0008/.orbstack/run/docker.sock"
ko build -B ./order-processor --platform=linux/arm64
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

dapr workflow history 7cab37fd-a76e-41ea-8601-ae73572a6c57 -k --app-id order-processor -o wide

dapr workflow run OrderProcessingWorkflow -k \
  --app-id order-processor \
  --instance-id "order-004" \
  --input '{"item_name":"Kubernetes Cluster","total_cost":5000,"quantity":1}'

dapr workflow history order-004 -k --app-id order-processor -o wide

```
