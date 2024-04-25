# Tower of Hanoi Challenge Project

## Prerequisites

1. Kubernetes Cluster - A cluster where Crossplane and the custom operator will be deployed.
2. Install Crossplane - Install Crossplane on your Kubernetes cluster.
3. Kubebuilder - tool to scaffold out the operator.

## Cluster Creation & Setup

1. Create our Kubernetes Cluster

```bash
kind create cluster --name <desiredclustername>
```

2. Install Crossplane

```bash
kubectl create namespace crossplane-system
helm repo add crossplane-stable https://charts.crossplane.io/stable
helm repo update
helm install crossplane --namespace crossplane-system crossplane-stable/crossplane
```
- For the sake of this project my cluster was locally on my machine however, if we wanted to we could configure crossplane with various cloud providers, AWS, GCP, Azure. T


```go
// TowerChallengeSpec defines the desired state of TowerChallenge
type TowerChallengeSpec struct {
	// Discs is the number of discs in the Tower of Hanoi challenge
	// +kubebuilder:validation:Minimum=1
	Discs int `json:"discs"` // This is what we will add, it won't be there from running the command above. :p
}
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
