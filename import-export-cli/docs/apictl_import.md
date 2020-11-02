## apictl import

Import an API/API Product/Application to an environment

### Synopsis

Import an API to the environment specified by flag (--environment, -e)
Import an API Product to the environment specified by flag (--environment, -e)
Import an Application to the environment specified by flag (--environment, -e)

```
apictl import [flags]
```

### Examples

```
apictl import api -f qa/TwitterAPI.zip -e dev
apictl import api-product -f qa/LeasingAPIProduct.zip -e dev
apictl import app -f qa/apps/sampleApp.zip -e dev
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
* [apictl import api](apictl_import_api.md)	 - Import API
* [apictl import api-product](apictl_import_api-product.md)	 - Import API Product
* [apictl import app](apictl_import_app.md)	 - Import App

