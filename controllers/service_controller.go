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

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	datav1alpha1 "github.com/MadeByMakers/kong-operator-for-k8s/api/v1alpha1"
	daopackage "github.com/MadeByMakers/kong-operator-for-k8s/dao"
	"github.com/go-logr/logr"
	"github.com/prometheus/common/log"
	"github.com/redhat-cop/operator-utils/pkg/util"
)

// ServiceReconciler reconciles a Service object
type ServiceReconciler struct {
	util.ReconcilerBase
	Log logr.Logger
	Dao daopackage.ServiceDAO
}

//+kubebuilder:rbac:groups=data.konghq.com,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=data.konghq.com,resources=services/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=data.konghq.com,resources=services/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Service object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *ServiceReconciler) Reconcile(context context.Context, req ctrl.Request) (ctrl.Result, error) {
	const controllerName = "ServiceController"
	log := r.Log.WithValues("Service", req.NamespacedName)
	r.Dao = daopackage.ServiceDAO{}

	instance := &datav1alpha1.Service{}
	err := r.GetInstance(context, req, instance)
	if err != nil {
		if err.Error() == "NOT_FOUND" {
			return reconcile.Result{}, nil
		}
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

		err := r.GetInstance(context, req, instance)
		if err != nil {
			return ctrl.Result{}, err
		}

		instance.Status.Code = result.Status.Code
		instance.Status.Message = result.Status.Message
		instance.Status.Response = result.Status.Response

		r.GetClient().Status().Update(context, instance)

		err = r.GetInstance(context, req, instance)
		if err != nil {
			return ctrl.Result{}, err
		}

		if !controllerutil.ContainsFinalizer(instance, controllerName) {
			controllerutil.AddFinalizer(instance, controllerName)
		}

		instance.Spec = result.Spec
		r.GetClient().Update(context, instance)

		return ctrl.Result{}, nil
	}
	return r.ManageSuccess(context, instance)
}

func (r *ServiceReconciler) GetInstance(context context.Context, req ctrl.Request, instance *datav1alpha1.Service) error {
	err := r.GetClient().Get(context, req.NamespacedName, instance)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return errors.New("NOT_FOUND")
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get Service.")
		return err
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&datav1alpha1.Service{}).
		Complete(r)
}

func (r *ServiceReconciler) IsValid(obj metav1.Object) (bool, error) {
	instance, ok := obj.(*datav1alpha1.Service)
	if !ok {
		return false, errors.New("not a Service object")
	}

	if instance.Spec.Name == "" {
		return false, errors.New("'name' cannot be empty")
	}

	if instance.Spec.Host == "" {
		return false, errors.New("'host' cannot be empty")
	}

	return true, nil
}

func (r *ServiceReconciler) manageCleanUpLogic(instance *datav1alpha1.Service) error {
	response := r.Dao.Delete(*instance)

	if response.Status.Code != 200 {
		return errors.New(response.Status.Message)
	}

	return nil
}

func (r *ServiceReconciler) manageOperatorLogic(instance *datav1alpha1.Service) (error, *datav1alpha1.Service) {
	var response datav1alpha1.Service

	// DELETE
	if instance.GetDeletionTimestamp() != nil {
		return r.manageCleanUpLogic(instance), nil
	} else {
		response = r.Dao.Save(*instance)
	}

	if response.Status.Code != 200 {
		return errors.New(response.Status.Message), nil
	}

	return nil, &response
}
