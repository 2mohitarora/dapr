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

# Create your second cluster 

# Configure Registry for second cluster
```
# Start a local registry on the same Docker network as your vind cluster
docker run -d --name registry-2 --network vind-cluster-2 -p 5051:5000 registry:2

# Configure registry for cluster-2 so that nodes can pull from insecure registry
./cluster-2-script.sh
```

# Debug Cilium Gateway
```
# See the service Cilium created
kubectl get svc -l io.cilium.gateway/owning-gateway=default-gateway -n cilium
# See the Cilium Envoy proxy pod
kubectl -n kube-system logs -l app.kubernetes.io/name=cilium-envoy -f -n cilium
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
