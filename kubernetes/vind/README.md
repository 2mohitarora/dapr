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

kubectl config get-contexts
kubectl get nodes
kubectl get namespaces

helm repo add cilium https://helm.cilium.io/

kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.2.0/experimental-install.yaml

helm install cilium cilium/cilium --version 1.19.1 --set kubeProxyReplacement=true --set gatewayAPI.enabled=true --namespace kube-system --set ipam.operator.clusterPoolIPv4PodCIDRList=10.1.0.0/16

# After CNI is installed, wait for pods to become Ready:
kubectl get pods --all-namespaces -w

# Check cilium status
cilium status --namespace kube-system

# Note: Make sure to configure the CNI plugin according to your cluster's pod CIDR
kubectl get configmap cilium-config -n kube-system -o yaml | grep -i cidr
kubectl get nodes -o jsonpath='{.items[*].spec.podCIDR}'
```

# Check Gateway API
```
kubectl get gatewayclasses
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
kubectl config get-contexts
kubectl get nodes
kubectl get namespaces

helm repo add cilium https://helm.cilium.io/

kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.2.0/experimental-install.yaml

helm install cilium cilium/cilium --version 1.19.1 --set kubeProxyReplacement=true --set gatewayAPI.enabled=true --namespace kube-system --set ipam.operator.clusterPoolIPv4PodCIDRList=10.2.0.0/16

# After CNI is installed, wait for pods to become Ready:
kubectl get pods --all-namespaces -w

# Check cilium status
cilium status --namespace kube-system

# Note: Make sure to configure the CNI plugin according to your cluster's pod CIDR
kubectl get configmap cilium-config -n kube-system -o yaml | grep -i cidr
kubectl get nodes -o jsonpath='{.items[*].spec.podCIDR}'
```

# Check Gateway API
```
kubectl get gatewayclasses
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
