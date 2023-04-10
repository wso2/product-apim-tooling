## mi get log-levels

Get information about a Logger configured in a Micro Integrator

### Synopsis

Get information about the Logger specified by command line argument [logger-name]
configured in a Micro Integrator in the environment specified by the flag --environment, -e

```
mi get log-levels [logger-name] [flags]
```

### Examples

```
To get details about a specific logger
   mi get log-levels org-apache-coyote -e dev
NOTE: The flag (--environment (-e)) is mandatory
```

### Options

```
  -e, --environment string   Environment to be searched
      --format string        Pretty-print using Go Templates. Use "{{ jsonPretty . }}" to list all fields
  -h, --help                 help for log-levels
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [mi get](mi_get.md)	 - Get information about artifacts deployed in a Micro Integrator instance

