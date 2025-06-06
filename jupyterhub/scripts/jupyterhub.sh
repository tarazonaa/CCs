#!/bin/bash

# To avoid using sudo and exporting the kubeconfig env var you can run in root shell
sudo KUBECONFIG=/etc/rancher/k3s/k3s.yaml helm repo add jupyterhub https://hub.jupyter.org/helm-chart/
sudo KUBECONFIG=/etc/rancher/k3s/k3s.yaml helm repo update

sudo KUBECONFIG=/etc/rancher/k3s/k3s.yaml helm upgrade --cleanup-on-fail \
  --install jupyter jupyterhub/jupyterhub \
  --namespace jupyter \
  --create-namespace \
  --values values.yaml
