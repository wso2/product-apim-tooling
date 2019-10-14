## apimctl list apis

Display a list of APIs in an environment

### Synopsis


Display a list of APIs in the environment specified by the flag --environment, -e

```
apimctl list apis [flags]
```

### Examples

```
apimctl apis list -e dev
apimctl apis list -e dev -q version:1.0.0
apimctl apis list -e prod -q provider:admin
apimctl apis list -e staging
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
* [apimctl list](apimctl_list.md)	 - List APIs/Applications in an environment or List the environments

