## apimctl update api

handle APIs in kubernetes cluster 

### Synopsis


Add, Update and Delete APIs in kubernetes cluster. JSON and YAML formats are accepted.
available modes are as follows
* kubernetes

```
apimctl update api [flags]
```

### Examples

```
apimctl add/update api -n petstore --from-file=./Swagger.json --replicas=3 --namespace=wso2
```

### Options

```
  -f, --from-file string   Path to swagger file
  -h, --help               help for api
  -n, --name string        Name of the API
      --namespace string   namespace of API
      --replicas int       replica set (default 1)
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO
* [apimctl update](apimctl_update.md)	 - Update an API to the kubernetes cluster

