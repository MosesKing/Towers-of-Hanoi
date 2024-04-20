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

output: 
```bash
kubectl get nodes
NAME       STATUS   ROLES           AGE   VERSION
minikube   Ready    control-plane   24m   v1.30.0
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
``` bash
C:\Windows\System32>helm install crossplane crossplane-stable/crossplane --namespace crossplane-system --create-namespace
NAME: crossplane
LAST DEPLOYED: Sat Apr 20 16:55:06 2024
NAMESPACE: crossplane-system
STATUS: deployed
REVISION: 1
TEST SUITE: None
NOTES:
Release: crossplane

Chart Name: crossplane
Chart Description: Crossplane is an open source Kubernetes add-on that enables platform teams to assemble infrastructure from multiple vendors, and expose higher level self-service APIs for application teams to consume.
Chart Version: 1.15.2
Chart Application Version: 1.15.2

Kube Version: v1.30.0 
```
Verify Installation
To ensure Crossplane has been installed successfully, check the Crossplane components in the crossplane-system namespace:
``` bash
kubectl get all -n crossplane-system
```

``` bash
C:\Windows\System32>kubectl get pods -n crossplane-system
NAME                                      READY   STATUS    RESTARTS   AGE
crossplane-6d84f5ccdf-pm9sk               1/1     Running   0          22s
crossplane-rbac-manager-fd57f7d55-ftbzc   1/1     Running   0          22s
```