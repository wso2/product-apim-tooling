## apictl mi get apis

Get information about apis deployed in a Micro Integrator

### Synopsis

Get information about the apis specified by command line argument [api-name]
If not specified, list all the apis deployed in a Micro Integrator in the environment specified by the flag --environment, -e

```
apictl mi get apis [api-name] [flags]
```

### Examples

```
To list all the apis
  apictl mi get apis -e dev
To get details about a specific apis
  apictl mi get apis SampleIntegrationAPI -e dev
NOTE: The flag (--environment (-e)) is mandatory
```

### Options

```
  -e, --environment string   Environment to be searched
      --format string        Pretty-print using Go Templates. Use "{{ jsonPretty . }}" to list all fields
  -h, --help                 help for apis
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl mi get](apictl_mi_get.md)	 - Get information about artifacts deployed in a Micro Integrator instance

