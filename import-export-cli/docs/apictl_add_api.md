## apictl add api

handle APIs in kubernetes cluster 

### Synopsis

Add, Update and Delete APIs in kubernetes cluster. JSON and YAML formats are accepted.
available modes are as follows
* kubernetes

```
apictl add api [flags]
```

### Examples

```
apictl add/update api -n petstore --from-file=./Swagger.json --replicas=3 --namespace=wso2
```

### Options

```
  -f, --from-file stringArray   Path to swagger file
  -h, --help                    help for api
  -m, --mode string             Property to override the deploying mode. Available modes: privateJet, sidecar
  -n, --name string             Name of the API
      --namespace string        namespace of API
      --override                Property to override the existing docker image with same name and version (default false)
      --replicas int            replica set (default 1)
  -v, --version string          Version of the API
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl add](apictl_add.md)	 - Add an API to the kubernetes cluster

