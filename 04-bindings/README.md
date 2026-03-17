# Bindings

# Install postgres on k8s
```
# 1. Add the Bitnami repo
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

# 2. Install Postgres with a specific password and database name
helm install postgresql bitnami/postgresql \
  --set auth.database=orders \
  --set auth.postgresPassword=<STRONG_PASSWORD>

# 3. Create the logs table
export POSTGRES_PASSWORD=$(kubectl get secret --namespace default postgresql -o jsonpath="{.data.postgres-password}" | base64 -d)

kubectl run postgresql-client --rm --tty -i --restart='Never' \
  --namespace default \
  --image docker.io/bitnami/postgresql:latest \
  --env="PGPASSWORD=$POSTGRES_PASSWORD" \
  --command -- psql --host postgresql -U postgres -d orders -c "CREATE TABLE logs (id SERIAL PRIMARY KEY, message TEXT, created_at TIMESTAMP);"
```

# Build front end service and genid service container images
```
export KO_DOCKER_REPO=ko.local
export DOCKER_HOST="unix:///Users/mua0008/.colima/default/docker.sock"
ko build -B ./cronsvc
```

# Delete existing deployments
```
kubectl delete deployment cronsvc
```

# Deploy manifests
```
kubectl apply -f ./manifest
```

# Check Dapr UI 

# Check dapr components
```
dapr components -k
```

# Check k8s deployment 
```
kubectl get deployments -l app=cronsvc -o wide
```

# Check the pods in deployment 
```
kubectl get pods -l app=cronsvc
```

# Check application logs
```
kubectl logs -l app=cronsvc --prefix
```

# Check postgres data 