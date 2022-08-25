## apictl get policies api

Display a list of API Policies

### Synopsis

Display a list of API Policies in the environment

```
apictl get policies api [flags]
```

### Examples

```
apictl get policies api -e dev
apictl get policies api -e dev --all
apictl get policies api -e dev -l 30
 NOTE: The flag (--environment (-e)) is mandatory
 NOTE: Flags (--all) and (--limit (-l)) cannot be used at the same time
```

### Options

```
      --all                  Get all API Policies
  -e, --environment string   Environment to be searched
      --format string        Pretty-print API Policies using Go Templates. Use "{{ jsonPretty . }}" to list all fields
  -h, --help                 help for api
  -l, --limit string         Maximum number of API Policies to return (default "25")
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl get policies](apictl_get_policies.md)	 - Get Policy list

