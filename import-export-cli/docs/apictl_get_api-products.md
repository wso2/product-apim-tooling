## apictl get api-products

Display a list of API Products in an environment

### Synopsis

Display a list of API Products in the environment specified by the flag --environment, -e

```
apictl get api-products [flags]
```

### Examples

```
apictl get api-products -e dev
apictl get api-products -e dev -q provider:devops
apictl get api-products -e prod -q provider:admin context:/myproduct
apictl get api-products -e prod -l 25
apictl get api-products -e staging
NOTE: The flag (--environment (-e)) is mandatory
```

### Options

```
  -e, --environment string   Environment to be searched
      --format string        Pretty-print API Products using Go Templates. Use "{{ jsonPretty . }}" to list all fields
  -h, --help                 help for api-products
  -l, --limit string         Maximum number of API Products to return (default "25")
  -q, --query string         Query pattern
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl get](apictl_get.md)	 - Get APIs/APIProducts/Applications in an environment or Get the environments

