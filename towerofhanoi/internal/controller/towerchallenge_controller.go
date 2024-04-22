package controller

import (
	"context"
	"errors"
	"fmt"
	"time"

	webappv1alpha1 "hanoi.com/towerofhanoi/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type TowerChallengeReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *TowerChallengeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var towerChallenge webappv1alpha1.TowerChallenge
	if err := r.Get(ctx, req.NamespacedName, &towerChallenge); err != nil {
		log.Error(err, "Unable to fetch TowerChallenge")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	startTime := time.Now()
	if towerChallenge.Status.StartTime.IsZero() {
		towerChallenge.Status.StartTime = metav1.Time{Time: startTime}
	}

	if err := validateTowerChallenge(towerChallenge); err != nil {
		towerChallenge.Status.Phase = "Failed"
		towerChallenge.Status.ErrorMessage = err.Error()
		_ = r.Status().Update(ctx, &towerChallenge)
		return ctrl.Result{}, err
	}

	steps := solveHanoi(towerChallenge.Spec.Discs, "A", "C", "B")
	configMapNames := manageConfigMaps(ctx, r, req.Namespace, towerChallenge, steps)
	validNames := make(map[string]bool)
	for _, name := range configMapNames {
		validNames[name] = true
	}

	if err := cleanupOldConfigMaps(ctx, r, req.Namespace, towerChallenge, validNames); err != nil {
		log.Error(err, "Failed to clean up old ConfigMaps")
		return ctrl.Result{}, err
	}

	towerChallenge.Status.ConfigMapNames = configMapNames
	towerChallenge.Status.Phase = "Completed"
	towerChallenge.Status.EndTime = metav1.Time{Time: time.Now()}

	if err := r.Status().Update(ctx, &towerChallenge); err != nil {
		log.Error(err, "Failed to update TowerChallenge status")
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

func manageConfigMaps(ctx context.Context, r *TowerChallengeReconciler, namespace string, tc webappv1alpha1.TowerChallenge, steps []string) []string {
	var configMapNames []string
	existingCMs := &corev1.ConfigMapList{}
	listOpts := []client.ListOption{
		client.InNamespace(namespace),
		client.MatchingLabels{"challenge": tc.Name},
	}
	if err := r.List(ctx, existingCMs, listOpts...); err != nil {
		log.FromContext(ctx).Error(err, "Unable to list ConfigMaps")
		return nil
	}

	existingCMsMap := make(map[string]*corev1.ConfigMap)
	for _, cm := range existingCMs.Items {
		existingCMsMap[cm.Name] = cm.DeepCopy()
	}

	for i, step := range steps {
		cmName := fmt.Sprintf("%s-move-%d", tc.Name, i+1)
		cm, found := existingCMsMap[cmName]
		if found {
			// Refetch the latest version of the ConfigMap to ensure updates are applied on the latest version
			latestCM := &corev1.ConfigMap{}
			err := r.Get(ctx, client.ObjectKey{Name: cmName, Namespace: namespace}, latestCM)
			if err != nil {
				log.FromContext(ctx).Error(err, "Failed to fetch the latest version of ConfigMap", "ConfigMap", cmName)
				continue // skip this iteration if we cannot fetch the latest version
			}
			latestCM.Data = map[string]string{"move": step}
			if err := r.Update(ctx, latestCM); err != nil {
				if kerrors.IsConflict(err) {
					log.FromContext(ctx).Info("Conflict detected, retrying update", "ConfigMap", cmName)
					continue // Optionally, add a limit to retries or delay retries
				}
				log.FromContext(ctx).Error(err, "Failed to update ConfigMap", "ConfigMap", cmName)
				return nil // Return or handle the error appropriately
			}
		} else {
			cm = &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      cmName,
					Namespace: namespace,
					Labels:    map[string]string{"challenge": tc.Name},
				},
				Data: map[string]string{"move": step},
			}
			if err := r.Create(ctx, cm); err != nil {
				log.FromContext(ctx).Error(err, "Failed to create ConfigMap", "ConfigMap", cmName)
			}
		}
		configMapNames = append(configMapNames, cmName)
	}

	return configMapNames
}

func cleanupOldConfigMaps(ctx context.Context, r *TowerChallengeReconciler, namespace string, tc webappv1alpha1.TowerChallenge, validNames map[string]bool) error {
	var allConfigMaps corev1.ConfigMapList
	listOpts := []client.ListOption{
		client.InNamespace(namespace),
		client.MatchingLabels{"challenge": tc.Name},
	}
	if err := r.List(ctx, &allConfigMaps, listOpts...); err != nil {
		log.FromContext(ctx).Error(err, "Failed to list ConfigMaps for cleanup")
		return err
	}

	for _, cm := range allConfigMaps.Items {
		if _, isValid := validNames[cm.Name]; !isValid {
			// If the ConfigMap name is not in the list of valid names, delete it
			if err := r.Delete(ctx, &cm); err != nil {
				if !kerrors.IsNotFound(err) {
					log.FromContext(ctx).Error(err, "Failed to delete ConfigMap", "ConfigMap", cm.Name)
					continue // continue with the next item
				}
			}
			log.FromContext(ctx).Info("Deleted old or invalid ConfigMap", "ConfigMap", cm.Name)
		}
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

func (r *TowerChallengeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webappv1alpha1.TowerChallenge{}).
		Owns(&corev1.ConfigMap{}).
		Complete(r)
}
