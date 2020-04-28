## apictl list products

Display a list of API Products in an environment

### Synopsis

Display a list of API Products in the environment specified by the flag --environment, -e

```
apictl list products [flags]
```

### Examples

```
apictl list products -e dev
apictl list products -e dev -q version:1.0.0
apictl list products -e prod -q provider:admin
apictl list products -e staging
NOTE: The flag (--environment (-e)) is mandatory
```

### Options

```
  -e, --environment string   Environment to be searched
      --format string        Pretty-print products using Go Templates. Use {{ jsonPretty . }} to list all fields
  -h, --help                 help for products
  -q, --query string         Query pattern
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl list](apictl_list.md)	 - List APIs/APIProducts/Applications in an environment or List the environments

