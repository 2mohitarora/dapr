# Service Invocation in different cluster

# Context Checks
```
docker context ls
kubectl config get-contexts
```

# Build front end service and genid service container images
```
export KO_DOCKER_REPO=ko.local
export DOCKER_HOST="unix:///Users/mua0008/.colima/ingress/docker.sock"
ko build -B ./genidsvc
```

# Delete existing deployments
```
kubectl delete deployment genidsvc
```

# Deploy manifests
```
kubectl apply -f ./manifest
```
# Discuss configuration in detail (rate limiting is just an example)

# Inject Dapr sidecar into Traefik deployment

```
kubectl patch deployment traefik -n kube-system -p '
{
  "spec": {
    "template": {
      "metadata": {
        "annotations": {
          "dapr.io/enabled": "true",
          "dapr.io/app-id": "traefik-ingress",
          "dapr.io/app-port": "8000",
          "dapr.io/log-level": "debug"
        }
      }
    }
  }
}'
```

# Check traffic service created by Dapr
```
kubectl get svc -n kube-system | grep traefik-ingress-dapr
```

# Invoke genid service by find Traefik external IP
```
kubectl get svc -n kube-system -l app.kubernetes.io/name=traefik
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
seq 1 20 | xargs -P 20 -I {} curl -s -o /dev/null -w "%{http_code}\n" \
 http://192.168.64.4/v1.0/invoke/genidsvc.default/method/genid -X POST
```