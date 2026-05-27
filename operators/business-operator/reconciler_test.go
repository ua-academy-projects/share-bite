package businessoperator

import (
	"context"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestBusinessAppProfileReconciler_Reconcile(t *testing.T) {
	s := runtime.NewScheme()
	_ = appsv1.AddToScheme(s)
	_ = AddToScheme(s)

	type testCase struct {
		name           string
		initialObjects []client.Object
		reqName        string
		reqNamespace   string
		wantErr        bool
		validate       func(t *testing.T, fakeClient client.Client)
	}

	tests := []testCase{
		{
			name: "Happy Path: Scale from 1 to 3",
			initialObjects: []client.Object{
				&appsv1.Deployment{
					ObjectMeta: metav1.ObjectMeta{Name: "business-api", Namespace: "default"},
					Spec:       appsv1.DeploymentSpec{Replicas: func() *int32 { r := int32(1); return &r }()},
				},
				&BusinessAppProfile{
					ObjectMeta: metav1.ObjectMeta{Name: "test-profile", Namespace: "default"},
					Spec:       Spec{Replicas: 3, Enabled: true},
				},
			},
			reqName:      "test-profile",
			reqNamespace: "default",
			wantErr:      false,
			validate: func(t *testing.T, fakeClient client.Client) {
				var dep appsv1.Deployment
				err := fakeClient.Get(context.Background(), types.NamespacedName{Name: "business-api", Namespace: "default"}, &dep)
				if err != nil {
					t.Fatalf("Failed to get deployment: %v", err)
				}
				if *dep.Spec.Replicas != int32(3) {
					t.Errorf("Expected 3 replicas, got %d", *dep.Spec.Replicas)
				}
			},
		},
		{
			name: "Scale to Zero when Enabled is false",
			initialObjects: []client.Object{
				&appsv1.Deployment{
					ObjectMeta: metav1.ObjectMeta{Name: "business-api", Namespace: "default"},
					Spec:       appsv1.DeploymentSpec{Replicas: func() *int32 { r := int32(3); return &r }()},
				},
				&BusinessAppProfile{
					ObjectMeta: metav1.ObjectMeta{Name: "test-profile", Namespace: "default"},
					Spec:       Spec{Replicas: 3, Enabled: false},
				},
			},
			reqName:      "test-profile",
			reqNamespace: "default",
			wantErr:      false,
			validate: func(t *testing.T, fakeClient client.Client) {
				var dep appsv1.Deployment
				err := fakeClient.Get(context.Background(), types.NamespacedName{Name: "business-api", Namespace: "default"}, &dep)
				if err != nil {
					t.Fatalf("Failed to get deployment: %v", err)
				}
				if *dep.Spec.Replicas != int32(0) {
					t.Errorf("Expected 0 replicas, got: %d", *dep.Spec.Replicas)
				}
			},
		},
		{
			name: "Deployment Not Found sets Condition False",
			initialObjects: []client.Object{
				&BusinessAppProfile{
					ObjectMeta: metav1.ObjectMeta{Name: "test-profile", Namespace: "default"},
					Spec:       Spec{Replicas: 3, Enabled: true},
				},
			},
			reqName:      "test-profile",
			reqNamespace: "default",
			wantErr:      false,
			validate: func(t *testing.T, fakeClient client.Client) {
				var profile BusinessAppProfile
				err := fakeClient.Get(context.Background(), types.NamespacedName{Name: "test-profile", Namespace: "default"}, &profile)
				if err != nil {
					t.Fatalf("Failed to get profile: %v", err)
				}
				if len(profile.Status.Conditions) == 0 {
					t.Fatalf("Expected status conditions, got none")
				}

				cond := profile.Status.Conditions[0]
				if cond.Reason != "DeploymentNotFound" || cond.Status != metav1.ConditionFalse {
					t.Errorf("Expected Reason DeploymentNotFound and Status False, got %s and %s", cond.Reason, cond.Status)
				}
			},
		},
		{
			name: "Custom Deployment Name scaling",
			initialObjects: []client.Object{
				&appsv1.Deployment{
					ObjectMeta: metav1.ObjectMeta{Name: "custom-name", Namespace: "default"},
					Spec:       appsv1.DeploymentSpec{Replicas: func() *int32 { r := int32(1); return &r }()},
				},
				&BusinessAppProfile{
					ObjectMeta: metav1.ObjectMeta{Name: "test-profile", Namespace: "default"},
					Spec: Spec{
						Replicas:       3,
						Enabled:        true,
						DeploymentName: func() *string { s := "custom-name"; return &s }(),
					},
				},
			},
			reqName:      "test-profile",
			reqNamespace: "default",
			wantErr:      false,
			validate: func(t *testing.T, fakeClient client.Client) {
				var dep appsv1.Deployment
				err := fakeClient.Get(context.Background(), types.NamespacedName{Name: "custom-name", Namespace: "default"}, &dep)
				if err != nil {
					t.Fatalf("Failed to get custom deployment: %v", err)
				}
				if *dep.Spec.Replicas != int32(3) {
					t.Errorf("Expected 3 replicas, got %d", *dep.Spec.Replicas)
				}
			},
		},
		{
			name: "Status updates to Ready=True when readyReplicas match",
			initialObjects: []client.Object{
				&appsv1.Deployment{
					ObjectMeta: metav1.ObjectMeta{Name: "business-api", Namespace: "default"},
					Spec:       appsv1.DeploymentSpec{Replicas: func() *int32 { r := int32(3); return &r }()},
					Status:     appsv1.DeploymentStatus{ReadyReplicas: 3},
				},
				&BusinessAppProfile{
					ObjectMeta: metav1.ObjectMeta{Name: "test-profile", Namespace: "default"},
					Spec:       Spec{Replicas: 3, Enabled: true},
				},
			},
			reqName:      "test-profile",
			reqNamespace: "default",
			wantErr:      false,
			validate: func(t *testing.T, fakeClient client.Client) {
				var profile BusinessAppProfile
				err := fakeClient.Get(context.Background(), types.NamespacedName{Name: "test-profile", Namespace: "default"}, &profile)
				if err != nil {
					t.Fatalf("Failed to get profile: %v", err)
				}
				if len(profile.Status.Conditions) == 0 {
					t.Fatalf("Expected status conditions, got none")
				}
				cond := profile.Status.Conditions[0]
				if cond.Status != metav1.ConditionTrue || cond.Reason != "Scaled" {
					t.Errorf("Expected Status True and Reason Scaled, got %s and %s", cond.Status, cond.Reason)
				}
			},
		},
		{
			name:           "CRD Not Found (Graceful exit on delete event)",
			initialObjects: []client.Object{},
			reqName:        "deleted-profile",
			reqNamespace:   "default",
			wantErr:        false,
			validate: func(t *testing.T, fakeClient client.Client) {

			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeClient := fake.NewClientBuilder().
				WithScheme(s).
				WithStatusSubresource(&BusinessAppProfile{}).
				WithObjects(tt.initialObjects...).
				Build()

			reconciler := BusinessAppProfileReconciler{
				Client: fakeClient,
			}

			req := ctrl.Request{
				NamespacedName: types.NamespacedName{
					Name:      tt.reqName,
					Namespace: tt.reqNamespace,
				},
			}

			_, err := reconciler.Reconcile(context.Background(), req)

			if (err != nil) != tt.wantErr {
				t.Errorf("Reconcile() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.validate != nil {
				tt.validate(t, fakeClient)
			}
		})
	}
}
