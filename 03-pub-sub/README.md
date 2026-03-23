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
ko build -B ./orderprocsvc --platform=linux/arm64
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
kubectl get deployments -l app=orderprocsvc -o wide
```

# Check the pods in deployment 
```
kubectl get pods -l app=frontendsvc
kubectl get pods -l app=orderprocsvc
```

# Test the application
```
curl -i -d '{ "items": ["bike"]}'  -H "Content-type: application/json" "http://192.168.97.254/orders/new"

curl -i  -H "Content-type: application/json" "http://192.168.97.254/orders/order/order-1e33be7b-2392-4a08-9901-01b59289895c"
```

# Check application logs
```
kubectl logs -l app=frontendsvc --prefix
kubectl logs -l app=orderprocsvc --prefix
```

# Debug
```
kubectl apply -f ./manifest/redis-cli.yaml
kubectl get secret --namespace default redis -o jsonpath="{.data.redis-password}" | base64 --decode
kubectl exec -it redis-client -- /bin/sh
redis-cli -h redis-master.default.svc.cluster.local -p 6379 -a '<password>'
keys *
get order-56bd2207-f333-4fa5-a5f4-afd6104b2df8
TYPE received-orders
XRANGE received-orders - +
```

# Move from Redis to RabbitMQ without code change. First lets install rabbitmq
```
kubectl apply -f ./manifest-ext/rabbitmq.yaml
```

# Access RabbitMQ management interface
```
kubectl port-forward --namespace default svc/rabbitmq 15672:15672
http://127.0.0.1:15672/
```

# Delete existing deployments
```
kubectl delete deployment frontendsvc
kubectl delete deployment orderprocsvc
```

# Delete redis strams pub sub Dapr component and create rabbitmq pubsub Dapr component
```
kubectl delete component orders-pubsub
kubectl apply -f ./manifest-ext/rabbitmq-pubsub.yaml
Check Dapr UI
```

# Restart applications for new components to be picked up 
```
kubectl apply -f ./manifest/orderproc.yaml
kubectl apply -f ./manifest/frontend.yaml
```

# Moment of truth
curl -i -d '{ "items": ["bike"]}'  -H "Content-type: application/json" "http://192.168.97.254/orders/new"

curl -i  -H "Content-type: application/json" "http://192.168.97.254/orders/order/order-da56bb75-7ca8-4e64-a7d1-bebb48321bf6"
```
