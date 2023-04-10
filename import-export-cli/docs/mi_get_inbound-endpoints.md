## mi get inbound-endpoints

Get information about inbound endpoints deployed in a Micro Integrator

### Synopsis

Get information about the inbound endpoints specified by command line argument [inbound-name]
If not specified, list all the inbound endpoints deployed in a Micro Integrator in the environment specified by the flag --environment, -e

```
mi get inbound-endpoints [inbound-name] [flags]
```

### Examples

```
To list all the inbound endpoints
   mi get inbound-endpoints -e dev
To get details about a specific inbound endpoints
   mi get inbound-endpoints SampleInboundEndpoint -e dev
NOTE: The flag (--environment (-e)) is mandatory
```

### Options

```
  -e, --environment string   Environment to be searched
      --format string        Pretty-print using Go Templates. Use "{{ jsonPretty . }}" to list all fields
  -h, --help                 help for inbound-endpoints
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [mi get](mi_get.md)	 - Get information about artifacts deployed in a Micro Integrator instance

