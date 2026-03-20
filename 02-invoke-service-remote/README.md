# Service Invocation in different cluster

# Create a new kubernetes cluster with ingress enabled
```
cd ~/.colima/_templates
cp default.yaml ingress.yaml 
```

# Edit ingress.yaml
```
Comment following lines
k3sArgs:
    - --disable=traefik
```

# Start colima with ingress enabled
```
colima start --kubernetes --k3s-arg "" -p ingress
```

# Check kuberneres contexts
```
kubectl config get-contexts
```

# Check traefik is installed
```
kubectl get pods -n kube-system -l app.kubernetes.io/name=traefik
```

# Install Dapr on the new cluster
```
dapr init -k
```

# Verify Dapr
```
dapr status -k
kubectl get pods -o wide -n dapr-system
```

# Docker checks
```
docker context ls
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

# Run dapr UI and check genid svc is running
```
dapr dashboard -k -p 9999
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