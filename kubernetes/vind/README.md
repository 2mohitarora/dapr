# Install kubectl
```
brew install kubectl
```
# Install helm
```
brew install helm
```
# Install colima
```
brew install colima
```
# Install docker
```
brew install docker
```
# Install Cilium CLI
```
brew install cilium-cli
```
# Change colima template
```
colima template
Change following settings
- network.address: true
- vmType: vz
- rosetta: true
- mountType: virtiofs
```
# Configure vind template
```
cd ~/.colima/_templates
cp default.yaml vind.yaml 
```
# Start colima with vind template and enough resources for two 3-node clusters
```
colima start --cpu 4 --memory 8 --disk 60 -p vind
```
# Configure docker
```
docker context list
docker context use colima-vind
```
# Install vcluster CLI
```
brew install vcluster
vcluster version
```
# Configure vcluster - Set Docker as default deamon
```
vcluster use driver docker
```

# Start vCluster Platform UI
```
vcluster platform start
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
cilium status --namespace cilium

# Note: Make sure to configure the CNI plugin according to your cluster's pod CIDR
# Check the pod CIDR
kubectl cluster-info dump | grep -m 1 cluster-cidr
```
This will:

- Pull the vCluster container image (first run takes a minute)
- Start the control plane container
- Start one worker node container
- Install Cilium CNI plugin
- Wait for all nodes to become Ready
- Automatically switch your kubeconfig context to cluster-1

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
cilium status --namespace cilium

# Note: Make sure to configure the CNI plugin according to your cluster's pod CIDR
# Check the pod CIDR
kubectl cluster-info dump | grep -m 1 cluster-cidr
```

# Verify both clusters are running
```
vcluster disconnect
vcluster list
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
colima delete -p vind
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
