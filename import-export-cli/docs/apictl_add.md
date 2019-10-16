## apictl add

Add an API to the kubernetes cluster

### Synopsis


Add an API from a Swagger file to the kubernetes cluster. JSON and YAML formats are accepted.
To execute kubernetes commands set mode to Kubernetes

### Examples

```
apictl add api -n petstore --from-file=./Swagger.json --replicas=1 --namespace=wso2

apictl add api -n petstore --from-file=./product-apim-tooling/import-export-cli/build/target/apictl/myapi --replicas=1 --namespace=wso2 --override=true
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
* [apictl](apictl.md)	 - CLI for Importing and Exporting APIs and Applications
* [apictl add api](apictl_add_api.md)	 - handle APIs in kubernetes cluster 

