package main

import (
	"log"

	bo "github.com/ua-academy-projects/share-bite/operators/business-operator"
	appsv1 "k8s.io/api/apps/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func main() {
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))
	cfg := ctrl.GetConfigOrDie()

	mngr, err := ctrl.NewManager(cfg, ctrl.Options{LeaderElection: true, LeaderElectionID: "business-operator-lock", LeaderElectionNamespace: "default"})
	if err != nil {
		log.Fatalf("Error creating config manager: %v", err)
	}
	if err := bo.AddToScheme(mngr.GetScheme()); err != nil {
		log.Fatalf("Error adding scheme: %v", err)
	}

	client := mngr.GetClient()

	reconciler := bo.BusinessAppProfileReconciler{
		Client: client,
	}
	err = ctrl.NewControllerManagedBy(mngr).For(&bo.BusinessAppProfile{}).Owns(&appsv1.Deployment{}).Complete(&reconciler)
	if err != nil {
		log.Fatalf("Error setting up controller: %v", err)
	}

	err = mngr.Start(ctrl.SetupSignalHandler())
	if err != nil {
		log.Fatalf("Error starting manager: %v", err)
	}
}
