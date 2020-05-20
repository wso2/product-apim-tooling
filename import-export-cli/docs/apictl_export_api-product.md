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
apictl export api-product -n CreditAPIProduct -v 1.0.0 -r admin -e production
NOTE: Both the flags (--name (-n) and --environment (-e)) are mandatory
```

### Options

```
  -e, --environment string   Environment to which the API Product should be exported
      --format string        File format of exported archive (json or yaml)
  -h, --help                 help for api-product
  -n, --name string          Name of the API Product to be exported
  -r, --provider string      Provider of the API Product
  -v, --version string       Version of the API Product to be exported
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl export](apictl_export.md)	 - Export an API Product in an environment

