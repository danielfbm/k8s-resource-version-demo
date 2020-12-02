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

package v1beta1

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var frobberlog = logf.Log.WithName("frobber-beta-resource")

func (r *Frobber) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-frobs-danielfbm-github-io-v1beta1-frobber,mutating=true,failurePolicy=fail,groups=frobs.danielfbm.github.io,resources=frobbers,verbs=create;update,versions=v1beta1,name=mfrobberv1beta1.kb.io

var _ webhook.Defaulter = &Frobber{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Frobber) Default() {
	frobberlog.Info("default", "name", r.Name)

	if r.Spec.Param == "" {
		r.Spec.Param = "another"
	}
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-frobs-danielfbm-github-io-v1beta1-frobber,mutating=false,failurePolicy=fail,groups=frobs.danielfbm.github.io,resources=frobbers,versions=v1beta1,name=vfrobberv1beta1.kb.io

var _ webhook.Validator = &Frobber{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Frobber) ValidateCreate() error {
	frobberlog.Info("validate create", "name", r.Name)

	return r.Validate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Frobber) ValidateUpdate(old runtime.Object) error {
	frobberlog.Info("validate update", "name", r.Name)

	return r.Validate()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Frobber) ValidateDelete() error {
	frobberlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
