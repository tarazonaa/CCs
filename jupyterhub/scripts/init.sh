#!/bin/bash

set -euo pipefail

# Check for root
if [[ $EUID -eq 0 ]]; then
  echo "Please run as a regular user with sudo, not as root."
  exit 1
fi

echo "Installing K3s..."
curl -sfL https://get.k3s.io | sh -

# Might want to add this to .profile
export KUBECONFIG=/etc/rancher/k3s/k3s.yaml

echo "Waiting for K3s to be ready..."
sudo kubectl wait --for=condition=ready node --all --timeout=180s

echo "Installing Helm..."
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

echo "Adding JupyterHub Helm repo..."
helm repo add jupyterhub https://jupyterhub.github.io/helm-chart/
helm repo update

PRIVATE_IP=$(ip addr show | awk '/inet / && $NF ~ /^e/{print $2}' | cut -d/ -f1 | head -n 1)
if [[ -z "$PRIVATE_IP" ]]; then
  echo "Failed to determine private IP."
  exit 1
fi

# Create a config.yaml for JupyterHub
cat <<EOF > jupyterhub-config.yaml
proxy:
  service:
    type: NodePort
    nodePorts:
      http: 30080
      https: 30443
  secretToken: "$(openssl rand -hex 32)"
EOF

echo "Installing JupyterHub with Helm..."
helm upgrade --install jhub jupyterhub/jupyterhub \
  --namespace jhub --create-namespace \
  --version=3.3.7 \
  -f jupyterhub-config.yaml

echo "JupyterHub installed."
echo "Access it via: http://$PRIVATE_IP:30080"
