# Towers-of-Hanoi Kubernetes Integration Documentation

This document provides a comprehensive guide on setting up and deploying a Towers of Hanoi challenge within a Kubernetes environment, leveraging the Crossplane control plane to manage custom resources and compositions.

# Overview

This README outlines a solution for the Tower of Hanoi problem using Crossplane within a Kubernetes environment. The solution leverages Crossplane's powerful infrastructure-as-code capabilities to dynamically manage Kubernetes resources.
Challenge One: Basic Solution
Initial Approach

The initial challenge involves creating a Crossplane custom resource definition (CRD) that takes an integer discs as input and produces a series of Kubernetes ConfigMaps. Each ConfigMap represents a move necessary to solve the Tower of Hanoi puzzle for the specified number of discs across three rods (A, B, C).

## Key Components:

    CRD (TowerChallenge): Defines the structure for the input.
    Composition: Uses a custom function to interpret the TowerChallenge and generate the required ConfigMaps.
    XR Yaml - Our Claim
    Tower of Hanoi - Our Function and Package

### Example Resource:

```yaml
apiVersion: hanoi.com/v1alpha1
kind: TowerChallenge
metadata:
name: challenge3discs
spec:
discs: 3
```

## Expected Outputs:

ConfigMaps are created for each move, e.g., "Move disk 1 from rod A to rod C".
Challenge Two: Nested Solution
Advanced Approach

For the second challenge, a nested composition approach is proposed. This approach can involve defining additional layers of resource definitions that allow for more complex or multi-step operations that can depend on the state or outputs of previous operations.

## Assumptions:

    Nested compositions allow for defining complex dependencies and operations that can be reused in different contexts.
    Each step in solving the Tower of Hanoi can be treated as a sub-problem, potentially allowing for parallel or optimized processing.

    - the nested composition is just a proposal and not implementation, but an idea:
    Composition

### Assumptions and Interpretations

    The system assumes that the input CRD and the number of discs are correctly specified by the user.
    The moves required for any number of discs are calculated based on the standard recursive solution for the Tower of Hanoi problem.
    The Kubernetes cluster has sufficient permissions and capabilities to dynamically create and manage ConfigMaps.

## Composition Functions Over a Custom Operator?

Using Crossplane composition functions for this challenge leverages the existing Kubernetes API and Crossplane and other complex orchestration needs without developing a separate operator. However, this approach isn't the most ideal. Here's Why :

    - Complexity for Simple Tasks: For straightforward tasks like generating a sequence of moves, using Crossplane and writing custom functions might introduce unnecessary complexity compared to a custom operator tailored specifically for this task.
    - Performance Overhead: Running a composition function involves more overhead than executing an operator directly within the cluster, as it may involve additional layers of communication and processing.
    - sLimited Flexibility: While Crossplane compositions are powerful, they are somewhat limited to the structures and behaviors defined by Crossplane, whereas a custom operator can be highly tailored to specific needs and optimizations.

This setup in Crossplane uses its strengths in declarative configuration and maintaining state but might be over-engineering for problems that can be solved more efficiently with a custom Kubernetes operator.

## Additional Notes

    Consider security implications, especially in terms of resource permissions and access controls within Kubernetes.
    Consider why would choose using composition functions over an operator.
    Overall I learned a lot from this excercise and I am now motivated to keep learning Crossplane on my own,
    I am super interested in an opportunity to discuss how you would implement this to happen as I was not able to find a way to use functions to do this type of task.

## Custom Function Development

The function, which could be a separate project or container, would be triggered by changes to the TowerChallenge resource. It needs to:

    Read the number of discs from the TowerChallenge object.
    Implement the Tower of Hanoi algorithm to calculate the move sequence.
    Generate a series of ConfigMap definitions with move descriptions.
