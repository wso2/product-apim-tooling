## apictl get apis

Display a list of APIs in an environment

### Synopsis

Display a list of APIs in the environment specified by the flag --environment, -e

```
apictl get apis [flags]
```

### Examples

```
apictl get apis -e dev
apictl get apis -e dev -q version:1.0.0
apictl get apis -e prod -q provider:admin -q version:1.0.0
apictl get apis -e prod -l 100
apictl get apis -e staging
NOTE: The flag (--environment (-e)) is mandatory
```

### Options

```
  -e, --environment string   Environment to be searched
      --format string        Pretty-print apis using Go Templates. Use "{{ jsonPretty . }}" to list all fields
  -h, --help                 help for apis
  -l, --limit string         Maximum number of apis to return (default "25")
  -q, --query strings        Query pattern
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl get](apictl_get.md)	 - Get APIs/APIProducts/Applications or revisions of a specific API/APIProduct in an environment or Get the log level of each API in an environment or Get the environments

