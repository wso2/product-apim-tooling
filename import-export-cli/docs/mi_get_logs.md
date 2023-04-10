## mi get logs

List all the available log files

### Synopsis

Download a log file by providing the file name and download location,
if not provided, list all the log files of the Micro Integrator in the environment specified by the flag --environment, -e

```
mi get logs [file-name] [flags]
```

### Examples

```
Example:
To list all the log files
   mi get logs -e dev
To download a selected log file
   mi get logs [file-name] -p [download-location] -e dev
NOTE: The flag (--environment (-e)) is mandatory
```

### Options

```
  -e, --environment string   Environment to be searched
      --format string        Pretty-print using Go Templates. Use "{{ jsonPretty . }}" to list all fields
  -h, --help                 help for logs
  -p, --path string          Path the file should be downloaded
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [mi get](mi_get.md)	 - Get information about artifacts deployed in a Micro Integrator instance

