/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	hanoiv1alpha1 "github.com/MosesKing/Towers-of-Hanoi/api/v1alpha1"
)

// TowerChallengeReconciler reconciles a TowerChallenge object
type TowerChallengeReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=hanoi.hanoi.com,resources=towerchallenges,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=hanoi.hanoi.com,resources=towerchallenges/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=hanoi.hanoi.com,resources=towerchallenges/finalizers,verbs=update

// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *TowerChallengeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.Log.WithValues("towerchallenge", req.NamespacedName)
	// Establish a logger instance for structured logging throughout the reconciliation process.
	// Log entries will be tagged with the resource type 'towerchallenge' and the specific instance being reconciled.
	// Fetch the TowerChallenge instance using the namespace and name from the Request object.
	
	towerChallenge := &hanoiv1alpha1.TowerChallenge{}
	err := r.Get(ctx, req.NamespacedName, towerChallenge)
	if err != nil {
		// Check if the TowerChallenge resource was not found in the cluster.
		if errors.IsNotFound(err) {
			// The resource could have been deleted after the reconcile request was queued.
			// Log a message indicating that the resource was not found, and return gracefully.
			// Kubernetes garbage collector handles the cleanup of all dependent resources automatically using owner references.
			log.Info("TowerChallenge resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Log the error encountered while trying to fetch the TowerChallenge resource,
		// then requeue the request to be tried again.
	log.Info("Handling TowerChallenge resource")

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TowerChallengeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&hanoiv1alpha1.TowerChallenge{}).
		Complete(r)
}
