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
kubectl run dapr-debug --image=curlimages/curl -it --rm --restart=Never -- /bin/sh 
curl -H "dapr-app-id: order-processor" http://order-processor-dapr:3500/v1.0/state/orders-store/cars
curl -H "dapr-app-id: order-processor" http://order-processor-dapr:3500/v1.0/workflows/dapr/OrderProcessingWorkflow/instances

curl -X POST http://order-processor-dapr:3500/v1.0/workflows/dapr/OrderProcessingWorkflow/start?instanceID=order-001 \
     -H "Content-Type: application/json" \
     -d '{
           "ItemName": "Kubernetes Cluster",
           "TotalAmount": 5000
         }'
```
