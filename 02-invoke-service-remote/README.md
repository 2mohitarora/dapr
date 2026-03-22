# Service Invocation in different cluster

# Context Checks
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
kubectl apply -f ./manifest/gateway.yaml

# Check the gateway status
kubectl get gateway default-gateway

# See the service Cilium created
kubectl get svc -l io.cilium.gateway/owning-gateway=default-gateway

# See the Envoy proxy pod
kubectl -n kube-system logs -l app.kubernetes.io/name=cilium-envoy -f
```
# Discuss configuration in detail (rate limiting is just an example)

# Invoke genid service by find Traefik external IP
```
curl -H "Host: genidsvc.ingress" http://192.168.64.4/v1.0/invoke/genidsvc.default/method/genid -X POST
```

# Port forwarding from mac
```
kubectl port-forward -n kube-system service/traefik -n kube-system 8083:80
curl -H "Host: genidsvc.ingress" http://localhost:8083/v1.0/invoke/genidsvc.default/method/genid -X POST
```

# From other colima cluster
```
kubectl run net-test --rm -it --image=nicolaka/netshoot -- /bin/bash
curl -H "Host: genidsvc.ingress" http://192.168.5.2:8083/v1.0/invoke/genidsvc.default/method/genid -X POST
```

# Check Rate Limit
```
# Uncomment dapr.io/config: "genidsvc-config" in genid.yaml
kubectl apply -f ./manifest/genid.yaml
kubectl apply -f ./manifest/rate-limit.yaml

seq 1 20 | xargs -P 20 -I {} curl -s -o /dev/null -w "%{http_code}\n" \
 http://192.168.64.4/v1.0/invoke/genidsvc.default/method/genid -X POST
```