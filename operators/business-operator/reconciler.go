package businessoperator

import (
	"context"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Spec struct {
	Replicas       int32   `json:"replicas"`
	Enabled        bool    `json:"enabled"`
	DeploymentName *string `json:"deploymentName,omitempty"`
}

type Status struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

type BusinessAppProfile struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   Spec   `json:"spec"`
	Status Status `json:"status,omitempty"`
}

type BusinessAppProfileList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BusinessAppProfile `json:"items"`
}

type BusinessAppProfileReconciler struct {
	Client client.Client
}

func (in *BusinessAppProfile) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}

	out := &BusinessAppProfile{
		TypeMeta: in.TypeMeta,
	}

	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec.Replicas = in.Spec.Replicas
	out.Spec.Enabled = in.Spec.Enabled

	var temp string
	if in.Spec.DeploymentName != nil {
		temp = *in.Spec.DeploymentName
		out.Spec.DeploymentName = &temp
	}

	if in.Status.Conditions != nil {
		out.Status.Conditions = make([]metav1.Condition, len(in.Status.Conditions))
		copy(out.Status.Conditions, in.Status.Conditions)
	}

	return out
}

func (in *BusinessAppProfileList) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}

	out := &BusinessAppProfileList{
		TypeMeta: in.TypeMeta,
	}
	in.ListMeta.DeepCopyInto(&out.ListMeta)

	if in.Items != nil {
		out.Items = make([]BusinessAppProfile, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}
	return out
}

func (in *BusinessAppProfile) DeepCopyInto(out *BusinessAppProfile) {
	clone := in.DeepCopyObject().(*BusinessAppProfile)
	*out = *clone
}

func (r *BusinessAppProfileReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var profile BusinessAppProfile
	if err := r.Client.Get(ctx, req.NamespacedName, &profile); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	targetName := profile.Name
	if profile.Spec.DeploymentName != nil {
		targetName = *profile.Spec.DeploymentName
	}

	var dep appsv1.Deployment

	if err := r.Client.Get(ctx, types.NamespacedName{Namespace: req.Namespace, Name: targetName}, &dep); err != nil {
		if !errors.IsNotFound(err) {
			return ctrl.Result{}, err
		}
		condition := metav1.Condition{
			Type:               "Ready",
			Status:             metav1.ConditionFalse,
			Reason:             "DeploymentNotFound",
			LastTransitionTime: metav1.Now(),
		}

		profile.Status.Conditions = []metav1.Condition{condition}
		err := r.Client.Status().Update(ctx, &profile)
		return ctrl.Result{}, err
	}

	targetReplicas := int32(0)
	if profile.Spec.Enabled == true {
		targetReplicas = profile.Spec.Replicas
	}

	if dep.Spec.Replicas == nil || *dep.Spec.Replicas != targetReplicas {
		patch := client.MergeFrom(dep.DeepCopy())
		dep.Spec.Replicas = &targetReplicas
		if err := r.Client.Patch(ctx, &dep, patch); err != nil {
			return ctrl.Result{}, err
		}
	}

	var condition metav1.Condition
	res := ctrl.Result{}

	if dep.Status.ReadyReplicas == targetReplicas {
		condition = metav1.Condition{
			Type:               "Ready",
			Status:             metav1.ConditionTrue,
			Reason:             "Scaled",
			LastTransitionTime: metav1.Now(),
		}
	} else {
		condition = metav1.Condition{
			Type:               "Ready",
			Status:             metav1.ConditionFalse,
			Reason:             "Scaling",
			LastTransitionTime: metav1.Now(),
		}
		res = ctrl.Result{RequeueAfter: 5 * time.Second}
	}

	profile.Status.Conditions = []metav1.Condition{condition}
	err := r.Client.Status().Update(ctx, &profile)
	return res, err
}
