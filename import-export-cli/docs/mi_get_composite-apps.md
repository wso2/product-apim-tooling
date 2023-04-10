## mi get composite-apps

Get information about composite apps deployed in a Micro Integrator

### Synopsis

Get information about the composite apps specified by command line argument [app-name]
If not specified, list all the composite apps deployed in a Micro Integrator in the environment specified by the flag --environment, -e

```
mi get composite-apps [app-name] [flags]
```

### Examples

```
To list all the composite apps
   mi get composite-apps -e dev
To get details about a specific composite apps
   mi get composite-apps SampleApp -e dev
NOTE: The flag (--environment (-e)) is mandatory
```

### Options

```
  -e, --environment string   Environment to be searched
      --format string        Pretty-print using Go Templates. Use "{{ jsonPretty . }}" to list all fields
  -h, --help                 help for composite-apps
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [mi get](mi_get.md)	 - Get information about artifacts deployed in a Micro Integrator instance

