# Towers-of-Hanoi

## Step 1: Setting Up a Kubernetes Cluster

Before you begin, ensure you have a Kubernetes cluster. You can set this up on a cloud provider like AWS, GCP, or Azure, or locally using Minikube or kind. The instructions below describe setting up a local cluster using Minikube for development purposes.

### Install Minikube

1. **Prepare your virtualization environment**: Confirm that you have a virtualization tool such as VirtualBox.
2. **Install Minikube**: Follow the installation instructions on the [Minikube GitHub page](https://github.com/kubernetes/minikube).

### Start the Minikube Cluster

```bash
minikube start --driver=docker
```

### Verify the Cluster is Running
``` bash
kubectl get nodes
```

## Step 2: Install and Set Up Crossplane

### Install Crossplane

To install Crossplane using a Helm chart, run the following commands:

```bash
helm repo add crossplane-stable https://charts.crossplane.io/stable
helm repo update
```

Install Crossplane into Your Cluster
After adding the Helm repository, you can install Crossplane into your Kubernetes cluster:
``` bash
helm install crossplane crossplane-stable/crossplane --namespace crossplane-system --create-namespace
```

Verify Installation
To ensure Crossplane has been installed successfully, check the Crossplane components in the crossplane-system namespace:
``` bash
kubectl get all -n crossplane-system
```