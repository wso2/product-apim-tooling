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
  -a, --apiEndPoint string      
  -e, --env stringArray         Environment variables to be passed to deployment
  -f, --from-file stringArray   Path to swagger file
  -h, --help                    help for api
      --hostname string         Ingress hostname that the API is being exposed
  -i, --image string            Image of the API. If specified, ignores the value of --override
  -m, --mode string             Property to override the deploying mode. Available modes: privateJet, sidecar
  -n, --name string             Name of the API
      --namespace string        namespace of API
      --override                Property to override the existing docker image with the given name and version
      --replicas int            replica set (default 1)
  -v, --version string          Property to override the API version
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl add](apictl_add.md)	 - Add an API to the kubernetes cluster

