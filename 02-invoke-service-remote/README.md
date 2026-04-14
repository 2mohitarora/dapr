# cluster-1
```
kubectl delete deployment frontendsvc
kubectl delete deployment genidsvc
kubectl apply -f ./manifest/cluster-1/remote-genid-endpoint.yaml
kubectl apply -f ./manifest/cluster-1/frontend.yaml
```

# Switch to cluster-2 and make sure you are using cluster-2
```
docker context ls
kubectl config get-contexts
```

# Delete existing deployments
```
kubectl delete deployment genidsvc
```

# Deploy manifests
```
kubectl apply -f ./manifest/cluster-2/genid.yaml
```

# Invoke genid service by find Traefik external IP
```
curl -v -H "Host: genidsvc.ingress" http://192.168.117.253/v1.0/invoke/genidsvc.default/method/genid -X POST

# Test front end application
curl -i -d '{ "items": ["automobile"]}'  -H "Content-type: application/json" "http://192.168.107.254/orders/new"

curl -i  -H "Content-type: application/json" "http://192.168.107.254/orders/order/order-04233f93-028a-41b7-ac84-44fb3b754599"
```

# Discuss rate limit and Check Rate Limit
```
# Uncomment dapr.io/config: "genidsvc-config" in genid.yaml
kubectl apply -f ./manifest/cluster-2/ratelimit.yaml
kubectl delete deployment genidsvc
kubectl apply -f ./manifest/cluster-2/genid.yaml

seq 1 20 | xargs -P 20 -I {} curl -s -o /dev/null -w "%{http_code}\n" \
 http://192.168.117.253/v1.0/invoke/genidsvc.default/method/genid -X POST

```