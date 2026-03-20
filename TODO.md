# Things to do

- [x] Service Invocation
- [x] Service Invocation with retries and timeouts
- [x] Service Invocation: Different name space
- [x] Service Invocation: Different cluster
```
Example: 
    var res, err := daprClient.InvokeMethod(ctx, "<TARGET_SERVICE_ENDPOINT>", "v1.0/invoke/<AppID>/method/<Method>", "VERB")
Comparison with Local Invoke Format: 
    var res, err := daprClient.InvokeMethod(ctx, "<AppID>", "<Method>", "VERB")
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


    
