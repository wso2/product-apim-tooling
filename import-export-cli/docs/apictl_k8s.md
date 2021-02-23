## apictl k8s

Kubernetes mode based commands

### Synopsis

Kubernetes mode based commands such as add, update and delete API

```
apictl k8s [flags]
```

### Examples

```
apictl k8s add api -n petstore -f Swagger.json --namespace=wso2
apictl k8s update api -n petstore -f Swagger.json --namespace=wso2
apictl k8s delete api -n petstore
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

* [apictl](apictl.md)	 - CLI for Importing and Exporting APIs and Applications and Managing WSO2 Micro Integrator
* [apictl k8s add](apictl_k8s_add.md)	 - Add an API to the kubernetes cluster
* [apictl k8s delete](apictl_k8s_delete.md)	 - Delete resources related to kubernetes
* [apictl k8s gen](apictl_k8s_gen.md)	 - Generate deployment directory for K8S operator
* [apictl k8s update](apictl_k8s_update.md)	 - Update an API to the kubernetes cluster

