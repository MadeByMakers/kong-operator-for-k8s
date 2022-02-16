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
	"encoding/json"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	datav1alpha1 "github.com/MadeByMakers/kong-operator-for-k8s/api/v1alpha1"
	util "github.com/MadeByMakers/kong-operator-for-k8s/util"
)

// PluginReconciler reconciles a Plugin object
type PluginReconciler struct {
	client.Client
	Scheme *runtime.Scheme
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
func (r *PluginReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.Log.WithValues("memcached", req.NamespacedName)
	logger.Info("Reconciling Memcached")

	plugin := &datav1alpha1.Plugin{}
	err := r.Get(ctx, req.NamespacedName, plugin)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			logger.Info("Memcached resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		logger.Error(err, "Failed to get Memcached.")
		return ctrl.Result{}, err
	}

	var bodyBytes []byte

	// DELETE
	if plugin.GetDeletionTimestamp() != nil {
		bodyBytes = util.DoDelete("https://kong-kong-admin.kong.svc:8444/plugin/" + plugin.Spec.Id)
	} else if plugin.Spec.Id != "" {
		bodyBytes = util.DoPost("https://kong-kong-admin.kong.svc:8444/plugin", plugin.Spec)
	} else {
		bodyBytes = util.DoPut("https://kong-kong-admin.kong.svc:8444/plugin", plugin.Spec)
	}

	var responseObject datav1alpha1.PluginSpec
	json.Unmarshal(bodyBytes, &responseObject)

	r.Status().Update(ctx, plugin)

	fmt.Printf("API Response as struct %+v\n", responseObject)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PluginReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&datav1alpha1.Plugin{}).
		Complete(r)
}
