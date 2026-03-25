This "Validation Suite" is designed to be your final sanity check before moving a cluster into the production fleet. It tests the three pillars of your architecture: Connectivity, Discovery, and L7 Security.

1. The Multi-Cluster Connectivity Check
Before anything else, use the built-in CLI tool. This is more robust than a simple ping because it tests the eBPF datapath directly between Pods in Cluster 1 and Cluster 2.

Bash
# Run from your management terminal
cilium connectivity test \
  --context cluster-east \
  --multi-cluster cluster-west \
  --test pod-to-pod,pod-to-service
What this proves: VXLAN tunnels are up, Node-to-Node reachability on port 8472 is open, and IPAM isn't overlapping.

2. The MCS-API DNS Validator
Since you are using the ServiceExport model, you need to verify that clusterset.local is resolving correctly. Use this "Tracer" pod to check the DNS path.

YAML
# mcs-dns-check.yaml
apiVersion: v1
kind: Pod
metadata:
  name: dns-validator
  namespace: default
spec:
  containers:
  - name: alpine
    image: alpine
    command: ["sleep", "infinite"]
The Test Commands:

Bash
# 1. Apply the pod
kubectl apply -f mcs-dns-check.yaml

# 2. Try to resolve the remote service (replace <svc> and <ns>)
kubectl exec dns-validator -- nslookup <service-name>.<namespace>.svc.clusterset.local
What this proves: Cilium’s MCS controller has successfully synced the ServiceImport and CoreDNS is correctly configured with the clusterset stub-domain.

3. The L7 Policy "Kill-Switch" Test
For your Kafka and API Gateway goals, you need to ensure that the Envoy proxy is actually intercepting and enforcing Layer 7 rules. We'll use a simple HTTP test to simulate this.

Apply this Policy:

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

Summary Checklist for your Runbook
[ ] cilium status shows "ClusterMesh: OK"

[ ] cilium identity list | wc -l is within expected limits (< 0.5 pod ratio)

[ ] nslookup ...clusterset.local returns a valid IP

[ ] L7 curl to a restricted path returns a 403