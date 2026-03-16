# Things to do

- [ ] Service Invocation: Different name space
- [ ] Service Invocation: Different cluster
```
apiVersion: dapr.io/v1alpha1
kind: HTTPEndpoint
metadata:
  name: target-service-gateway
spec:
  baseUrl: https://target-service-gateway.hosting-cluster.com
  headers:
  - name: "dapr-api-token" # If your gateway requires a secret
    value: "your-gateway-token"
```
```
Format: InvokeMethod(ctx, "HTTPEndpointName", "v1.0/invoke/RemoteAppID/method/RemoteMethod", "VERB")
```
- [ ] Service Invocation: External HTTP service
- [ ] Service Invocation: Different protocol (HTTP vs gRPC)
- [ ] Zero Trust Security
- [ ] Observability
    - [ ] Tracing
    - [ ] Metrics
    - [ ] Logs
    - [ ] Health
- [ ] Bindings
    - [ ] Input Bindings
    - [ ] Output Bindings
- [ ] Actors
- [ ] Secrets Management
- [ ] Configuration Management
- [ ] Workflow
- [ ] Conversational AI
- [ ] AI Agents
- [ ] Jobs 
- [ ] Pub/Sub with CloudEvents


    
