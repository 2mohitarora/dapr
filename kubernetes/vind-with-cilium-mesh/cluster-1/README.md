# Instructions for first cluster
```
vcluster create cluster-1 --driver docker --values cluster-1.yaml
```
# Add the Jetstack repo
```
helm repo add cilium https://helm.cilium.io/
helm repo add jetstack https://charts.jetstack.io
helm repo update
```
# Upgrade Coredns
```
kubectl set image deployment/coredns \
  -n kube-system \
  coredns=registry.k8s.io/coredns/coredns:v1.14.2
```
# Install cilium 
```
helm install cilium cilium/cilium --version 1.19.1 --namespace cilium --create-namespace --set ipam.operator.clusterPoolIPv4PodCIDRList=10.1.0.0/16 --set routingMode=tunnel  --set tunnelProtocol=vxlan --set ipam.mode=cluster-pool
```
# After CNI is installed, wait for pods to become Ready:
```
kubectl get pods --all-namespaces -w
```
# Install cert-manager
```
helm install cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --create-namespace \
  --set crds.enabled=true
```
# After cert-manager, wait for pods to become Ready:
```
kubectl get pods --all-namespaces -w
```
# Install Gateway API
```
kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.2.0/experimental-install.yaml
```
# Install cert-manager issuer
```
kubectl apply -f cert-manager-issuer.yaml
kubectl get clusterissuer cilium-ca-issuer 
kubectl get certificates -n cert-manager 
```

# Upgrade Cilium with ClusterMesh
```
helm upgrade cilium cilium/cilium --version 1.19.1 \
  --namespace cilium \
  -f cilium-1-helm.yaml

# After Mesh is installed, wait for pods to become Ready:
kubectl get pods --all-namespaces -w

# Check cilium status
cilium status --namespace cilium
```
# Update CoreDNS ConfigMap. We have to do this because helm flag corednsAutoConfigure is not working as expected
```
kubectl edit configmap coredns -n kube-system

# Add the following to the configmap
rewrite name regex "(.*)\.clusterset\.local" "{1}.cluster.local"
rewrite name regex .*\.nodes\.vcluster\.com kubernetes.default.svc.cluster.local
kubernetes cluster.local clusterset.local in-addr.arpa ip6.arpa {
multicluster clusterset.local
```

# Check Gateway Class and Create Cilium Gateway
```
kubectl get gatewayclasses -o wide
kubectl apply -f ../cilium-gateway.yaml
kubectl get gateways -n cilium

# Create a route for front end service
kubectl apply -f ./frontendsvc-route.yaml
```

# Install Dapr Inject Dapr sidecar into Traefik deployment
```
dapr init -k

dapr status -k

# Make sure Dapr containers are ready
kubectl get pods -o wide -n dapr-system

# Run DAPR UI

dapr dashboard -k

```
# Start a local registry on the same Docker network as your vind cluster
```
docker run -d --name registry-1 --network vind-cluster-1 -p 5050:5000 registry:2
```
# Configure registry for cluster-1 so that nodes can pull from insecure registry
```
./cluster-1-script.sh
```
# Disconnect from first cluster
```
vcluster disconnect
```

