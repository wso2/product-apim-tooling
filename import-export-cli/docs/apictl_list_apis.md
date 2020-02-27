## apictl list apis

Display a list of APIs in an environment

### Synopsis

Display a list of APIs in the environment specified by the flag --environment, -e

```
apictl list apis [flags]
```

### Examples

```
apictl apis list -e dev
apictl apis list -e dev -q version:1.0.0
apictl apis list -e prod -q provider:admin
apictl apis list -e staging
```

### Options

```
  -e, --environment string   Environment to be searched
      --format string        Pretty-print apis using Go Templates. Use {{ jsonPretty . }} to list all fields
  -h, --help                 help for apis
  -q, --query string         Query pattern
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl list](apictl_list.md)	 - List APIs/Applications in an environment or List the environments

