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
apictl k8s add/update api -n petstore --from-file=./Swagger.json --replicas=3 --namespace=wso2
```

### Options

```
  -f, --from-file stringArray   Path to swagger file
  -h, --help                    help for api
  -m, --mode string             Property to override the deploying mode. Available modes: privateJet, sidecar
  -n, --name string             Name of the API
      --namespace string        namespace of API
      --replicas int            replica set (default 1)
  -v, --version string          Property to override the existing docker image with same name and version
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl k8s update](apictl_k8s_update.md)	 - Update an API to the kubernetes cluster

