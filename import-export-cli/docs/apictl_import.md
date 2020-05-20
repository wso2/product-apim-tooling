## apictl import

Import an API Product to an environment

### Synopsis

Import an API Product to the environment specified by flag (--environment, -e)

```
apictl import [flags]
```

### Examples

```
apictl import api-product -f qa/LeasingAPIProduct.zip -e dev
```

### Options

```
  -h, --help   help for import
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl](apictl.md)	 - CLI for Importing and Exporting APIs and Applications
* [apictl import api-product](apictl_import_api-product.md)	 - Import API Product

