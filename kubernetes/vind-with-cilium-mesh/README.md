# Install Tools
```
brew install kubectl helm docker go ko dapr/tap/dapr-cli cilium-cli vcluster
```

# Install Orbstack
```
brew install --cask orbstack

# Add local registries that will be created later
# Add registry to docker daemon in ~/.docker/daemon.json
{
  "insecure-registries": ["localhost:5050", "localhost:5051"]
}

# Start Orbstack
```

# Configure docker
```
docker context use orbstack
export DOCKER_HOST="unix:///Users/mua0008/.orbstack/run/docker.sock"
docker context list
```

# Create your first vcluster
```
vcluster create cluster-1 --driver docker --values cluster-1.yaml

# Add the Jetstack repo
helm repo add cilium https://helm.cilium.io/
helm repo add jetstack https://charts.jetstack.io
helm repo update

# Upgrade Coredns
kubectl set image deployment/coredns \
  -n kube-system \
  coredns=registry.k8s.io/coredns/coredns:v1.14.2

# Install cilium 

helm install cilium cilium/cilium --version 1.19.1 --set kubeProxyReplacement=true --namespace cilium --create-namespace --set ipam.operator.clusterPoolIPv4PodCIDRList=10.1.0.0/16 --set routingMode=tunnel  --set tunnelProtocol=vxlan --set ipam.mode=cluster-pool

# After CNI is installed, wait for pods to become Ready:
kubectl get pods --all-namespaces -w

# Install cert-manager
helm install cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --create-namespace \
  --set installCRDs=true

# After cert-manager, wait for pods to become Ready:
kubectl get pods --all-namespaces -w

kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.2.0/experimental-install.yaml

kubectl apply -f cert-manager-issuer.yaml

kubectl get clusterissuer cilium-ca-issuer 
kubectl get certificates -n cilium 

helm upgrade cilium cilium/cilium --version 1.19.1 \
  --namespace cilium \
  -f cilium-1-helm.yaml \

# After Mesh is installed, wait for pods to become Ready:
kubectl get pods --all-namespaces -w

# Check cilium status
cilium status --namespace cilium

# Note: Make sure to configure the CNI plugin according to your cluster's pod CIDR
kubectl get configmap cilium-config -n cilium -o yaml | grep -i cidr
```

# Checks
```
# Multi-cluster connectivity check
# Run from your management terminal
cilium connectivity test \
  --context cluster-east \
  --multi-cluster cluster-west \
  --test pod-to-pod,pod-to-service

# The MCS-API DNS Validator

```

# Check Gateway Class and Create Cilium Gateway
```
kubectl get gatewayclasses -o wide

kubectl apply -f cilium-gateway.yaml

kubectl get gateways -n cilium

# Debug Cilium Gateway
# See the service Cilium created
kubectl get svc -l io.cilium.gateway/owning-gateway=default-gateway -n cilium
# See the Cilium Envoy proxy pod
kubectl -n kube-system logs -l app.kubernetes.io/name=cilium-envoy -f -n cilium

# Create a route for front end service
kubectl apply -f ./frontendsvc-route.yaml
```

# Add Traefik Ingress for Dapr communication between clusters
```
helm repo add traefik https://traefik.github.io/charts
helm repo update
helm install traefik traefik/traefik \
  --namespace traefik \
  --create-namespace \
  --set providers.kubernetesCRD.allowExternalNameServices=true --skip-crds

# Install CRDs manually
kubectl apply -f https://raw.githubusercontent.com/traefik/traefik-helm-chart/master/traefik/crds/traefik.io_ingressroutes.yaml
kubectl apply -f https://raw.githubusercontent.com/traefik/traefik-helm-chart/master/traefik/crds/traefik.io_ingressroutetcps.yaml
kubectl apply -f https://raw.githubusercontent.com/traefik/traefik-helm-chart/master/traefik/crds/traefik.io_ingressrouteudps.yaml
kubectl apply -f https://raw.githubusercontent.com/traefik/traefik-helm-chart/master/traefik/crds/traefik.io_middlewares.yaml
kubectl apply -f https://raw.githubusercontent.com/traefik/traefik-helm-chart/master/traefik/crds/traefik.io_middlewaretcps.yaml
kubectl apply -f https://raw.githubusercontent.com/traefik/traefik-helm-chart/master/traefik/crds/traefik.io_serverstransports.yaml
kubectl apply -f https://raw.githubusercontent.com/traefik/traefik-helm-chart/master/traefik/crds/traefik.io_serverstransporttcps.yaml
kubectl apply -f https://raw.githubusercontent.com/traefik/traefik-helm-chart/master/traefik/crds/traefik.io_tlsoptions.yaml
kubectl apply -f https://raw.githubusercontent.com/traefik/traefik-helm-chart/master/traefik/crds/traefik.io_tlsstores.yaml
kubectl apply -f https://raw.githubusercontent.com/traefik/traefik-helm-chart/master/traefik/crds/traefik.io_traefikservices.yaml

# Create Dapr IngressRoute
kubectl apply -f traefik-dapr-ingress.yaml
```

# Install Dapr Inject Dapr sidecar into Traefik deployment
```
dapr init -k

dapr status -k

# Make sure Dapr containers are ready
kubectl get pods -o wide -n dapr-system

# Run DAPR UI

dapr dashboard -k

# Patch Traefik deployment to inject Dapr sidecar

kubectl patch deployment traefik -n traefik -p '
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

# Check traffic service created by Dapr
kubectl get svc -n traefik | grep traefik-ingress-dapr
```

