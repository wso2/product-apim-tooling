## mi get data-services

Get information about data services deployed in a Micro Integrator

### Synopsis

Get information about the data services specified by command line argument [dataservice-name]
If not specified, list all the data services deployed in a Micro Integrator in the environment specified by the flag --environment, -e

```
mi get data-services [dataservice-name] [flags]
```

### Examples

```
To list all the data services
   mi get data-services -e dev
To get details about a specific data services
   mi get data-services SampleDataService -e dev
NOTE: The flag (--environment (-e)) is mandatory
```

### Options

```
  -e, --environment string   Environment to be searched
      --format string        Pretty-print using Go Templates. Use "{{ jsonPretty . }}" to list all fields
  -h, --help                 help for data-services
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [mi get](mi_get.md)	 - Get information about artifacts deployed in a Micro Integrator instance

