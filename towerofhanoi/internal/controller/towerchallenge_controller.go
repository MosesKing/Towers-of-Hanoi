package controller

import (
	"context"
	"errors"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	webappv1alpha1 "hanoi.com/towerofhanoi/api/v1alpha1"
)

// TowerChallengeReconciler reconciles a TowerChallenge object
type TowerChallengeReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=webapp.hanoi.com,resources=towerchallenges,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=webapp.hanoi.com,resources=towerchallenges/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=webapp.hanoi.com,resources=towerchallenges/finalizers,verbs=update

// Reconcile is part of the main Kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *TowerChallengeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var towerChallenge webappv1alpha1.TowerChallenge
	if err := r.Get(ctx, req.NamespacedName, &towerChallenge); err != nil {
		log.Error(err, "Unable to fetch TowerChallenge")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Validate the TowerChallenge specs
	if err := validateTowerChallenge(towerChallenge); err != nil {
		log.Error(err, "Invalid TowerChallenge specifications")
		towerChallenge.Status.Message = "Error: " + err.Error()
		_ = r.Status().Update(ctx, &towerChallenge)
		return ctrl.Result{}, err
	}

	// Solve Tower of Hanoi
	steps := solveHanoi(towerChallenge.Spec.Discs, "A", "C", "B")

	// Update status with the solution steps
	towerChallenge.Status.Steps = steps
	if err := r.Status().Update(ctx, &towerChallenge); err != nil {
		log.Error(err, "Failed to update TowerChallenge status")
		return ctrl.Result{}, err
	}

	log.Info("Reconciled TowerChallenge", "steps", steps)
	return ctrl.Result{}, nil
}

// validateTowerChallenge ensures that the TowerChallenge resource has valid data
func validateTowerChallenge(tc webappv1alpha1.TowerChallenge) error {
	if tc.Spec.Discs <= 0 {
		return errors.New("the number of discs must be positive")
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TowerChallengeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webappv1alpha1.TowerChallenge{}).
		Complete(r)
}

func solveHanoi(n int, from, to, aux string) []string {
	if n == 1 {
		return []string{fmt.Sprintf("Move disk 1 from %s to %s", from, to)}
	}
	var moves []string
	moves = append(moves, solveHanoi(n-1, from, aux, to)...)
	moves = append(moves, fmt.Sprintf("Move disk %d from %s to %s", n, from, to))
	moves = append(moves, solveHanoi(n-1, aux, to, from)...)
	return moves
}
