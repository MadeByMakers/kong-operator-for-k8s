/*
Copyright 2022.

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

package controllers

import (
	"context"
	"errors"

	"github.com/go-logr/logr"
	"github.com/redhat-cop/operator-utils/pkg/util"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	datav1alpha1 "github.com/MadeByMakers/kong-operator-for-k8s/api/v1alpha1"
	daopackage "github.com/MadeByMakers/kong-operator-for-k8s/dao"
)

// PluginReconciler reconciles a Plugin object
type PluginReconciler struct {
	util.ReconcilerBase
	Log logr.Logger
	Dao daopackage.PluginDAO
}

//+kubebuilder:rbac:groups=data.data.konghq.com,resources=plugins,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=data.data.konghq.com,resources=plugins/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=data.data.konghq.com,resources=plugins/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Plugin object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *PluginReconciler) Reconcile(context context.Context, req ctrl.Request) (ctrl.Result, error) {
	const controllerName = "PluginController"
	log := r.Log.WithValues("Plugin", req.NamespacedName)
	r.Dao = daopackage.PluginDAO{}

	instance := &datav1alpha1.Plugin{}
	err := r.GetClient().Get(context, req.NamespacedName, instance)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("Plugin resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get Memcached.")
		return ctrl.Result{}, err
	}

	if ok, err := r.IsValid(instance); !ok {
		return r.ManageError(context, instance, err)
	}

	if util.IsBeingDeleted(instance) {
		if !util.HasFinalizer(instance, controllerName) {
			return reconcile.Result{}, nil
		}
		err := r.manageCleanUpLogic(instance)
		if err != nil {
			log.Error(err, "unable to delete instance", "instance", instance)
			return r.ManageError(context, instance, err)
		}
		util.RemoveFinalizer(instance, controllerName)
		err = r.GetClient().Update(context, instance)
		if err != nil {
			log.Error(err, "unable to update instance", "instance", instance)
			return r.ManageError(context, instance, err)
		}
		return reconcile.Result{}, nil
	}

	err, result := r.manageOperatorLogic(instance)
	if err != nil {
		return r.ManageError(context, instance, err)
	}

	if result != nil {
		r.GetClient().Update(context, result)
	}

	return r.ManageSuccess(context, instance)
}

// SetupWithManager sets up the controller with the Manager.
func (r *PluginReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&datav1alpha1.Plugin{}).
		Complete(r)
}

func (r *PluginReconciler) IsValid(obj metav1.Object) (bool, error) {
	instance, ok := obj.(*datav1alpha1.Plugin)
	if !ok {
		return false, errors.New("not a Plugin object")
	}

	if instance.Spec.Name == "" {
		return false, errors.New("'name' cannot be empty")
	}

	if instance.Spec.Config == "" {
		return false, errors.New("'config' cannot be empty")
	}

	return true, nil
}

func (r *PluginReconciler) manageCleanUpLogic(instance *datav1alpha1.Plugin) error {
	response := r.Dao.Delete(*instance)

	if response.Status.Code == 200 {
		return errors.New(response.Status.Message)
	}

	return nil
}

func (r *PluginReconciler) manageOperatorLogic(instance *datav1alpha1.Plugin) (error, *datav1alpha1.Plugin) {
	var response datav1alpha1.Plugin

	// DELETE
	if instance.GetDeletionTimestamp() != nil {
		return r.manageCleanUpLogic(instance), nil
	} else if instance.Spec.Id != "" {
		response = r.Dao.Update(*instance)
	} else {
		response = r.Dao.Create(*instance)
	}

	if response.Status.Code != 200 {
		return errors.New(response.Status.Message), nil
	}

	return nil, &response
}
