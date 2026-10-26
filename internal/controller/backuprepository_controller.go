package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	respawnv1alpha1 "github.com/bblaker/respawn/api/v1alpha1"
)

// BackupRepositoryReconciler reconciles a BackupRepository object
type BackupRepositoryReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=respawn.bblaker.com,resources=backuprepositories,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=respawn.bblaker.com,resources=backuprepositories/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=respawn.bblaker.com,resources=backuprepositories/finalizers,verbs=update

// Reconcile moves the cluster state toward the desired state of a BackupRepository.
// It is a no-op until M5; see section 4.8 of docs/PROJECT_PLAN.md.
func (r *BackupRepositoryReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = logf.FromContext(ctx)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *BackupRepositoryReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&respawnv1alpha1.BackupRepository{}).
		Named("backuprepository").
		Complete(r)
}
