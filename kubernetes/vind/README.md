# Install kubectl
```
brew install kubectl
```
# Install helm
```
brew install helm
```
# Install Orbstack
```
brew install --cask orbstack
# Add local registries that will be created late before Starting Orbstack
# Add registry to docker daemon in ~/.docker/daemon.json
{
  "insecure-registries": ["localhost:5000", "localhost:5001"]
}
# Start Orbstack
```
# Install docker
```
brew install docker
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
helm install cilium cilium/cilium --version 1.19.1

# After CNI is installed, wait for pods to become Ready:
kubectl get pods --all-namespaces -w

# Check cilium status
cilium status --namespace default

# Note: Make sure to configure the CNI plugin according to your cluster's pod CIDR
kubectl get configmap  cilium-config -o yaml | grep -i cidr
kubectl get nodes -o jsonpath='{.items[*].spec.podCIDR}'
```
This will:

- Pull the vCluster container image (first run takes a minute)
- Start the control plane container
- Start one worker node container
- Install Cilium CNI plugin
- Wait for all nodes to become Ready
- Automatically switch your kubeconfig context to cluster-1

# Configure Registry for first cluster
```
# Start a local registry on the same Docker network as your vind cluster
docker run -d --name registry --network vind-cluster-1 -p 5000:5000 registry:2
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
helm install cilium cilium/cilium --version 1.19.1

# After CNI is installed, wait for pods to become Ready:
kubectl get pods --all-namespaces -w

# Check cilium status
cilium status --namespace default

# Note: Make sure to configure the CNI plugin according to your cluster's pod CIDR
kubectl get configmap  cilium-config -o yaml | grep -i cidr
kubectl get nodes -o jsonpath='{.items[*].spec.podCIDR}'
```

# Verify both clusters are running
```
vcluster disconnect
vcluster list
```

# Configure Registry for second cluster
```
# Start a local registry on the same Docker network as your vind cluster
docker run -d --name registry --network vind-cluster-2 -p 5001:5000 registry:2
```


# Check all the docker containers that are running
```
docker ps --format "table {{.Names}}\t{{.Status}}"
```

# Switching Between Clusters
```
vcluster connect cluster-1
vcluster disconnect
vcluster connect cluster-2
vcluster disconnect (Will bring to original kubecontext)
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
