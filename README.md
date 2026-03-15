# Install kubectl
brew install kubectl

# Install helm
brew install helm

# Install colima
brew install colima

# Change colima template
```colima template
Change few settings```

# Start colima with k3s
colima start --kubernetes

# Verify kubernetes
kubectl get svc -o wide -A 
colima ssh 

# Docker checks
docker ps 
docker context ls

# Install Dapr
brew install dapr/tap/dapr-cli

# Initialize Dapr on Kubernetes
dapr init -k

# Verify Dapr
dapr status -k
kubectl get pods -o wide -n dapr-system

# Run dapr UI
dapr dashboard -k

# Install ko
brew install ko
