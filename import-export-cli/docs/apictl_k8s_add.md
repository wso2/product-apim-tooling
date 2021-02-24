## apictl k8s add

Add an API to the kubernetes cluster

### Synopsis

Add an API either from a Swagger file, project zip for API project to the kubernetes cluster. 
JSON, YAML, zip and API project formats are accepted.

### Examples

```
apictl k8s add api -n petstore -f Swagger.json --namespace=wso2
apictl k8s add api -n petstore -f product-apim-tooling/import-export-cli/build/target/apictl/myapi.zip --namespace=wso2
apictl k8s add api -n petstore -f myapi --namespace=wso2
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

