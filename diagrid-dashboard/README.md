# Create configmap for dashboard from redis-store.yaml in workflow directory
```
kubectl create configmap diagrid-dashboard-config --from-file=../05-workflow/manifest/redis-store.yaml
# Edit cofigmap with real password
kubectl get secret --namespace default redis -o jsonpath="{.data.redis-password}" | base64 --decode
kubectl edit configmap diagrid-dashboard-config
```
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

