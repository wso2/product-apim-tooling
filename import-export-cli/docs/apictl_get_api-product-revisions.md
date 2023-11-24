## apictl get api-product-revisions

Display a list of Revisions for the API Products

### Synopsis

Display a list of Revisions available for the API Product in the environment specified

```
apictl get api-product-revisions [flags]
```

### Examples

```
apictl get api-product-revisions -n PizzaProduct -v 1.0.0 -e dev
apictl get api-product-revisions -n ShopProduct -v 1.0.0 -r admin -e dev
apictl get api-product-revisions -n PizzaProduct -q deployed:true -e dev
NOTE: All the 3 flags (--name (-n), --version (-v) and --environment (-e)) are mandatory.
```

### Options

```
  -e, --environment string   Environment to be searched
      --format string        Pretty-print revisions using Go Templates. Use "{{ jsonPretty . }}" to list all fields
  -h, --help                 help for api-product-revisions
  -n, --name string          Name of the API Product to get the revision
  -r, --provider string      Provider of the API Product
  -q, --query strings        Query pattern
  -v, --version string       Version of the API Product to get the revision
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl get](apictl_get.md)	 - Get APIs/APIProducts/Applications or revisions of a specific API/APIProduct in an environment or Get the Correlation Log Configurations or Get the log level of each API in an environment or Get the environments

