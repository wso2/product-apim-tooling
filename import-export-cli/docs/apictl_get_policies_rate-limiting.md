## apictl get policies rate-limiting

Display a list of APIs in an environment

### Synopsis

Display a list of APIs in the environment specified by the flag --environment, -e

```
apictl get policies rate-limiting [flags]
```

### Examples

```
apictl get apictl get policies rate-limiting -e production -q type:sub rate-limiting -e dev
apictl get policies rate-limiting -e prod -q type:api
apictl get policies rate-limiting -e prod -q type:sub
apictl get policies rate-limiting -e staging -q type:global
NOTE: The flag (--environment (-e)) is mandatory
```

### Options

```
  -e, --environment string   Environment to be searched
      --format string        Pretty-print throttle policies using Go Templates. Use "{{ jsonPretty . }}" to list all fields
  -h, --help                 help for rate-limiting
  -q, --query strings        Query pattern
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl get policies](apictl_get_policies.md)	 - Get Policy list

