# Workflow

# Build order-processor service container images
```
ko build -L -B ./order-processor --platform=linux/arm64
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

dapr workflow history 43efc694-37d3-4714-9f7a-e3870af6898e -k --app-id order-processor -o wide
```
