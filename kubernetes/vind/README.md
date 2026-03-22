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
colima start --cpu 4 --memory 8 --disk 60 --runtime docker -p vind
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
sudo vcluster create cluster-1 --values cluster-1.yaml

kubectl get nodes
kubectl get namespaces
```
This will:

- Pull the vCluster container image (first run takes a minute)
- Start the control plane container
- Start the two worker node containers
- Wait for all nodes to become Ready
- Automatically switch your kubeconfig context to cluster-1

# Disconnect from first cluster
```
vcluster disconnect
```

# Create second vcluster
```
sudo vcluster create cluster-2 --values cluster-2.yaml
kubectl get nodes
kubectl get namespaces
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
- Each cluster uses a separate Docker network (vind-cluster1 and vind-cluster2) to keep them isolated from each other.
- LoadBalancer is enabled for both clusters, so you can access services using `kubectl port-forward` or by exposing them via Docker port mappings.
