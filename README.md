# Install kubernetes cluster
```
Option 1: vind (Instructions in kubernetes/vind)
Option 2: k3s (Instructions in kubernetes/k3s)
```

# Install Dapr
```
brew install dapr/tap/dapr-cli
```

# Initialize Dapr on Kubernetes
```
kubectl config get-contexts
dapr init -k
```

# Verify Dapr
```
dapr status -k
kubectl get pods -o wide -n dapr-system
```

# Run dapr UI
```
dapr dashboard -k # For second cluster use -p 9999
```

# Install ko
```
brew install ko
```