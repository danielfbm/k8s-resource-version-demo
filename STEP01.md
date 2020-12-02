## 1. Bootstrap the project

*Tag:* `step-1`

### 1.1 Basic structure

Outside of `GOPATH`

```bash
go mod init mod.example.com
kubebuilder init --domain example.com
```

### 1.2 Adding the first version of Frobber

Create your first API object

```bash
kubebuilder create api --group frobs --version v1alpha1 --kind Frobber
```

Following the [**On compability**](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api_changes.md#on-compatibility) part we should add a few things to our type

on `api/v1alpha1/frobber_types.go`

```golang
// FrobberSpec defines the desired state of Frobber
type FrobberSpec struct {
	Height int    `json:"height"`
    Param  string `json:"param"`
}

// FrobberStatus defines the observed state of Frobber
type FrobberStatus struct {
	Phase string `json:"phase"`
}
```

Update manifests and any generated code

```bash
make generate
make manifests
```

Now we can run the controller and check it out. For  specifics please check [RUNNING](RUNNING.md)


### 1.3 Adding webhooks

To work with controllers only it is fundamental to add mutation and validation webhooks to implement `defaults` and `validation`. This is done by using [Admission controllers kubernetes feature](https://kubernetes.io/blog/2019/03/21/a-guide-to-kubernetes-admission-controllers/).

#### 1.3.1 Scafold webhooks

To scafold all the necessary code for defaulting and validation just run the code below. For the details check [this section](https://book.kubebuilder.io/cronjob-tutorial/webhook-implementation.html) of the kubebuilder book.

```bash
kubebuilder create webhook --group frobs --version v1alpha1 --kind Frobber --defaulting --programmatic-validation
```

#### 1.3.2 Install certmanager

Kubernetes requires a https server to call its webhooks, currently our deployments does not have any, so we will invoke cert-manager to manage certificates and inject them into our manager.


Install a certmanager in your cluster: [Official](https://cert-manager.io/docs/installation/)

or 

```bash
kubectl apply -f config/cert-manager.v0.16.0.yaml
```

#### 1.3.2 Enable webhook and certmanager

inside `config/default/kustomization.yaml` enable all the `[WEBHOOK]` and `[CERTMANAGER]` sections. This should result in something like:

```yaml

bases:
- ../crd
- ../rbac
- ../manager
# [WEBHOOK] To enable webhook, uncomment all the sections with [WEBHOOK] prefix including the one in 
# crd/kustomization.yaml
- ../webhook
# [CERTMANAGER] To enable cert-manager, uncomment all sections with 'CERTMANAGER'. 'WEBHOOK' components are required.
- ../certmanager
# [PROMETHEUS] To enable prometheus monitor, uncomment all sections with 'PROMETHEUS'. 
#- ../prometheus

patchesStrategicMerge:
  # Protect the /metrics endpoint by putting it behind auth.
  # If you want your controller-manager to expose the /metrics
  # endpoint w/o any authn/z, please comment the following line.
- manager_auth_proxy_patch.yaml

# [WEBHOOK] To enable webhook, uncomment all the sections with [WEBHOOK] prefix including the one in 
# crd/kustomization.yaml
- manager_webhook_patch.yaml

# [CERTMANAGER] To enable cert-manager, uncomment all sections with 'CERTMANAGER'.
# Uncomment 'CERTMANAGER' sections in crd/kustomization.yaml to enable the CA injection in the admission webhooks.
# 'CERTMANAGER' needs to be enabled to use ca injection
- webhookcainjection_patch.yaml

# the following config is for teaching kustomize how to do var substitution
vars:
# [CERTMANAGER] To enable cert-manager, uncomment all sections with 'CERTMANAGER' prefix.
- name: CERTIFICATE_NAMESPACE # namespace of the certificate CR
  objref:
    kind: Certificate
    group: cert-manager.io
    version: v1alpha2
    name: serving-cert # this name should match the one in certificate.yaml
  fieldref:
    fieldpath: metadata.namespace
- name: CERTIFICATE_NAME
  objref:
    kind: Certificate
    group: cert-manager.io
    version: v1alpha2
    name: serving-cert # this name should match the one in certificate.yaml
- name: SERVICE_NAMESPACE # namespace of the service
  objref:
    kind: Service
    version: v1
    name: webhook-service
  fieldref:
    fieldpath: metadata.namespace
- name: SERVICE_NAME
  objref:
    kind: Service
    version: v1
    name: webhook-service
```

#### 1.3.3 Implement default and validation

Although not best practice, for simplicity sake lets change `api/v1alpha1/frobber_types.go` adding a `Validate` method that checks that `Height` should be greater than zero

```golang
// Validate basic Frobber validation
func (r *Frobber) Validate() (err error) {
	var errs field.ErrorList
	defer func() {
		err = errs.ToAggregate()
	}()
	if r.Spec.Height <= 0 {
		errs = append(errs, field.Invalid(field.NewPath("spec", "height"), r.Spec.Height, `should have height greater than zero`)
	}
	return
}
```

In our `api/v1alpha1/frobber_webhook.go` file implement the following change, which adds default value for `spec.param` and invoke the `Validate` method during `Create` and `Update`

```golang
[...]

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Frobber) Default() {
	frobberlog.Info("default", "name", r.Name)

	if r.Spec.Param == "" {
		r.Spec.Param = "default"
	}
}

[...]


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

[...]
```



### 1.4 Adding data

Modify your `config/samples/frobs_v1alpha1_frobber.yaml` to reflect the structure

```yaml
apiVersion: frobs.danielfbm.github.io/v1alpha1
kind: Frobber
metadata:
  name: frobber-sample
spec:
  # Add fields here
  height: 3
  param: ""
```

apply 

```bash
kubectl apply -f config/samples/frobs_v1alpha1_frobber.yaml
```

check that it is all well, with our "default" for `param`


```bash
kubectl get frobber frobber-sample -o yaml
```

Implement a basic controller that prints `height` and `param` changing `controllers/frobber_controller.go`'s `Reconcile` function:

```golang
[...]

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

[...]
```

create a second sample as `config/samples/frobs_v1alpha1_frobber_invalid.yaml`

```yaml
apiVersion: frobs.danielfbm.github.io/v1alpha1
kind: Frobber
metadata:
  name: frobber-invalid
spec:
  height: 0
  param: "inches"
```

apply and check if the validation is working as expected

```bash
kubectl apply -f config/samples/frobs_v1alpha1_frobber_invalid.yaml

$ Error from server (spec.height: Invalid value: 0: should have height greater than zero): error when creating "config/samples/frobs_v1alpha1_frobber_invalid.yaml": admission webhook "vfrobber.kb.io" denied the request: spec.height: Invalid value: 0: should have height greater than zero
```

Now we are finally ready to start


## Next

| Previous | Next |
|----------|------|
| [README](README.md) | [Step 2](STEP02.md) |