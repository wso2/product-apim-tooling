## apictl k8s

Kubernetes mode based commands

### Synopsis

Kubernetes mode based commands such as install, uninstall, add/update api, change registry.

```
apictl k8s [flags]
```

### Examples

```
apictl k8s install api-operator
apictl k8s uninstall api-operator
apictl k8s add api -n petstore --from-file=./Swagger.json --replicas=1 --namespace=wso2
apictl k8s update api -n petstore --from-file=./Swagger.json --replicas=1 --namespace=wso2
apictl k8s change registry
```

### Options

```
  -h, --help   help for k8s
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl](apictl.md)	 - CLI for Importing and Exporting APIs and Applications
* [apictl k8s add](apictl_k8s_add.md)	 - Add an API to the kubernetes cluster
* [apictl k8s change](apictl_k8s_change.md)	 - Change a configuration in K8s cluster resource
* [apictl k8s delete](apictl_k8s_delete.md)	 - Delete resources related to kubernetes
* [apictl k8s install](apictl_k8s_install.md)	 - Install an operator in the configured K8s cluster
* [apictl k8s uninstall](apictl_k8s_uninstall.md)	 - Uninstall an operator in the configured K8s cluster
* [apictl k8s update](apictl_k8s_update.md)	 - Update an API to the kubernetes cluster