# Configure Registry for first cluster
```
# Start a local registry on the same Docker network as your vind cluster
docker run -d --name registry-1 --network vind-cluster-1 -p 5050:5000 registry:2

# Configure registry for cluster-1 so that nodes can pull from insecure registry
./cluster-1-script.sh
```

# Disconnect from first cluster
```
vcluster disconnect
```

# Create second vcluster
```
vcluster create cluster-2 --driver docker --values cluster-2.yaml

helm repo add cilium https://helm.cilium.io/

kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.2.0/experimental-install.yaml

helm install cilium cilium/cilium --version 1.19.1 --set kubeProxyReplacement=true --set gatewayAPI.enabled=true --namespace cilium --create-namespace --set ipam.operator.clusterPoolIPv4PodCIDRList=10.2.0.0/16

# After CNI is installed, wait for pods to become Ready:
kubectl get pods --all-namespaces -w

# Check cilium status
cilium status --namespace cilium

# Note: Make sure to configure the CNI plugin according to your cluster's pod CIDR
kubectl get configmap cilium-config -n cilium -o yaml | grep -i cidr
```

# Check Gateway Class and Create Cilium Gateway
```
kubectl get gatewayclasses -o wide

kubectl apply -f cilium-gateway.yaml

kubectl get gateways -n cilium

# Debug Cilium Gateway
# See the service Cilium created
kubectl get svc -l io.cilium.gateway/owning-gateway=default-gateway -n cilium
# See the Cilium Envoy proxy pod
kubectl -n kube-system logs -l app.kubernetes.io/name=cilium-envoy -f -n cilium
```

# Add Traefik Ingress for Dapr communication between clusters
```
helm repo add traefik https://traefik.github.io/charts
helm repo update
helm install traefik traefik/traefik \
  --namespace traefik \
  --create-namespace \
  --set providers.kubernetesCRD.allowExternalNameServices=true --skip-crds

# Install CRDs manually
kubectl apply -f https://raw.githubusercontent.com/traefik/traefik-helm-chart/master/traefik/crds/traefik.io_ingressroutes.yaml
kubectl apply -f https://raw.githubusercontent.com/traefik/traefik-helm-chart/master/traefik/crds/traefik.io_ingressroutetcps.yaml
kubectl apply -f https://raw.githubusercontent.com/traefik/traefik-helm-chart/master/traefik/crds/traefik.io_ingressrouteudps.yaml
kubectl apply -f https://raw.githubusercontent.com/traefik/traefik-helm-chart/master/traefik/crds/traefik.io_middlewares.yaml
kubectl apply -f https://raw.githubusercontent.com/traefik/traefik-helm-chart/master/traefik/crds/traefik.io_middlewaretcps.yaml
kubectl apply -f https://raw.githubusercontent.com/traefik/traefik-helm-chart/master/traefik/crds/traefik.io_serverstransports.yaml
kubectl apply -f https://raw.githubusercontent.com/traefik/traefik-helm-chart/master/traefik/crds/traefik.io_serverstransporttcps.yaml
kubectl apply -f https://raw.githubusercontent.com/traefik/traefik-helm-chart/master/traefik/crds/traefik.io_tlsoptions.yaml
kubectl apply -f https://raw.githubusercontent.com/traefik/traefik-helm-chart/master/traefik/crds/traefik.io_tlsstores.yaml
kubectl apply -f https://raw.githubusercontent.com/traefik/traefik-helm-chart/master/traefik/crds/traefik.io_traefikservices.yaml

# Create Dapr IngressRoute
kubectl apply -f traefik-dapr-ingress.yaml
```

# Install Dapr Inject Dapr sidecar into Traefik deployment
```
dapr init -k

dapr status -k

# Make sure Dapr containers are ready
kubectl get pods -o wide -n dapr-system

# Run DAPR UI

dapr dashboard -k -p 9999

# Patch Traefik deployment to inject Dapr sidecar

kubectl patch deployment traefik -n traefik -p '
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

# Check traffic service created by Dapr
kubectl get svc -n traefik | grep traefik-ingress-dapr
```

# Configure Registry for second cluster
```
# Start a local registry on the same Docker network as your vind cluster
docker run -d --name registry-2 --network vind-cluster-2 -p 5051:5000 registry:2

# Configure registry for cluster-2 so that nodes can pull from insecure registry
./cluster-2-script.sh
```

# Verify both clusters are running
```
vcluster list
```

# Check all the docker containers that are running
```
docker ps --format "table {{.Names}}\t{{.Status}}"
```

# Describe clusters
```
vcluster describe cluster-1
vcluster describe cluster-1
```

# Cleanup
```
vcluster delete cluster-1
vcluster delete cluster-2
```

# Sample command for logs
```
# View control plane logs
docker exec vcluster.cp.cluster-1 journalctl -u vcluster --no-pager

# View worker node kubelet logs
docker exec vcluster.node.cluster-1.worker-1 journalctl -u kubelet --no-pager
```

# Networking
- Each cluster uses a separate Docker network (vind-cluster-1 and vind-cluster-2) to keep them isolated from each other.
- LoadBalancer is enabled for both clusters
