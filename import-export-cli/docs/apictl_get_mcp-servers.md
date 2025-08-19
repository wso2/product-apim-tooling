## apictl get mcp-servers

Display a list of MCP Servers in an environment

### Synopsis

Display a list of MCP Servers in the environment specified by the flag --environment, -e

```
apictl get mcp-servers [flags]
```

### Examples

```
apictl get mcp-servers -e dev
apictl get mcp-servers -e dev -q version:1.0.0
apictl get mcp-servers -e prod -q provider:admin -q version:1.0.0
apictl get mcp-servers -e prod -l 100
apictl get mcp-servers -e staging
NOTE: The flag (--environment (-e)) is mandatory
```

### Options

```
  -e, --environment string   Environment to be searched
      --format string        Pretty-print mcp-servers using Go Templates. Use "{{ jsonPretty . }}" to list all fields
  -h, --help                 help for mcp-servers
  -l, --limit string         Maximum number of MCP servers to return (default "25")
  -q, --query strings        Query pattern
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl get](apictl_get.md)	 - Get APIs/MCPServers/APIProducts/Applications or revisions of a specific API/MCPServers/APIProduct in an environment or Get the Correlation Log Configurations or Get the log level of each API/MCPServers in an environment or Get the environments

