# cluster-1
```
kubectl apply -f ./manifest/cluster-1/remote-genid-endpoint.yaml
```

# Make sure you are using cluster-2
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
kubectl apply -f ./manifest/genid.yaml
```
# Discuss configuration in detail (rate limiting is just an example)

# Invoke genid service by find Traefik external IP
```
curl -H "Host: genidsvc.ingress" http://192.168.107.253/v1.0/invoke/genidsvc.default/method/genid -X POST
```

# Check Rate Limit
```
# Uncomment dapr.io/config: "genidsvc-config" in genid.yaml
kubectl apply -f ./manifest/ratelimit.yaml
kubectl delete deployment genidsvc
kubectl apply -f ./manifest/genid.yaml

seq 1 20 | xargs -P 20 -I {} curl -s -o /dev/null -w "%{http_code}\n" \
 http://192.168.107.253/v1.0/invoke/genidsvc.default/method/genid -X POST

```