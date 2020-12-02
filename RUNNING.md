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

4. Force deployment

If as-is the tag will always be the same and will not trigger any deployment for new versions. To `*force*` a redeploy you can execute `kubectl -n k8s-resource-version-demo-system delete pods --all` **AFTER** the deploy commands below

Add the `kill` command to your `Makefile`

```Makefile
kill:
	kubectl -n k8s-resource-version-demo-system delete pods --all
```

> This can be avoided if a `IMG` variable with a full repository:tag is provided. i.e `make docker-build IMG=danielfbm/controller:v1`


## Docker for Mac 

```bash
make docker-build
make deploy
make kill
```

## Kind

Set `KIND` variable with your `kind` cluster name

```bash
make docker-build
make KIND=kind deploy
make kill
```


