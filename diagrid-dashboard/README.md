# Install Diagrid Dashboard
```
 kubectl apply -f ./manifest
```
# Access the dashboard
```
kubectl port-forward svc/diagrid-dashboard 8082:8080
```
# Open the dashboard
```
http://localhost:8082
```

