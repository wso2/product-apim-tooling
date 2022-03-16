## apictl export api-product

Export API Product

### Synopsis

Export an API Product in an environment

```
apictl export api-product (--name <name-of-the-api-product> --provider <provider-of-the-api-product> --environment <environment-from-which-the-api-product-should-be-exported>) [flags]
```

### Examples

```
apictl export api-product -n LeasingAPIProduct -e dev
apictl export api-product -n CreditAPIProduct -r admin -e production
NOTE: Both the flags (--name (-n) and --environment (-e)) are mandatory
```

### Options

```
  -e, --environment string   Environment to which the API Product should be exported
      --format string        File format of exported archive (json or yaml) (default "YAML")
  -h, --help                 help for api-product
      --latest               Export the latest revision of the API Product
  -n, --name string          Name of the API Product to be exported
      --preserve-status      Preserve API Product status when exporting. Otherwise API Product will be exported in CREATED status (default true)
  -r, --provider string      Provider of the API Product
      --rev string           Revision number of the API Product to be exported
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl export](apictl_export.md)	 - Export an API/API Product/Application in an environment

