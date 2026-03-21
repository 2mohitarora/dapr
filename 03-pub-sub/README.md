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
ko build -B ./orderprocsvc
```

# Delete existing deployments
```
kubectl delete deployment frontendsvc
kubectl delete deployment orderprocsvc
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
kubectl get deployments -l app=genidsvc -o wide
kubectl get deployments -l app=orderprocsvc -o wide
```

# Check the pods in deployment 
```
kubectl get pods -l app=frontendsvc
kubectl get pods -l app=genidsvc
kubectl get pods -l app=orderprocsvc
```

# Test application using port forwarding 
```
kubectl port-forward deployment/frontendsvc 8081:8080
```

# Test the application
```
curl -i -d '{ "items": ["bike"]}'  -H "Content-type: application/json" "http://localhost:8081/orders/new"

curl -i  -H "Content-type: application/json" "http://localhost:8081/orders/order/order-56bd2207-f333-4fa5-a5f4-afd6104b2df8"
```

# Check application logs
```
kubectl logs -l app=frontendsvc --prefix
kubectl logs -l app=orderprocsvc --prefix
```

# Debug
```
kubectl apply -f redis-cli.yaml
kubectl exec -it redis-client -- /bin/sh
redis-cli -h redis-master.default.svc.cluster.local -p 6379 -a '<password>'
keys *
get order-56bd2207-f333-4fa5-a5f4-afd6104b2df8
```

# Move from Redis to RabbitMQ without code change
```
kubectl apply -f ./manifest-ext/rabbitmq.yaml
```

# Access RabbitMQ management interface
```
kubectl port-forward --namespace default svc/rabbitmq 15672:15672
http://127.0.0.1:15672/
```

