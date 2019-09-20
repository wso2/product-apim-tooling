## apimcli update

Update an API to the kubernetes cluster

### Synopsis


Update an existing API with  Swagger file in the kubernetes cluster. JSON and YAML formats are accepted.

### Examples

```
apimcli update api -n petstore --from-file=./Swagger.json --replicas=1 --namespace=wso2

apimcli update api -n petstore --from-file=./product-apim-tooling/import-export-cli/build/target/apimcli/myapi --replicas=1 --namespace=wso2
```

### Options

```
  -h, --help   help for update
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO
* [apimcli](apimcli.md)	 - CLI for Importing and Exporting APIs and Applications
* [apimcli update api](apimcli_update_api.md)	 - handle APIs in kubernetes cluster 

