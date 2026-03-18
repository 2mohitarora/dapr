# Things to do

- [x] Service Invocation
- [x] Service Invocation with retries and timeouts
- [ ] Service Invocation: Different name space
- [ ] Service Invocation: Different cluster
```
- Make Ingress Controller "Dapr-aware"
- Allow the sidecar to receive traffic from the Ingress controller's IP: dapr.io/sidecar-listen-addresses: "0.0.0.0, [::]"
- Ensure the dapr-app-id header is preserved and correctly handled by the gateway.
```
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
Remote Invoke Format: 
    InvokeMethod(ctx, "HTTPEndpointName", "v1.0/invoke/RemoteAppID/method/RemoteMethod", "VERB")
Comparison with Local Invoke Format: 
    InvokeMethod(ctx, "LocalAppID", "Method", "VERB")
```
```
Example: 
    var res, err := daprClient.InvokeMethod(ctx, "target-service-gateway", "v1.0/invoke/target-service/method/hello", "GET")
Comparison with Local Invoke Format: 
    var res, err := daprClient.InvokeMethod(ctx, "target-service", "hello", "GET")
```

- [ ] Service Invocation: External HTTP service
- [ ] Service Invocation: Different protocol (HTTP vs gRPC)
```
    dapr.io/app-port: "50051" # Your app's gRPC port
    dapr.io/app-protocol: "grpc" # <--- This ensures the final hop is gRPC
```
- [ ] Zero Trust Security
- [ ] Observability
    - [ ] Tracing
    - [ ] Metrics
    - [ ] Logs
    - [ ] Health
- [x] Bindings
    - [x] Input Bindings
    - [x] Output Bindings
- [ ] Actors
- [ ] Secrets Management
- [ ] Configuration Management
- [x] Workflow
- [x] Conversational AI
- [ ] AI Agents
- [ ] Jobs 
- [x] Pub/Sub
- [ ] Zero code Change Pub/Sub Move from Redis to Kafka
- [x] State Management
- [x] Outbox Pattern


    
