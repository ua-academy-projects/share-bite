package controller

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	guestv1alpha1 "github.com/ua-academy-projects/share-bite/operators/guest-operator/api/v1alpha1"
)

type GuestAppProfileReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *GuestAppProfileReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var profile guestv1alpha1.GuestAppProfile
	if err := r.Get(ctx, req.NamespacedName, &profile); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	deployName := profile.Spec.DeploymentName
	if deployName == "" {
		deployName = "guest-api"
	}

	var desiredReplicas int32 = 0
	if profile.Spec.Enabled {
		if profile.Spec.Replicas < 0 {
			if err := r.updateStatus(ctx, &profile, metav1.ConditionFalse, "InvalidSpec", "spec.replicas cannot be negative"); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		}
		desiredReplicas = profile.Spec.Replicas
	}

	var deploy appsv1.Deployment
	deployReq := types.NamespacedName{Name: deployName, Namespace: profile.Namespace}
	if err := r.Get(ctx, deployReq, &deploy); err != nil {
		if errors.IsNotFound(err) {
			logger.Info("Deployment not found, waiting...", "Deployment", deployName)
			if statusErr := r.updateStatus(ctx, &profile, metav1.ConditionFalse, "DeploymentNotFound", "Target deployment is missing"); statusErr != nil {
				return ctrl.Result{}, statusErr
			}
			return ctrl.Result{RequeueAfter: 5000000000}, nil
		}
		return ctrl.Result{}, err
	}

	if deploy.Spec.Replicas == nil || *deploy.Spec.Replicas != desiredReplicas {
		logger.Info("Updating Deployment replicas", "Current", deploy.Spec.Replicas, "Desired", desiredReplicas)
		deploy.Spec.Replicas = &desiredReplicas
		if err := r.Update(ctx, &deploy); err != nil {
			return ctrl.Result{}, err
		}
	}

	if deploy.Status.ReadyReplicas == desiredReplicas {
		if err := r.updateStatus(ctx, &profile, metav1.ConditionTrue, "Scaled", "Deployment reached desired replicas"); err != nil {
			return ctrl.Result{}, err
		}
	} else {
		if err := r.updateStatus(ctx, &profile, metav1.ConditionFalse, "Scaling", "Waiting for pods to be ready"); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{RequeueAfter: 3000000000}, nil
	}

	return ctrl.Result{}, nil
}

func (r *GuestAppProfileReconciler) updateStatus(ctx context.Context, profile *guestv1alpha1.GuestAppProfile, status metav1.ConditionStatus, reason, message string) error {
	condition := metav1.Condition{
		Type:               "Ready",
		Status:             status,
		Reason:             reason,
		Message:            message,
		LastTransitionTime: metav1.Now(),
	}
	meta.SetStatusCondition(&profile.Status.Conditions, condition)
	return r.Status().Update(ctx, profile)
}

func (r *GuestAppProfileReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&guestv1alpha1.GuestAppProfile{}).
		Complete(r)
}
