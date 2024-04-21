package controller

import (
	"context"
	"fmt"

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

//+kubebuilder:rbac:groups=hanoi.com,resources=towerchallenges,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=hanoi.com,resources=towerchallenges/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=hanoi.com,resources=towerchallenges/finalizers,verbs=update

// Reconcile handles the actual reconciliation logic of the TowerChallenge controller.
func (r *TowerChallengeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.Log.WithValues("towerchallenge", req.NamespacedName)

	// Fetch the TowerChallenge instance using the namespace and name from the Request object.
	towerChallenge := &hanoiv1alpha1.TowerChallenge{}
	err := r.Get(ctx, req.NamespacedName, towerChallenge)
	if err != nil {
		if errors.IsNotFound(err) {
			// The resource could have been deleted after the reconcile request was queued.
			// Log a message indicating that the resource was not found, and return gracefully.
			// Kubernetes garbage collector handles the cleanup of all dependent resources automatically using owner references.
			log.Info("TowerChallenge resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Log the error encountered while trying to fetch the TowerChallenge resource,
		// then requeue the request to be tried again.
		log.Error(err, "Failed to get TowerChallenge")
		return ctrl.Result{}, err
	}

	// Log a message indicating that the TowerChallenge resource is being handled.
	log.Info("Handling TowerChallenge resource")

	// Implement Tower of Hanoi algorithm here
	moves := TowerOfHanoi(towerChallenge.Spec.NumDisks, "Source", "Auxiliary", "Destination")

	// Update ConfigMaps based on game state
	// For now, let's just log the moves
	for _, move := range moves {
		log.Info(move)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TowerChallengeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&hanoiv1alpha1.TowerChallenge{}).
		Complete(r)
}

// TowerOfHanoi implements the Tower of Hanoi algorithm recursively
func TowerOfHanoi(numDisks int, source, auxiliary, destination string) []string {
	var moves []string
	if numDisks == 1 {
		moves = append(moves, fmt.Sprintf("Move disk 1 from %s to %s", source, destination))
		return moves
	}
	moves = append(moves, TowerOfHanoi(numDisks-1, source, destination, auxiliary)...)
	moves = append(moves, fmt.Sprintf("Move disk %d from %s to %s", numDisks, source, destination))
	moves = append(moves, TowerOfHanoi(numDisks-1, auxiliary, source, destination)...)
	return moves
}
