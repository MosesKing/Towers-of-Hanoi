# Towers-of-Hanoi

# Setup Guide for Kubernetes Cluster with Minikube

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