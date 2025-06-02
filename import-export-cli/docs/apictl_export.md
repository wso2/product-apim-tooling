## apictl export

Export an API/API Product/Application/Policy in an environment

### Synopsis

Export an API available in the environment specified by flag (--environment, -e)
Export APIs available in the environment specified by flag (--environment, -e)
Export an API Product available in the environment specified by flag (--environment, -e)
Export an Application of a specific user (--owner, -o) in the environment specified by flag (--environment, -e)

```
apictl export [flags]
```

### Examples

```
apictl export api -n TwitterAPI -v 1.0.0 -r admin -e dev
apictl export apis -e dev
apictl export api-product -n LeasingAPIProduct -v 1.0.0 -e dev
apictl export app -n SampleApp -o admin -e dev
```

### Options

```
  -h, --help   help for export
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl](apictl.md)	 - CLI for Importing and Exporting APIs and Applications and Managing WSO2 Micro Integrator
* [apictl export api](apictl_export_api.md)	 - Export API
* [apictl export api-product](apictl_export_api-product.md)	 - Export API Product
* [apictl export apis](apictl_export_apis.md)	 - Export APIs for migration
* [apictl export app](apictl_export_app.md)	 - Export App
* [apictl export apps](apictl_export_apps.md)	 - Export Applications
* [apictl export policy](apictl_export_policy.md)	 - Export/Import a Policy

