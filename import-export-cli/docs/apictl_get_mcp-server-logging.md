## apictl get mcp-server-logging

Display a list of MCP Server loggers in an environment

### Synopsis

Display a list of MCP Server loggers available for the MCP Servers in the environment specified

```
apictl get mcp-server-logging [flags]
```

### Examples

```
apictl get mcp-server-logging -e dev --tenant-domain carbon.super
apictl get mcp-server-logging --mcp-server-id bf36ca3a-0332-49ba-abce-e9992228ae06 -e dev --tenant-domain carbon.super
```

### Options

```
  -e, --environment string     Environment of the MCP Servers which the MCP Server loggers should be displayed
      --format string          Pretty-print MCP Server loggers using Go Templates. Use "{{ jsonPretty . }}" to list all fields
  -h, --help                   help for mcp-server-logging
  -i, --mcp-server-id string   MCP Server ID
      --tenant-domain string   Tenant Domain
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl get](apictl_get.md)	 - Get APIs/MCPServers/APIProducts/Applications or revisions of a specific API/MCPServers/APIProduct in an environment or Get the Correlation Log Configurations or Get the log level of each API/MCPServers in an environment or Get the environments

