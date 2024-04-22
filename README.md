# Towers-of-Hanoi Kubernetes Integration Documentation

This document provides a comprehensive guide on setting up and deploying a Towers of Hanoi challenge within a Kubernetes environment, leveraging the Crossplane control plane to manage custom resources and compositions.

## Overview
The project involves deploying a Kubernetes-based solution for managing the "Tower of Hanoi" puzzle. It includes creating a custom Kubernetes controller, defining the necessary Custom Resource Definitions (CRDs), and configuring Crossplane compositions to facilitate the orchestration and operational management of puzzle instances.

## Step 1: Setting Up a Kubernetes Cluster
### Initial Setup
For development and testing, a local Kubernetes cluster is set up using Minikube. This involves installing virtualization tools and Minikube itself.

1. **Virtualization Environment Setup**:
   - Verify and set up a virtualization tool like VirtualBox.
  
2. **Minikube Installation**:
   - Follow the installation guide on the [Minikube GitHub page](https://github.com/kubernetes/minikube).

### Start and Verify the Minikube Cluster
- Start Minikube with Docker as the driver and verify that the cluster is operational by checking the status of the nodes.

```bash
minikube start --driver=docker
kubectl get nodes
```

## Step 2: Install and Set Up Crossplane
### Crossplane Installation
Install Crossplane using Helm by adding the stable repository, updating it, and then deploying Crossplane into the Kubernetes cluster.

```bash
helm repo add crossplane-stable https://charts.crossplane.io/stable
helm repo update
helm install crossplane crossplane-stable/crossplane --namespace crossplane-system --create-namespace
```

### Verify Crossplane Installation
Ensure that all Crossplane components are running correctly in the `crossplane-system` namespace.

```bash
kubectl get all -n crossplane-system
```

## Step 3: Implement Tower of Hanoi Logic with CRDs and Controller
### Custom Resource Definition (CRD)
- The `TowerChallenge` CRD is introduced to manage the lifecycle of the Tower of Hanoi challenges within the cluster.

### Apply the CRD
Deploy the CRD to the Kubernetes cluster:

```bash
kubectl apply -f towerchallenge.yaml
```

## Step 4: Controller Implementation
### TowerChallengeReconciler
This custom Kubernetes controller orchestrates the resolution of the Tower of Hanoi puzzle. It handles resource lifecycle management, including the creation, update, and cleanup of associated ConfigMaps.

### Key Features
- **Reconciliation Logic**: Manages the initialization, move calculation, and cleanup of resources.
- **Validation**: Ensures valid puzzle configurations.
- **ConfigMap Management**: Updates and creates ConfigMaps as needed.

## Step 5: Crossplane Composition Setup
### Composition Definition
Defines how the `CompositeResourceTowerChallenge` is composed of underlying resources, specifically focusing on managing puzzle instances and associated logging configurations.

### Deployment and Configuration
Deploy and configure the composition to link challenge instances with their logging mechanisms, ensuring comprehensive management and observability.

```bash
kubectl apply -f composition.yaml
```
