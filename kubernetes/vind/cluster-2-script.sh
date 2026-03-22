for node in vcluster.cp.cluster-2 vcluster.node.cluster-2.worker-1; do
  docker exec "$node" mkdir -p /etc/containerd/certs.d/registry-2:5000
  docker exec "$node" sh -c 'cat > /etc/containerd/certs.d/registry-2:5000/hosts.toml << EOF
server = "http://registry-2:5000"

[host."http://registry-2:5000"]
  capabilities = ["pull", "resolve"]
  skip_verify = true
EOF'
done

for node in vcluster.cp.cluster-2 vcluster.node.cluster-2.worker-1; do
  docker exec "$node" systemctl restart containerd
done