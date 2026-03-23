# Install kubectl
```
brew install kubectl
```
# Install helm
```
brew install helm
```
# Install docker
```
brew install docker
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
# Install Cilium CLI
```
brew install cilium-cli
```
# Configure docker
```
docker context use orbstack
export DOCKER_HOST="unix:///Users/mua0008/.orbstack/run/docker.sock"
docker context list
```
# Install vcluster CLI
```
brew install vcluster
vcluster version
```
# Create your first vcluster
```
vcluster create cluster-1 --driver docker --values cluster-1.yaml

helm repo add cilium https://helm.cilium.io/

kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.2.0/experimental-install.yaml

helm install cilium cilium/cilium --version 1.19.1 --set kubeProxyReplacement=true --set gatewayAPI.enabled=true --namespace kube-system --set ipam.operator.clusterPoolIPv4PodCIDRList=10.1.0.0/16

# After CNI is installed, wait for pods to become Ready:
kubectl get pods --all-namespaces -w

# Check cilium status
cilium status --namespace kube-system

# Note: Make sure to configure the CNI plugin according to your cluster's pod CIDR
kubectl get configmap cilium-config -n kube-system -o yaml | grep -i cidr

# Cilium doesn't use this. Cilium with its own IPAM mode
kubectl get nodes -o jsonpath='{.items[*].spec.podCIDR}'
```

# Add Envoy Gateway for Dapr communication between clusters
```
kubectl apply --server-side --force-conflicts -f https://github.com/envoyproxy/gateway/releases/download/v1.2.0/install.yaml

kubectl apply -f 1_eg.yaml

# For any issues
kubectl logs -n envoy-gateway-system -l control-plane=envoy-gateway --tail=20
```

# Patch Gateway Class
```
kubectl apply -f 2_eg-dapr.yaml

kubectl patch gatewayclass eg --type merge -p '{
  "spec": {
    "parametersRef": {
      "group": "gateway.envoyproxy.io",
      "kind": "EnvoyProxy",
      "name": "dapr-config",
      "namespace": "default"
    }
  }
}'
```

# Check Gateway Classes# Check Gateway API and create gateways
```
kubectl get gatewayclasses -o wide

kubectl apply -f 3_eg-gateway.yaml

kubectl apply -f cilium-gateway.yaml

kubectl get gateways

# Debug Cilium Gateway
# See the service Cilium created
kubectl get svc -l io.cilium.gateway/owning-gateway=default-gateway
# See the Cilium Envoy proxy pod
kubectl -n kube-system logs -l app.kubernetes.io/name=cilium-envoy -f

# Debug Envoy Gateway
# See the service EG created
kubectl get svc -l gateway.envoyproxy.io/owning-gateway-name=eg-gateway -n envoy-gateway-system
# See the EG Envoy proxy pod
kubectl -n envoy-gateway-system logs -l gateway.envoyproxy.io/owning-gateway-name=eg-gateway -f
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

helm install cilium cilium/cilium --version 1.19.1 --set kubeProxyReplacement=true --set gatewayAPI.enabled=true --namespace kube-system --set ipam.operator.clusterPoolIPv4PodCIDRList=10.2.0.0/16

# After CNI is installed, wait for pods to become Ready:
kubectl get pods --all-namespaces -w

# Check cilium status
cilium status --namespace kube-system

# Note: Make sure to configure the CNI plugin according to your cluster's pod CIDR
kubectl get configmap cilium-config -n kube-system -o yaml | grep -i cidr
```

# Add Envoy Gateway for Dapr communication between clusters
```
kubectl apply --server-side --force-conflicts -f https://github.com/envoyproxy/gateway/releases/download/v1.2.0/install.yaml

kubectl apply -f 1_eg.yaml

# For any issues
kubectl logs -n envoy-gateway-system -l control-plane=envoy-gateway --tail=20
```

# Patch Gateway Class
```
kubectl apply -f 2_eg-dapr.yaml

kubectl patch gatewayclass eg --type merge -p '{
  "spec": {
    "parametersRef": {
      "group": "gateway.envoyproxy.io",
      "kind": "EnvoyProxy",
      "name": "dapr-config",
      "namespace": "default"
    }
  }
}'
```

# Check Gateway API and create gateways
```
kubectl get gatewayclasses -o wide

kubectl apply -f 3_eg-gateway.yaml

kubectl apply -f cilium-gateway.yaml

kubectl get gateways

# Debug Cilium Gateway
# See the service Cilium created
kubectl get svc -l io.cilium.gateway/owning-gateway=default-gateway
# See the Cilium Envoy proxy pod
kubectl -n kube-system logs -l app.kubernetes.io/name=cilium-envoy -f

# Debug Envoy Gateway
# See the service EG created
kubectl get svc -l gateway.envoyproxy.io/owning-gateway-name=eg-gateway -n envoy-gateway-system
# See the EG Envoy proxy pod
kubectl -n envoy-gateway-system logs -l gateway.envoyproxy.io/owning-gateway-name=eg-gateway -f
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
- LoadBalancer is enabled for both clusters, so you can access services using `kubectl port-forward` or by exposing them via Docker port mappings.
