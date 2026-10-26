package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	respawnv1alpha1 "github.com/bblaker/respawn/api/v1alpha1"
)

// GameServerReconciler reconciles a GameServer object
type GameServerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=respawn.bblaker.com,resources=gameservers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=respawn.bblaker.com,resources=gameservers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=respawn.bblaker.com,resources=gameservers/finalizers,verbs=update

// Reconcile moves the cluster state toward the desired state of a GameServer.
// It is a no-op until M2 (core lifecycle), M3 (sleep), and M4 (wake); see section 4.4 of docs/PROJECT_PLAN.md.
func (r *GameServerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = logf.FromContext(ctx)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GameServerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&respawnv1alpha1.GameServer{}).
		Named("gameserver").
		Complete(r)
}
