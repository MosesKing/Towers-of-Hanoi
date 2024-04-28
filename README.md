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
3. Install Kubebuilder

```bash
curl -L -o kubebuilder "https://go.kubebuilder.io/dl/latest/$(go env GOOS)/$(go env GOARCH)"\nchmod +x kubebuilder && sudo mv kubebuilder /usr/local/bin/
```


### Example Resource:

```yaml
apiVersion: hanoi.com/v1alpha1
kind: TowerChallenge
metadata:
name: challenge3discs
spec:
discs: 3
```
## Custom Function Development

The function, which could be a separate project or container, would be triggered by changes to the TowerChallenge resource. It needs to:

- Read the number of discs from the TowerChallenge object.
- Implement the Tower of Hanoi algorithm to calculate the move sequence.
- Generate a series of ConfigMap definitions with move descriptions.

## Composition Functions Over a Custom Operator?

Using Crossplane composition functions for this challenge leverages the existing Kubernetes API and Crossplane and other complex orchestration needs without developing a separate operator. However, this approach isn't the most ideal. Here's Why:

- Complexity for Simple Tasks: For straightforward tasks like generating a sequence of moves, using Crossplane and writing custom functions might introduce unnecessary complexity compared to a custom operator tailored specifically for this task.
- Performance Overhead: Running a composition function involves more overhead than executing an operator directly within the cluster, as it may involve additional layers of communication and processing.
- Limited Flexibility: While Crossplane compositions functions are powerful, as they are able to add a level of dynamicness to static nature of these files. they are still somewhat limited to the structures and behaviors defined by Crossplane, whereas a custom operator can be highly tailored to specific needs and optimizations.

This setup in Crossplane uses its strengths in declarative configuration and maintaining state but might be over-engineering for problems that can be solved more efficiently with a custom Kubernetes operator. I was still able to use Crossplane extensions, but not for the aspect of calculating of Hanoi Tower puzzle.

## Expected Outputs:

ConfigMaps are created for each move, e.g., "Move disk 1 from rod A to rod C".
Challenge Two: Nested Solution
Advanced Approach

# The Second Challenge: 
- a nested composition approach is proposed. This approach can involve defining additional layers of resource definitions that allow for more complex or multi-step operations that can depend on the state or outputs of previous operations.

## Assumptions:
- Nested compositions allow for defining complex dependencies and operations that can be reused in different contexts.
- Each step in solving the Tower of Hanoi can be treated as a sub-problem, potentially allowing for parallel or optimized processing.

### Assumptions and Interpretations
- The system assumes that the input CRD and the number of discs are correctly specified by the user.
- The moves required for any number of discs are calculated based on the standard recursive solution for the Tower of Hanoi problem.
- The Kubernetes cluster has sufficient permissions and capabilities to dynamically create and manage ConfigMaps.


## Additional Notes

- Consider security implications, especially in terms of resource permissions and access controls within Kubernetes.
- Consider why would choose using composition functions over an operator.
- Overall I learned a lot from this excercise and I am now motivated to keep learning Crossplane on my own, I am super interested in an opportunity to discuss how you would implement this to happen as I was not able to find a way to use functions to do this type of task.
