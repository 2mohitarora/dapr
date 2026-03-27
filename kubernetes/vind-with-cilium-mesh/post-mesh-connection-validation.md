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

What this proves: Cilium’s MCS controller has successfully synced the ServiceImport and CoreDNS is correctly configured with the clusterset stub-domain.
```



--------NOT EXPLORED YET--------

YAML
# l7-test-policy.yaml
apiVersion: "cilium.io/v2"
kind: CiliumNetworkPolicy
metadata:
  name: limit-to-health-check
spec:
  endpointSelector:
    matchLabels:
      app: my-service
  ingress:
  - fromEndpoints:
    - matchLabels:
        app: guest-pod
    toPorts:
    - ports:
      - port: "80"
        protocol: TCP
      rules:
        http:
        - method: "GET"
          path: "/health"
Run the Validation:

Bash
# This SHOULD work (200 OK)
kubectl exec guest-pod -- curl -I http://my-service/health

# This SHOULD fail (403 Forbidden - enforced by Envoy)
kubectl exec guest-pod -- curl -I http://my-service/admin
What this proves: The Cilium-Envoy daemon is operational and correctly attached to your pods' network namespace.

4. The Hubble "Flow-Watch" Command
Finally, use Hubble to see the "Identity" magic in action. This command allows you to see traffic filtered by the numeric Identity we discussed earlier.

Bash
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
