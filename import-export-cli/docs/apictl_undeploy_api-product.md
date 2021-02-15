## apictl undeploy api-product

Undeploy API Product

### Synopsis

Undeploy an API Product revision from the given gateway environment

```
apictl undeploy api-product (--name <name-of-the-api-product> --version <version-of-the-api-product> --rev<revision-number-of-the-api-product> --gateway <gateway-environment> --environment <environment-from-which-the-api-product-should-be-undeployed>) [flags]
```

### Examples

```
apictl undeploy api-product -n TwitterAPIProduct -v 1.0.0 -r admin -g Label1 -e dev
apictl undeploy api-product -n StoreProduct -v 2.1.0 --rev 6 --all-gateways -e production
apictl undeploy api-product -n FacebookProduct -v 2.1.0 --rev 2 -g Label1 -e production
NOTE: All the 4 flags (--name (-n), --version (-v), --rev, --environment (-e)) and one from (--gateway (-g) or --all-gateways) are mandatory.
```

### Options

```
      --all-gateways         Undeploy the revision from all the deployed gateways at once
  -e, --environment string   Environment of which the API Product should be undeployed
  -g, --gateway string       Gateway which the revision has to be undeployed
  -h, --help                 help for api-product
  -n, --name string          Name of the API Product to be exported
  -r, --provider string      Provider of the API
      --rev string           Revision number of the API Product to undeploy
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl undeploy](apictl_undeploy.md)	 - Undeploy an API/API Product revision from a gateway environment

