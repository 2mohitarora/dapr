for node in vcluster.cp.cluster-1 vcluster.node.cluster-1.worker-1; do
  docker exec "$node" mkdir -p /etc/containerd/certs.d/registry-1:5000
  docker exec "$node" sh -c 'cat > /etc/containerd/certs.d/registry-1:5000/hosts.toml << EOF
server = "http://registry-1:5000"

[host."http://registry-1:5000"]
  capabilities = ["pull", "resolve"]
  skip_verify = true
EOF'
done

for node in vcluster.cp.cluster-1 vcluster.node.cluster-1.worker-1; do
  docker exec "$node" systemctl restart containerd
done