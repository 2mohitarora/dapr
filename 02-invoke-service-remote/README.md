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

# Build front end service and genid service container images
```
export KO_DOCKER_REPO=localhost:5051
export DOCKER_HOST="unix:///Users/mua0008/.orbstack/run/docker.sock"
ko build -B ./genidsvc --platform=linux/arm64
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
curl -H "Host: genidsvc.ingress" http://192.168.107.253/v1.0/invoke/genidsvc.default/method/genid -X POST
```

# Discuss rate limit and Check Rate Limit
```
# Uncomment dapr.io/config: "genidsvc-config" in genid.yaml
kubectl apply -f ./manifest/cluster-2/ratelimit.yaml
kubectl delete deployment genidsvc
kubectl apply -f ./manifest/cluster-2/genid.yaml

seq 1 20 | xargs -P 20 -I {} curl -s -o /dev/null -w "%{http_code}\n" \
 http://192.168.107.253/v1.0/invoke/genidsvc.default/method/genid -X POST

```