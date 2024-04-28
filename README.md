# Tower of Hanoi Challenge Project

## Prerequisites

1. Kubernetes Cluster - A cluster where Crossplane and the custom operator will be deployed.
2. Install Crossplane - Install Crossplane on your Kubernetes cluster.
3. Go Function Docker - moe12358/hanoi-function:latest 

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

## Function Creation: 
```bash
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv" // Import strconv package for string conversion utilities.
)

// Move represents a single operation in the Tower of Hanoi puzzle.
type Move struct {
	Disc int    `json:"disc"` // The disc number being moved.
	From string `json:"from"` // The rod from which the disc is moved.
	To   string `json:"to"`   // The rod to which the disc is moved.
}

// solveHanoi recursively calculates the moves required to solve the Tower of Hanoi puzzle.
func solveHanoi(n int, from, to, aux string, moves *[]Move) {
	if n == 1 {
		// Base case: only one disc to move directly from the source to destination.
		*moves = append(*moves, Move{Disc: n, From: from, To: to})
		return
	}
	// Recursive case: Move n-1 discs to the auxiliary rod.
	solveHanoi(n-1, from, aux, to, moves)
	// Move the nth disc to the destination rod.
	*moves = append(*moves, Move{Disc: n, From: from, To: to})
	// Move the n-1 discs from the auxiliary rod to the destination rod.
	solveHanoi(n-1, aux, to, from, moves)
}

// handler handles the HTTP requests for solving the Tower of Hanoi puzzle.
func handler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the number of discs from the query parameter "discs"
	discsStr := r.URL.Query().Get("discs")
	if discsStr == "" {
		// Respond with an error if the "discs" parameter is missing.
		http.Error(w, "Missing discs parameter", http.StatusBadRequest)
		return
	}

	// Convert the discs string to an integer
	discs, err := strconv.Atoi(discsStr)
	if err != nil {
		// Respond with an error if the conversion fails (non-integer value).
		http.Error(w, "Invalid number of discs: must be an integer", http.StatusBadRequest)
		return
	}

	// Ensure the number of discs is a positive integer
	if discs <= 0 {
		// Respond with an error if the number of discs is not positive.
		http.Error(w, "Invalid number of discs: must be a positive integer", http.StatusBadRequest)
		return
	}

	// Compute the moves for the specified number of discs
	var moves []Move
	solveHanoi(discs, "A", "C", "B", &moves)

	// Set the response header to application/json and encode the moves into JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(moves)
}

// main initializes the HTTP server.
func main() {
	// Set up HTTP routing and print server start-up message.
	http.HandleFunc("/", handler)
	fmt.Println("Starting server at port 8080...")
	// Start an HTTP server listening on port 8080 and log any errors.
	log.Fatal(http.ListenAndServe(":8080", nil))
}

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

## 1. Define TowerChallenge Custom Resource (CR)
- Create a CRD for TowerChallenge to accept the number of discs as input.
```yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: towerchallenges.hanoi.com
spec:
  group: hanoi.com
  names:
    kind: TowerChallenge
    listKind: TowerChallengeList
    plural: towerchallenges
    singular: towerchallenge
  scope: Namespaced
  versions:
  - name: v1alpha1
    served: true
    storage: true
    schema:
      openAPIV3Schema:
        type: object
        properties:
          spec:
            type: object
            properties:
              discs:
                type: integer
                description: "Number of discs in the Tower of Hanoi challenge."
          status:
            type: object
            properties:
              moves:
                type: array
                items:
                  type: string

```
## 2. Install Provider Kubernetes
- Ensure the provider-kubernetes is installed to manage arbitrary Kubernetes resources. This kubernetes provider comes from this place, I found this after engaging with the Crossplane Slack Community, it will help me generate `ConfigMap`.
```bash
kubectl crossplane install provider crossplane/provider-kubernetes:<version>
```

## 3. Setup Provider Configuration: 
- Create a `ProviderConfig` using credentials to interact w/ the Kubernetes API
```yaml
apiVersion: kubernetes.crossplane.io/v1alpha1
kind: ProviderConfig
metadata:
  name: kubernetes-provider
spec:
  credentials:
    source: Secret
    secretRef:
      namespace: crossplane-system
      name: kubernetes-creds
      key: kubeconfig

