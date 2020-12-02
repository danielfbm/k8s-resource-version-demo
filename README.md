
Following [API conventions from K8s sig-architecture](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api_changes.md) how to do it using kubebuilder?

You can `git checkout` specific **tags** to see how the code should be at any specific step

## 1. Bootstrap the project

Outside of `GOPATH`

```bash
go mod init mod.example.com
kubebuilder init --domain example.com
```

## 2. Adding the first version of Frobber

*Tag:* `step-2`

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


