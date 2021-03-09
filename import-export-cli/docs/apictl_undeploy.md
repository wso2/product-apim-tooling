## apictl undeploy

Undeploy an API/API Product revision from a gateway environment

### Synopsis

Undeploy an API/API Product revision available in the environment specified by flag (--environment, -e) from the gateway specified by flag (--gateway, -g)

```
apictl undeploy [flags]
```

### Examples

```
apictl undeploy api -n TwitterAPI -v 1.0.0 -r admin --rev 1 -g Label1 Label2 -e dev
apictl undeploy api -n PizzaAPI -v 1.0.0 --rev 2 --all-gateways -e dev
apictl undeploy api -n LeasingAPIProduct -e dev
```

### Options

```
  -h, --help   help for undeploy
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl](apictl.md)	 - CLI for Importing and Exporting APIs and Applications and Managing WSO2 Micro Integrator
* [apictl undeploy api](apictl_undeploy_api.md)	 - Undeploy API
* [apictl undeploy api-product](apictl_undeploy_api-product.md)	 - Undeploy API Product

