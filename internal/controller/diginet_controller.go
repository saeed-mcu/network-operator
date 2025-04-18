/*
Copyright 2025.

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

package controller

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"

	networkv1 "github.com/saeed-mcu/network-operator/api/v1"
)

// DigiNetReconciler reconciles a DigiNet object
type DigiNetReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=network.digicloud.ir,resources=diginets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=network.digicloud.ir,resources=diginets/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=network.digicloud.ir,resources=diginets/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DigiNet object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *DigiNetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	logger := log.FromContext(ctx)

	digiNet := &networkv1.DigiNet{}
	err := r.Get(ctx, req.NamespacedName, digiNet)
	if err != nil && errors.IsNotFound(err) {
		return ctrl.Result{}, nil
	} else if err != nil {

		meta.SetStatusCondition(&digiNet.Status.Conditions, metav1.Condition{
			Type:               "OperatorDegraded",
			Status:             metav1.ConditionTrue,
			Reason:             networkv1.ReasonDeploymentNotAvailable,
			LastTransitionTime: metav1.NewTime(time.Now()),
			Message:            fmt.Sprintf("unable to get operator custom resource: %s", err.Error()),
		})

		return ctrl.Result{}, utilerrors.NewAggregate([]error{err, r.Status().Update(ctx, netConfig)})
	}

	digiNet.Status.Applied = "False"
	r.Status().Update(ctx, digiNet)

	meta.SetStatusCondition(&digiNet.Status.Conditions, metav1.Condition{
		Type:               "OperatorDegraded",
		Status:             metav1.ConditionTrue,
		Reason:             networkv1.ReasonSucceeded,
		LastTransitionTime: metav1.NewTime(time.Now()),
		Message:            "Operator successfully reconciling",
	})

	digiNet.Status.Applied = "True"
	logger.Info("DigiNet Done")
	return ctrl.Result{}, utilerrors.NewAggregate([]error{err, r.Status().Update(ctx, digiNet)})
}

// SetupWithManager sets up the controller with the Manager.
func (r *DigiNetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&networkv1.DigiNet{}).
		Complete(r)
}