```
## 4. Create Composition:
- Define `Composition` that specifies how the `TowerChallenge` is composed of `ConfigMap` resouces 
## Custom Function Development
```yaml
apiVersion: apiextensions.crossplane.io/v1
kind: Composition
metadata:
  name: tower-hanoi-composition
spec:
  compositeTypeRef:
    apiVersion: hanoi.com/v1alpha1
    kind: TowerChallenge
  resources:
  - name: moves-configmap
    base:
      apiVersion: v1
      kind: ConfigMap
      metadata:
        generateName: hanoi-moves-
      data:
        moves: ""
    patches:
    - fromFieldPath: "spec.discs"
      toFieldPath: "data.discs"

```
## 5. Deploy and Manage Resources
Create instances of TowerChallenge to see Crossplane orchestrate the creation of ConfigMap resources based on the moves calculated by the solution to the Tower of Hanoi problem.
```bash
kubectl crossplane install provider crossplane/provider-kubernetes:<version>
```

The function, which could be a separate project or container, would be triggered by changes to the TowerChallenge resource. It needs to:

- Read the number of discs from the TowerChallenge object.
- Implement the Tower of Hanoi algorithm to calculate the move sequence.
- Generate a series of ConfigMap definitions with move descriptions.

## Composition Functions Over a Custom Operator?

Using Crossplane composition functions for this challenge leverages the existing Kubernetes API and Crossplane and other complex orchestration needs without developing a separate operator. However, this approach isn't the most ideal. Here's Why:

- **Scalability**: Custom functions in Crossplane are designed to handle transformations or configurations at a declarative level within a specific framework, often lacking the dynamic scaling capabilities. For example, a custom function would not automatically adjust computing resources or scale across multiple clusters based on real-time demands or changes in the environment.

- **Automation and Lifecycle Management**: Custom functions generally do not manage state or handle complex operational tasks like software updates, backup, or recovery. If the Tower of Hanoi application needed regular updates or required persistent state management across restarts, a custom function would not autonomously manage these processes, potentially leading to manual interventions for updates and maintenance.

- **Flexibility**: While custom functions can be adjusted or extended, they primarily focus on specific tasks and lack the broader operational control that operators offer. For instance, integrating new monitoring tools or changing deployment strategies based on evolving requirements would require manual reconfiguration or rewriting of the function, as opposed to an operator which could adapt its behavior based on predefined policies or external triggers.

This setup in Crossplane uses its strengths in declarative configuration and maintaining state but might be over-engineering for problems that can be solved more efficiently with a custom Kubernetes operator. I was still able to use Crossplane extensions, but not for the aspect of calculating of Hanoi Tower puzzle.

## Expected Outputs:

ConfigMaps are created for each move, e.g., "Move disk 1 from rod A to rod C".
Challenge Two: Nested Solution
Advanced Approach

# The Second Challenge: 
-  To do the second crossplane we will define a hierarchy of resources where each level can influence the configuration and behavior of subsequent resources.

## Proposals:
Steps for Nested Composition for Tower of Hanoi
- **Define Resource Definitions for Each Move**:
Each disk move can be represented as a unique ConfigMap, where each ConfigMap could contain the details of the move such as disk number, source rod, and destination rod.
- **Create Intermediate Compositions**:
Define intermediate compositions that represent each phase of the puzzle. For example, moving a disk from rod A to rod C could trigger a condition in the next composition that deals with moving a disk from rod A to rod B, based on the previous state.
- **Link Compositions Using Outputs and References**:
Utilize outputs from one composition as inputs to the next. This can be done by referencing the state or outputs of the previous ConfigMap in subsequent compositions, ensuring that moves follow the logical rules of the Tower of Hanoi.
- **Implement Dependency and Order Management**:
Manage dependencies to ensure that no move violates the puzzle rules. This may involve creating explicit dependencies in the resource definitions to ensure that the moves are executed in a legal and orderly manner.
- **Automate Puzzle Completion Checks:**
- Each composition perhaps could check the puzzle is completed (i.e., all disks are moved to the final rod in the correct order). If not, it will trigger the next set of moves.

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

## Helpful Resources:
https://youtu.be/otwg-bO757A?si=K_tXxjCJygTpxQgy
https://youtu.be/jjtpEhvwgMw?si=jyzW4gr9t8QOCaZe
https://youtu.be/XSzKs97Ls4g?si=JhejywibCcVGC43K
https://youtu.be/m6-xIhQcCe4?si=ZINf2w04MitJm_i0
