## apictl mi get proxy-services

Get information about proxy services deployed in a Micro Integrator

### Synopsis

Get information about the proxy services specified by command line argument [proxy-name]
If not specified, list all the proxy services deployed in a Micro Integrator in the environment specified by the flag --environment, -e

```
apictl mi get proxy-services [proxy-name] [flags]
```

### Examples

```
To list all the proxy services
  apictl mi get proxy-services -e dev
To get details about a specific proxy services
  apictl mi get proxy-services SampleProxy -e dev
NOTE: The flag (--environment (-e)) is mandatory
```

### Options

```
  -e, --environment string   Environment to be searched
      --format string        Pretty-print using Go Templates. Use "{{ jsonPretty . }}" to list all fields
  -h, --help                 help for proxy-services
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl mi get](apictl_mi_get.md)	 - Get information about artifacts deployed in a Micro Integrator instance

