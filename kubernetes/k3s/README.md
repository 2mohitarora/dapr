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

# Change colima template
```
colima template
Change following settings
- network.address: true
- vmType: vz
- rosetta: true
- mountType: virtiofs
```

# Start colima with k3s
```
colima start --kubernetes
```

# Verify kubernetes
```
kubectl get svc -o wide -A 
```

# Checks
```
docker ps 
docker context ls
kubectl config get-contexts
```

# Instructions for second cluster (It's not required initially)

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

# Check Contexts
```
kubectl config get-contexts
docker context ls
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

# Run dapr UI
```
dapr dashboard -k -p 9999
```
