# Both clusters should see each other
cilium clustermesh status --context vcluster-docker_cluster-1 --namespace cilium
cilium clustermesh status --context vcluster-docker_cluster-2 --namespace cilium

# Each cluster should see nodes from the other
kubectl --context vcluster-docker_cluster-1 exec -n cilium ds/cilium -- cilium node list
kubectl --context vcluster-docker_cluster-2 exec -n cilium ds/cilium -- cilium node list

# Identities from both clusters should be visible
kubectl --context vcluster-docker_cluster-1 exec -n cilium ds/cilium -- cilium identity list | head -20
kubectl --context vcluster-docker_cluster-2 exec -n cilium ds/cilium -- cilium identity list | head -20

# Full connectivity test across the mesh
cilium connectivity test --context vcluster-docker_cluster-1 --multi-cluster vcluster-docker_cluster-2 --test pod-to-pod,pod-to-service --namespace cilium

# Clean up test resources
kubectl --context vcluster-docker_cluster-1 delete ns cilium-test-1
kubectl --context vcluster-docker_cluster-2 delete ns cilium-test-1

# Create dns-validator pod on both clusters
kubectl --context vcluster-docker_cluster-1 apply -f mcs-dns-check.yaml
kubectl --context vcluster-docker_cluster-2 apply -f mcs-dns-check.yaml

# MCS-API Validator, Scearion 1 : Service only on cluster-2
```
# Create mcs-test namespace of cluster-2
kubectl --context vcluster-docker_cluster-2 create namespace mcs-test

# Create deployment, service and serviceexport on cluster-2
kubectl --context vcluster-docker_cluster-2 apply -f mcs-test.yaml

# Check serviceexport and serviceimport objects on cluster-2
kubectl --context vcluster-docker_cluster-2 get serviceexport -n mcs-test
kubectl --context vcluster-docker_cluster-2 get serviceimport -n mcs-test

# Try to resolve the remote service
kubectl --context vcluster-docker_cluster-1 exec dns-validator -- nslookup web.mcs-test.svc.clusterset.local

# Create namespace on cluster-1
kubectl --context vcluster-docker_cluster-1 create namespace mcs-test

# Check serviceimport objects appearing on cluster-1
kubectl --context vcluster-docker_cluster-1 get serviceimport -n mcs-test

# Try to resolve the remote service again
kubectl --context vcluster-docker_cluster-1 exec dns-validator -- nslookup web.mcs-test.svc.clusterset.local

kubectl --context vcluster-docker_cluster-1 run curl-test --rm -it --image=curlimages/curl --restart=Never -n mcs-test --labels="app=curl-test" -- curl -v -s --max-time 5 http://web.mcs-test.svc.clusterset.local -o /dev/null -w "Response from: %{remote_ip}\n"

What this proves: Cilium’s MCS controller has successfully synced the ServiceImport and CoreDNS is correctly configured with the clusterset stub-domain.
```

# MCS-API Validator, Scearion 2 : Headless service on cluster-1 only and we will resolve it from cluster-2
```
# Create deployment, headless service and serviceexport on cluster-1
kubectl --context vcluster-docker_cluster-1 apply -f mcs-headless-test.yaml

# Check serviceexport and serviceimport objects on cluster-1
kubectl --context vcluster-docker_cluster-1 get serviceexport -n mcs-test
kubectl --context vcluster-docker_cluster-1 get serviceimport -n mcs-test

# Check serviceimport objects appearing on cluster-2 (Had to restart operator cluster-2)
kubectl --context vcluster-docker_cluster-2 get serviceimport -n mcs-test

# Try to resolve the remote service
kubectl --context vcluster-docker_cluster-2 exec dns-validator -- nslookup web-headless.mcs-test.svc.clusterset.local

kubectl --context vcluster-docker_cluster-2 run curl-test --rm -it --image=curlimages/curl --restart=Never -n mcs-test --labels="app=curl-test" -- curl -v -s --max-time 5 http://web-headless.mcs-test.svc.clusterset.local -o /dev/null -w "Response from: %{remote_ip}\n"
```

# MCS-API Validator, Scearion 3 : CiliumNetworkPolicy and L7 Policy and CiliumClusterwideNetworkPolicy 
 CiliumNetworkPolicy (CNP) and CiliumClusterwideNetworkPolicy (CCNP) are not automatically replicated across clusters in a Cilium Cluster Mesh. Cluster Mesh synchronizes identities, pods, and services to allow cross-cluster communication, security policies must be managed separately in each cluster

# Deny all remote cluster traffic (CiliumClusterwideNetworkPolicy)
```
kubectl --context vcluster-docker_cluster-1 apply -f policies/deny-all-cluster-1.yaml
kubectl --context vcluster-docker_cluster-2 apply -f policies/deny-all-cluster-2.yaml

kubectl --context vcluster-docker_cluster-1 get ccnp
kubectl --context vcluster-docker_cluster-2 get ccnp

# Try connecting from cluster-1 to cluster-2 and cluster-2 to cluster-1, it will fail
```

# Allow cluster-1 to access web (CiliumNetworkPolicy)
```
kubectl --context vcluster-docker_cluster-2 apply -f policies/allow-web.yaml
kubectl --context vcluster-docker_cluster-2 get cnp -n mcs-test
```

# Allow cluster-2 to access web-headless (CiliumNetworkPolicy)
```
kubectl --context vcluster-docker_cluster-2 apply -f policies/allow-web-headless.yaml
kubectl --context vcluster-docker_cluster-2 get cnp -n mcs-test
```

# Only allow prod to prod (CiliumClusterwideNetworkPolicy)
```
kubectl --context vcluster-docker_cluster-1 apply -f policies/only-allow-prod-to-prod.yaml
kubectl --context vcluster-docker_cluster-2 apply -f policies/only-allow-prod-to-prod.yaml
```



--------NOT EXPLORED YET--------
# The Hubble "Flow-Watch" Command
Finally, use Hubble to see the "Identity" magic in action. This command allows you to see traffic filtered by the numeric Identity we discussed earlier.

# Port-forward Hubble if not already open
cilium hubble port-forward &

# Watch flows between clusters, showing the Identity IDs
hubble observe --follow --output json | jq '{
  time: .time,
  source: .source.identity,
  dest: .destination.identity,
  verdict: .verdict
}'
What this proves: You are seeing the actual uint32 identities assigned by the KVStore, confirming that your 100-cluster identity sync is healthy.

----------------