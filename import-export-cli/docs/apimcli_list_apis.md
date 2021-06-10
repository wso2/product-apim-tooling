## apimcli list apis

Display a list of APIs in an environment

### Synopsis


Display a list of APIs in the environment specified by the flag --environment, -e

```
apimcli list apis [flags]
```

### Examples

```
apimcli apis list -e dev
apimcli apis list -e dev -q version:1.0.0
apimcli apis list -e prod -q provider:admin
apimcli list apis -e prod -l 100
apimcli apis list -e staging
```

### Options

```
  -e, --environment string   Environment to be searched
      --format string        Pretty-print apis using Go Templates. Use {{ jsonPretty . }} to list all fields
  -h, --help                 help for apis
  -l, --limit string         Maximum number of APIs to return
  -q, --query string         Query pattern
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO
* [apimcli list](apimcli_list.md)	 - List APIs/Applications in an environment or List the environments

