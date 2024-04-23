# Towers-of-Hanoi Kubernetes Integration Documentation

This document provides a comprehensive guide on setting up and deploying a Towers of Hanoi challenge within a Kubernetes environment, leveraging the Crossplane control plane to manage custom resources and compositions.

## Overview
The project involves deploying a Kubernetes-based solution for managing the "Tower of Hanoi" puzzle. It includes creating a custom Kubernetes controller, defining the necessary Custom Resource Definitions (CRDs), and configuring Crossplane compositions to facilitate the orchestration and operational management of puzzle instances.

Tower of Hanoi Kubernetes Custom Resource and Composition Documentation

Overview:
This documentation provides a comprehensive guide to utilizing Kubernetes Custom Resources (CRs) and Compositions to tackle the Tower of Hanoi problem within a Kubernetes environment. The Tower of Hanoi problem is a classic mathematical puzzle that involves moving a stack of disks from one rod to another, adhering to specific rules.

Components:

    TowerChallenge Custom Resource Definition (CRD):
        Resource Name: TowerChallenge
        Description: This CRD defines a custom resource named TowerChallenge with a discs field indicating the number of discs in the Tower of Hanoi puzzle.
        Usage: Users can create instances of TowerChallenge by specifying the number of discs, allowing Kubernetes to manage Tower of Hanoi challenges as custom resources.

    Composition for Tower Challenge:
        Composition Name: tower-hanoi-composition
        Description: This composition orchestrates the solving of Tower of Hanoi challenges by defining a pipeline that executes a function to calculate moves.
        Usage: When a TowerChallenge instance is created, this composition is triggered to calculate the optimal moves required to solve the Tower of Hanoi puzzle with the specified number of discs.

    Composition for Tower Hanoi Solution:
        Composition Name: tower-hanoi-composition (with label purpose: tower-hanoi-solution)
        Description: This composition is responsible for generating the solution to the Tower of Hanoi puzzle by dynamically setting the moves in a ConfigMap.
        Usage: Upon completion of move calculation, this composition updates a ConfigMap named tower-hanoi-moves with the optimal move sequence required to solve the Tower of Hanoi puzzle.

    Composite Resource Definition for XTowerChallenge:
        Resource Name: XTowerChallenge
        Description: This CRD defines a composite resource that corresponds to the TowerChallenge custom resource.
        Usage: Kubernetes operators interact with instances of XTowerChallenge to manage Tower of Hanoi challenges, providing a standardized interface for creating, updating, and deleting TowerChallenge instances.

Deployment:
To utilize the Tower of Hanoi Kubernetes Custom Resource and Composition setup:

    Apply the provided YAML manifests to your Kubernetes cluster to create the necessary CRDs, compositions, and composite resource definitions.
    Create instances of TowerChallenge custom resources with the desired number of discs to initiate Tower of Hanoi challenges.
    The defined compositions will automatically trigger the calculation of moves and update the solution in the tower-hanoi-moves ConfigMap.
    Monitor the status of TowerChallenge instances and tower-hanoi-moves ConfigMap for the progress and solution of Tower of Hanoi challenges.
