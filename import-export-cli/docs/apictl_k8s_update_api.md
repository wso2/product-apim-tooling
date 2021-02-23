## apictl k8s update api

Handle APIs in kubernetes cluster 

### Synopsis

Add, Update and Delete APIs in kubernetes cluster. JSON and YAML formats are accepted.
available modes are as follows
* kubernetes

```
apictl k8s update api [flags]
```

### Examples

```
apictl k8s add/update api -n petstore -f Swagger.json --namespace=wso2
```

### Options

```
  -f, --file string        Path to swagger, zip file or API project
  -h, --help               help for api
  -n, --name string        Name of the API
      --namespace string   namespace of API
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl k8s update](apictl_k8s_update.md)	 - Update an API to the kubernetes cluster

