## apictl get correlation-logging

Display a list of correlation logging components in an environment

### Synopsis

Display a list of correlation logging components available in the environment specified
NOTE: The flag (--environment (-e)) is mandatory.

```
apictl get correlation-logging [flags]
```

### Examples

```
apictl get correlation-logging -e dev 
```

### Options

```
  -e, --environment string   Environment which the correlation logging components should be displayed
      --format string        Pretty-print correlation logging components using Go Templates. Use "{{ jsonPretty . }}" to list all fields
  -h, --help                 help for correlation-logging
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl get](apictl_get.md)	 - Get APIs/MCPServers/APIProducts/Applications or revisions of a specific API/MCPServers/APIProduct in an environment or Get the Correlation Log Configurations or Get the log level of each API/MCPServers in an environment or Get the environments

