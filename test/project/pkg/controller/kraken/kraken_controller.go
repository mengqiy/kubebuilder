/*
Copyright 2019 The Kubernetes authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kraken

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	creaturesv2alpha1 "sigs.k8s.io/kubebuilder/test/project/pkg/apis/creatures/v2alpha1"
)

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Kraken Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileKraken{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	err := builder.ControllerManagedBy(mgr).
		For(&creaturesv2alpha1.Kraken{}).
		// TODO(user): Modify this to be the types you create
		Owns(&appsv1.Deployment{}).
		Complete(r)
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileKraken{}

// ReconcileKraken reconciles a Kraken object
type ReconcileKraken struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Kraken object and makes changes based on the state read
// and what is in the Kraken.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  The scaffolding writes
// a Deployment as an example
// +kubebuilder:rbac:groups=creatures.testproject.org,resources=krakens,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=creatures.testproject.org,resources=krakens/status,verbs=get;update;patch
func (r *ReconcileKraken) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the Kraken instance
	instance := &creaturesv2alpha1.Kraken{}
	err := r.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
