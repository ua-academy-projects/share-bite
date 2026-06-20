package controller

import (
	"context"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	adminv1alpha1 "github.com/ua-academy-projects/share-bite/operators/admin-operator/api/v1alpha1"
)

func TestReconcile_EnabledFalse(t *testing.T) {
	s := scheme.Scheme
	_ = adminv1alpha1.AddToScheme(s)

	profile := &adminv1alpha1.AdminAppProfile{
		ObjectMeta: metav1.ObjectMeta{Name: "test-profile", Namespace: "default"},
		Spec:       adminv1alpha1.AdminAppProfileSpec{Replicas: 3, Enabled: false},
	}

	var initialReplicas int32 = 3
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: "admin-auth-api", Namespace: "default"},
		Spec:       appsv1.DeploymentSpec{Replicas: &initialReplicas},
	}

	cl := fake.NewClientBuilder().WithScheme(s).WithObjects(profile, deployment).WithStatusSubresource(profile).Build()
	r := &AdminAppProfileReconciler{Client: cl, Scheme: s}

	req := ctrl.Request{NamespacedName: types.NamespacedName{Name: "test-profile", Namespace: "default"}}
	_, err := r.Reconcile(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updatedDep := &appsv1.Deployment{}
	_ = cl.Get(context.Background(), types.NamespacedName{Name: "admin-auth-api", Namespace: "default"}, updatedDep)

	if *updatedDep.Spec.Replicas != 0 {
		t.Errorf("expected 0 replicas, got %d", *updatedDep.Spec.Replicas)
	}
}

func TestReconcile_MissingDeployment(t *testing.T) {
	s := scheme.Scheme
	_ = adminv1alpha1.AddToScheme(s)

	profile := &adminv1alpha1.AdminAppProfile{
		ObjectMeta: metav1.ObjectMeta{Name: "test-profile", Namespace: "default"},
		Spec:       adminv1alpha1.AdminAppProfileSpec{Replicas: 2, Enabled: true},
	}

	cl := fake.NewClientBuilder().WithScheme(s).WithObjects(profile).WithStatusSubresource(profile).Build()
	r := &AdminAppProfileReconciler{Client: cl, Scheme: s}

	req := ctrl.Request{NamespacedName: types.NamespacedName{Name: "test-profile", Namespace: "default"}}
	res, err := r.Reconcile(context.Background(), req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.RequeueAfter != 5000000000 {
		t.Errorf("expected RequeueAfter 5s, got %v", res.RequeueAfter)
	}
}

func TestReconcile_HappyPath(t *testing.T) {
	s := scheme.Scheme
	_ = adminv1alpha1.AddToScheme(s)

	profile := &adminv1alpha1.AdminAppProfile{
		ObjectMeta: metav1.ObjectMeta{Name: "test-profile", Namespace: "default"},
		Spec:       adminv1alpha1.AdminAppProfileSpec{Replicas: 3, Enabled: true},
	}

	var initialReplicas int32 = 1
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: "admin-auth-api", Namespace: "default"},
		Spec:       appsv1.DeploymentSpec{Replicas: &initialReplicas},
	}

	cl := fake.NewClientBuilder().WithScheme(s).WithObjects(profile, deployment).WithStatusSubresource(profile).Build()
	r := &AdminAppProfileReconciler{Client: cl, Scheme: s}

	req := ctrl.Request{NamespacedName: types.NamespacedName{Name: "test-profile", Namespace: "default"}}
	_, err := r.Reconcile(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updatedDep := &appsv1.Deployment{}
	_ = cl.Get(context.Background(), types.NamespacedName{Name: "admin-auth-api", Namespace: "default"}, updatedDep)

	if *updatedDep.Spec.Replicas != 3 {
		t.Errorf("expected 3 replicas, got %d", *updatedDep.Spec.Replicas)
	}
}
