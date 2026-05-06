# Service Invocation

# Install redis on k8s
```
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
helm install redis bitnami/redis
```

# Build front end service and genid service container images
```
ko build -B -L ./frontendsvc --platform=linux/arm64
ko build -B -L ./orderprocsvc --platform=linux/arm64
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
curl -i -d '{ "items": ["bike"]}'  -H "Content-type: application/json" "http://192.168.147.254/orders/new"

curl -i  -H "Content-type: application/json" "http://192.168.147.254/orders/order/order-16f8e4b6-b91d-4292-a78c-cf501e660a40"
```

# Check application logs
```
kubectl logs -l app=frontendsvc --prefix
kubectl logs -l app=orderprocsvc --prefix
```

# Debug
```
kubectl get secret --namespace default redis -o jsonpath="{.data.redis-password}" | base64 --decode
kubectl exec -it redis-client -- /bin/sh
redis-cli -h redis-master.default.svc.cluster.local -p 6379 -a '<password>'
keys *
HGETALL order-16f8e4b6-b91d-4292-a78c-cf501e660a40
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
dapr components -k
```

# Restart applications for new components to be picked up 
```
kubectl apply -f ./manifest/orderproc.yaml
kubectl apply -f ./manifest/frontend.yaml
```

# Moment of truth
curl -i -d '{ "items": ["bike"]}'  -H "Content-type: application/json" "http://192.168.147.254/orders/new"

curl -i  -H "Content-type: application/json" "http://192.168.147.254/orders/order/order-5602c0f9-d276-4bc7-b042-ac6e666dd943"
```
