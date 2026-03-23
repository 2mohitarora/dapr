
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