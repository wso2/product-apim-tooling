## apictl get policies api

Display a list of API Policies in an environment

### Synopsis

Display a list of API Policies in the environment specified by the flag --environment, -e

```
apictl get policies api [flags]
```

### Examples

```
apictl get policies api -e production
apictl get policies api -e prod --all
apictl get policies api -e prod -l 10
apictl get policies api -e prod --format jsonArray
NOTE: The flag (--environment (-e)) is mandatory
NOTE: Flags (--all) and (--limit (-l)) cannot be used at the same time
```

### Options

```
  -e, --environment string   Environment of the API Policies to be fetched
      --format string        Pretty-print API Policies using Go Templates. Use "jsonArray" to list all fields
  -h, --help                 help for API Policies
  -l, --limit string         Limit the number of policies fetched
  --all                      Fetch all available API Policies
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl get policies](apictl_get_policies.md)	 - Get Policy list

