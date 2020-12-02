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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// FrobberSpec defines the desired state of Frobber
type FrobberSpec struct {
	Height int    `json:"height"`
	Param  string `json:"param"`
}

// FrobberStatus defines the observed state of Frobber
type FrobberStatus struct {
	Phase string `json:"phase"`
}

// +kubebuilder:object:root=true

// Frobber is the Schema for the frobbers API
type Frobber struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FrobberSpec   `json:"spec,omitempty"`
	Status FrobberStatus `json:"status,omitempty"`
}

// Validate basic Frobber validation
func (r *Frobber) Validate() (err error) {
	var errs field.ErrorList
	defer func() {
		err = errs.ToAggregate()
	}()
	if r.Spec.Height <= 0 {
		errs = append(errs, field.Invalid(field.NewPath("spec", "height"), r.Spec.Height, `should have height greater than zero`))
	}
	return
}

// +kubebuilder:object:root=true

// FrobberList contains a list of Frobber
type FrobberList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Frobber `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Frobber{}, &FrobberList{})
}
