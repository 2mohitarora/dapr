# Things to do

- [x] Service Invocation
- [x] Service Invocation with retries and timeouts
- [x] Service Invocation: Different name space
- [x] Service Invocation: Different cluster
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
    InvokeMethod(ctx, "<HTTPEndpoint>", "v1.0/invoke/<AppID>/method/<Method>", "VERB")
Comparison with Local Invoke Format: 
    InvokeMethod(ctx, "<AppID>", "<Method>", "VERB")
```
```
Example: 
    var res, err := daprClient.InvokeMethod(ctx, "target-service-gateway", "v1.0/invoke/target-service/method/hello", "GET")
Comparison with Local Invoke Format: 
    var res, err := daprClient.InvokeMethod(ctx, "target-service", "hello", "GET")
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


    
