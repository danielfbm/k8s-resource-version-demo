## 2. Upgrade to beta and stable

In [Step 1](STEP01.md) we did implement a simple example of a resource and its controller with integrated defaulting and validation. 

As time goes by new ideas and requirements are added which will naturally result in changes in the codebase and our resource definition.

In order to simulate a more complex scenario and keep complying with the [On compability](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api_changes.md#on-compatibility) section of the k8s community documentation, we will also first raise our api to beta, and further to stable.

### 2.1 Road to beta

As we are still working in our `v1alpha1`, the natural evolution is to raise it to `v1beta1`. For more details please check [Alpha, Beta, and Stable versions](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api_changes.md#alpha-beta-and-stable-versions)

In the compability section we can find a few important points that we have observe before starting:

> 1. Any API call (e.g. a structure POSTed to a REST endpoint) that succeeded before your change must succeed after your change.
> 2. Any API call that does not use your change must behave the same as it did before your change.
> 3. Any API call that uses your change must not cause problems (e.g. crash or degrade behavior) when issued against an API servers that do not include your change.
> 4. It must be possible to round-trip your change (convert to different API versions and back) with no loss of information.
> 5. Existing clients need not be aware of your change in order for them to continue to function as they did previously, even when your change is in use.
> 6. It must be possible to rollback to a previous version of API server that does not include your change and have no impact on API objects which do not use your change. API objects that use your change will be impacted in case of a rollback. 

In order to meet these, we need to add a new version `v1beta1` while keeping `v1alpha1` to keep old client compatibility. Also, we will explicitly set `v1alpha1` as our storage version. More details on why storage works this way please check [Operational overview](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api_changes.md#operational-overview)


### 2.2 Add beta

1. Create new version resource
*PS: Do not create controller just yet**


```bash
kubebuilder create api --group frobs --version v1beta1 --kind Frobber
```

2. Create webhooks

```bash
kubebuilder create webhook --group frobs --version v1beta1 --kind Frobber --defaulting --programmatic-validation
```

3. Copy `_types.go` and `_webhook.go` contents from `v1alpha1` to `v1beta1`


4. Change any `v1alpha1` text to `v1beta1` text, including and specially comments with `// +kubebuilder` prefix 

Warning on `// +kubebuilder` comments: Change the `name=` of each webhook comment to something different, i.e `name=mfrobberv1beta1.kb.io`

5. Add `// +kubebuilder:storageversion` to `Frobber` in `v1alpha1` 

```golang
// +kubebuilder:object:root=true
// +kubebuilder:storageversion

// Frobber is the Schema for the frobbers API
type Frobber struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FrobberSpec   `json:"spec,omitempty"`
	Status FrobberStatus `json:"status,omitempty"`
}
```

6. Run manifest and generate commands

```bash
make generate
make manifests
```

7. Build and deploy ([RUNNING](RUNNING.md))


### 2.3 Convertion

It is pretty common for APIs change during development, so naturally versions will have differences and need to be converted. A really good explanation for this is written in the [Hubs, spokes and other wheel metaphors](https://book.kubebuilder.io/multiversion-tutorial/conversion-concepts.html) document.

the `TL;DR` version is: One version as the center, all other version are converted from/to this version. For our usecase we make `v1alpha` the `Hub`, or central version, and convert `v1beta1` to it.


1. Create a `apis/v1alpha1/frobber_conversion.go` and declare this version as `Hub`

```golang
package v1alpha1

func (*Frobber) Hub() {}

```

2. Create a `apis/v1beta1/frobber_conversion.go`

```golang
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

```

3. Setup webhooks

Because we already setup webhooks in the step1, nothing needs to be added. The [kubebuilder book](https://book.kubebuilder.io/multiversion-tutorial/webhooks.html#webhook-setup) has the following statement:

> This setup doubles as setup for our conversion webhooks: as long as our
> types implement the [Hub](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/conversion?tab=doc#Hub) and [Convertible](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/conversion?tab=doc#Convertible) interfaces, a conversion webhook will be registered.

As we are not changing types here we just need to run

4. Try the alpha and beta apis

*PS: change accordingly to your domain/api group*

```bash
kubectl get frobbers.v1beta1.frobs.danielfbm.github.io

kubectl get frobbers.v1alpha1.frobs.danielfbm.github.io
```

### 2.4 Differences

For most cases having exactly the same APIs with different versions is not useful. The motivation to create new versions is generally:

- Slight changes in the API
- Change of behaviour
- Deprecation cycle of apis/features

## Next

| Previous | Next |
|----------|------|
| [Step 1](STEP01.md) | [Step 2](STEP02.md) |