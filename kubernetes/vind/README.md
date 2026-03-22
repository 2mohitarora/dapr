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
vcluster create cluster-1

kubectl get nodes
kubectl get namespaces
```

# Connect to your vcluster
```
vcluster connect cluster-1
```