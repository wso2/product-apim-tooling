## apictl mi get endpoints

Get information about endpoints deployed in a Micro Integrator

### Synopsis

Get information about the endpoints specified by command line argument [endpoint-name]
If not specified, list all the endpoints deployed in a Micro Integrator in the environment specified by the flag --environment, -e

```
apictl mi get endpoints [endpoint-name] [flags]
```

### Examples

```
To list all the endpoints
  apictl mi get endpoints -e dev
To get details about a specific endpoints
  apictl mi get endpoints SampleEndpoint -e dev
NOTE: The flag (--environment (-e)) is mandatory
```

### Options

```
  -e, --environment string   Environment to be searched
      --format string        Pretty-print endpoints using Go Templates. Use "{{ jsonPretty . }}" to list all fields
  -h, --help                 help for endpoints
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl mi get](apictl_mi_get.md)	 - Get information about artifacts deployed in a Micro Integrator instance

