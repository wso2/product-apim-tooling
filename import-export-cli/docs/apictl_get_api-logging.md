## apictl get api-logging

Display a list of API loggers in an environment

### Synopsis

Display a list of API loggers available for the APIs in the environment specified

```
apictl get api-logging [flags]
```

### Examples

```
apictl get api-logging -e dev --tenant-domain carbon.super
apictl get api-logging --api-id bf36ca3a-0332-49ba-abce-e9992228ae06 -e dev --tenant-domain carbon.super
```

### Options

```
  -i, --api-id string          API ID
  -e, --environment string     Environment of the APIs which the API loggers should be displayed
      --format string          Pretty-print API loggers using Go Templates. Use "{{ jsonPretty . }}" to list all fields
  -h, --help                   help for api-logging
      --tenant-domain string   Tenant Domain
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl get](apictl_get.md)	 - Get APIs/MCPServers/APIProducts/Applications or revisions of a specific API/MCPServers/APIProduct in an environment or Get the Correlation Log Configurations or Get the log level of each API/MCPServers in an environment or Get the environments

