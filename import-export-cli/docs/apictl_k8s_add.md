## apictl k8s add

Add an API to the kubernetes cluster

### Synopsis

Add an API from a Swagger file to the kubernetes cluster. JSON and YAML formats are accepted.
To execute kubernetes commands set mode to Kubernetes

### Examples

```
apictl k8s add api -n petstore --from-file=./Swagger.json --replicas=1 --namespace=wso2

apictl k8s add api -n petstore --from-file=./product-apim-tooling/import-export-cli/build/target/apictl/myapi --replicas=1 --namespace=wso2 --override=true
```

### Options

```
  -h, --help   help for add
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl k8s](apictl_k8s.md)	 - Kubernetes mode based commands
* [apictl k8s add api](apictl_k8s_add_api.md)	 - Handle APIs in kubernetes cluster 

