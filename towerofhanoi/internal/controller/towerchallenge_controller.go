package controller

import (
	"context"
	"errors"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;create;update;patch;delete

// Reconcile is part of the main Kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *TowerChallengeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var towerChallenge webappv1alpha1.TowerChallenge
	if err := r.Get(ctx, req.NamespacedName, &towerChallenge); err != nil {
		log.Error(err, "Unable to fetch TowerChallenge")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	startTime := time.Now()
	towerChallenge.Status.StartTime = metav1.Time{Time: startTime}

	if err := validateTowerChallenge(towerChallenge); err != nil {
		towerChallenge.Status.Phase = "Failed"
		towerChallenge.Status.ErrorMessage = err.Error()
		_ = r.Status().Update(ctx, &towerChallenge)
		return ctrl.Result{}, err
	}

	steps := solveHanoi(towerChallenge.Spec.Discs, "A", "C", "B")
	towerChallenge.Status.Steps = steps

	configMapNames := []string{}
	existingCMs := &corev1.ConfigMapList{}
	listOpts := []client.ListOption{
		client.InNamespace(req.Namespace),
		client.MatchingLabels{"challenge": towerChallenge.Name},
	}
	if err := r.List(ctx, existingCMs, listOpts...); err != nil {
		log.Error(err, "Unable to list ConfigMaps")
		return ctrl.Result{}, err
	}

	for i, step := range steps {
		cmName := fmt.Sprintf("%s-move-%d", towerChallenge.Name, i+1)
		cm, found := existingCMsMap[cmName]
		if found {
			cm.Data = map[string]string{"move": step}
			if err := r.Update(ctx, cm); err != nil {
				log.Error(err, "Failed to update ConfigMap for the move", "ConfigMap", cmName)
				return ctrl.Result{}, err
			}
		} else {
			cm = &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      cmName,
					Namespace: towerChallenge.Namespace,
					Labels:    map[string]string{"challenge": towerChallenge.Name},
				},
				Data: map[string]string{"move": step},
			}
			if err := r.Create(ctx, cm); err != nil {
				log.Error(err, "Failed to create ConfigMap for the move", "ConfigMap", cmName)
				return ctrl.Result{}, err
			}
		}
		configMapNames = append(configMapNames, cmName)
	}

	towerChallenge.Status.ConfigMapNames = configMapNames
	towerChallenge.Status.ConfigMapsCreated = true
	towerChallenge.Status.EndTime = metav1.Time{Time: time.Now()}

	if err := r.Status().Update(ctx, &towerChallenge); err != nil {
		log.Error(err, "Failed to update TowerChallenge status at the end")
		return ctrl.Result{}, err
	}

	log.Info("Reconciled TowerChallenge successfully")
	return ctrl.Result{}, nil
}

func validateTowerChallenge(tc webappv1alpha1.TowerChallenge) error {
	if tc.Spec.Discs <= 0 {
		return errors.New("the number of discs must be positive")
	}
	return nil
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
