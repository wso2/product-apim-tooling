## apictl list api-products

Display a list of API Products in an environment

### Synopsis

Display a list of API Products in the environment specified by the flag --environment, -e

```
apictl list api-products [flags]
```

### Examples

```
apictl list api-products -e dev
apictl list api-products -e dev -q version:1.0.0
apictl list api-products -e prod -q provider:admin context:/myproduct
apictl list apis -e prod -l 25
apictl list api-products -e staging
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

* [apictl list](apictl_list.md)	 - List APIs/APIProducts/Applications in an environment or List the environments

