## apictl mi get local-entries

Get information about local entries deployed in a Micro Integrator

### Synopsis

Get information about the local entries specified by command line argument [localentry-name]
If not specified, list all the local entries deployed in a Micro Integrator in the environment specified by the flag --environment, -e

```
apictl mi get local-entries [localentry-name] [flags]
```

### Examples

```
To list all the local entries
  apictl mi get local-entries -e dev
To get details about a specific local entries
  apictl mi get local-entries SampleLocalEntry -e dev
NOTE: The flag (--environment (-e)) is mandatory
```

### Options

```
  -e, --environment string   Environment to be searched
      --format string        Pretty-print using Go Templates. Use "{{ jsonPretty . }}" to list all fields
  -h, --help                 help for local-entries
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl mi get](apictl_mi_get.md)	 - Get information about artifacts deployed in a Micro Integrator instance

