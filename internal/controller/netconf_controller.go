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
	"os"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	networkv1 "github.com/saeed-mcu/network-operator/api/v1"

	"github.com/saeed-mcu/network-operator/api/shared"
	"github.com/saeed-mcu/network-operator/pkg/nmstatectl"
)

const (
	finalizerName = "ae.digicloud/netplan"

	defaultGwProbeTimeout = 60 * time.Second
	apiServerProbeTimeout = 60 * time.Second
	// DesiredStateConfigurationTimeout doubles the default gw ping probe and API server
	// connectivity check timeout to ensure the Checkpoint is alive before rolling it back
	// https://nmstate.github.io/cli_guide#manual-transaction-control
	DesiredStateConfigurationTimeout = (defaultGwProbeTimeout + apiServerProbeTimeout)
)

// NetConfReconciler reconciles a NetConf object
type NetConfReconciler struct {
	client.Client
	APIClient client.Client
	Scheme    *runtime.Scheme
}

// +kubebuilder:rbac:groups=network.digicloud.ir,resources=netconfs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=network.digicloud.ir,resources=netconfs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=network.digicloud.ir,resources=netconfs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the NetConf object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *NetConfReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	nodeName := os.Getenv("NODE_NAME")

	_, err := nmstatectl.Version()
	if err != nil {
		logger.Error(err, "failed retrieving nmstate version")
		return ctrl.Result{}, err
	}

	instance := &networkv1.NetConf{}
	err = r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, err
	}

	if instance.Spec.NodeName != nodeName {
		return ctrl.Result{}, nil
	}

	// --- Handle Deletion ---
	if !instance.ObjectMeta.DeletionTimestamp.IsZero() {
		logger.Info("Resource is being deleted")
		if controllerutil.ContainsFinalizer(instance, finalizerName) {
			if err := r.cleanupResource(ctx, instance); err != nil {
				return ctrl.Result{}, err
			}

			// Remove finalizer to allow deletion
			controllerutil.RemoveFinalizer(instance, finalizerName)
			if err := r.Update(ctx, instance); err != nil {
				return ctrl.Result{}, err
			}
			logger.Info("RemoveFinalizer")
		}
		return ctrl.Result{}, nil
	}

	desiredState := shared.NewState(instance.Spec.NetConf)
	nmstateOutput, err := ApplyDesiredState(ctx, r.APIClient, desiredState)

	if err != nil {
		logger.Info("nmstate", "output", nmstateOutput)
		instance.Status.Applied = "False"
		instance.Status.State = err.Error()

		meta.SetStatusCondition(&instance.Status.Conditions, metav1.Condition{
			Type:               "OperatorDegraded",
			Status:             metav1.ConditionTrue,
			Reason:             networkv1.ReasonOperandDeploymentFailed,
			LastTransitionTime: metav1.NewTime(time.Now()),
			Message:            "Operator Failed",
		})

		r.Status().Update(ctx, instance)
		return ctrl.Result{}, nil
	}

	// Add finalizer if not present
	if !controllerutil.ContainsFinalizer(instance, finalizerName) {
		controllerutil.AddFinalizer(instance, finalizerName)
		if err := r.Update(ctx, instance); err != nil {
			logger.Info("Add finalizer")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	meta.SetStatusCondition(&instance.Status.Conditions, metav1.Condition{
		Type:               "OperatorDegraded",
		Status:             metav1.ConditionTrue,
		Reason:             networkv1.ReasonSucceeded,
		LastTransitionTime: metav1.NewTime(time.Now()),
		Message:            "Operator successfully reconciling",
	})

	instance.Status.Applied = "True"
	instance.Status.State = networkv1.NoError
	r.Status().Update(ctx, instance)
	logger.Info("Apply nmstate Done !!!")
	return ctrl.Result{Requeue: false}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NetConfReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&networkv1.NetConf{}).
		Complete(r)
}

func (r *NetConfReconciler) cleanupResource(ctx context.Context, instance *networkv1.NetConf) error {

	logger := log.FromContext(ctx)
	logger.Info("Cleanup Done")
	return nil
}

func ApplyDesiredState(ctx context.Context, cli client.Client, desiredState shared.State) (string, error) {

	logger := log.FromContext(ctx)

	if string(desiredState.Raw) == "" {
		return "Ignoring empty desired state", nil
	}

	logger.Info("ApplyDesiredState", "desiredState", string(desiredState.Raw))

	// Before apply we get the probes that are working fine, they should be
	// working fine after apply
	//probes := probe.Select(cli)

	// Rollback before Apply to remove pending checkpoints (for example handler pod restarted
	// before Commit)
	nmstatectl.Rollback()

	setOutput, err := nmstatectl.Set(desiredState, DesiredStateConfigurationTimeout)
	if err != nil {
		logger.Info("nmstatectl.Set", "setOutput", setOutput, "err", err.Error())
		return setOutput, err
	}
	logger.Info("nmstatectl.Set DONE")

	//err = probe.Run(cli, probes)
	//if err != nil {
	//	return "", rollback(cli, probes, errors.Wrap(err, "failed runnig probes after network changes"))
	//}

	commitOutput, err := nmstatectl.Commit()
	if err != nil {
		// We cannot rollback if commit fails, just return the error
		logger.Info("nmstatectl.Commit", "commitOutput", commitOutput, "err", err.Error())
		return commitOutput, err
	}
	logger.Info("nmstatectl.Commit Done")

	commandOutput := fmt.Sprintf("setOutput: %s \n", setOutput)
	return commandOutput, nil
}
