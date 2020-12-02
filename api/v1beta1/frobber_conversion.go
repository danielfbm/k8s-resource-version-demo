package v1beta1

import (
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts to the Hub version (v1alpha1).
func (r *Frobber) ConvertTo(dstRaw conversion.Hub) error {
	return nil
}

// ConvertFrom converts from the Hub version (v1alpha1) to this version.
func (r *Frobber) ConvertFrom(srcRaw conversion.Hub) error {
	return nil
}
