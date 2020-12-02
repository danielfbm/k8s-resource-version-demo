# Running options for kubebuilder controllers

In your `Makefile` do the follow changes:

1. Add `PKG` and `KIND` variables to the top

```Makefile
# Image URL to use all building/pushing image targets
IMG ?= controller:latest
PKG ?= ko://danielfbm.github.io/k8s-resource-version
KIND ?= 
```

2. Add `KIND` support to the `deploy` command

```Makefile
# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests
	cd config/manager && kustomize edit set image controller=${IMG}
ifneq (,${KIND})
	kind load docker-image --name ${KIND} ${IMG}
endif 
	kustomize build config/default | kubectl apply -f -
```

This will give support to running with [`kind`](https://sigs.k8s.io/kind)


3. Change your manager yaml on `config/manager/manager.yaml` specifying an `imagePullPolicy`

```yaml
[...]
      containers:
      - command:
        - /manager
        args:
        - --enable-leader-election
        image: controller:latest
        imagePullPolicy: IfNotPresent
        name: manager
[...]
```


## Docker for Mac 

```bash
make docker-build
make deploy
```

## Kind

Set `KIND` variable with your `kind` cluster name

```bash
make docker-build
make KIND=kind deploy
```


