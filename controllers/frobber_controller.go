/*


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

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	frobsv1alpha1 "danielfbm.github.io/k8s-resource-version/api/v1alpha1"
)

// FrobberReconciler reconciles a Frobber object
type FrobberReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=frobs.danielfbm.github.io,resources=frobbers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=frobs.danielfbm.github.io,resources=frobbers/status,verbs=get;update;patch

func (r *FrobberReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("frobber", req.NamespacedName)

	frobber := &frobsv1alpha1.Frobber{}
	err := r.Client.Get(ctx, req.NamespacedName, frobber)
	if err != nil {
		err = client.IgnoreNotFound(err)
		return ctrl.Result{}, err
	}
	// your logic here
	log.Info("Print", "height", frobber.Spec.Height, "param", frobber.Spec.Param)

	return ctrl.Result{}, nil
}

func (r *FrobberReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&frobsv1alpha1.Frobber{}).
		Complete(r)
}
