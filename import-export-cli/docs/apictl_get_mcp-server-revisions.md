## apictl get mcp-server-revisions

Display a list of Revisions for the MCP Server

### Synopsis

Display a list of Revisions available for the MCP Server in the environment specified

```
apictl get mcp-server-revisions [flags]
```

### Examples

```
apictl get mcp-server-revisions -n ChoreoConnect -v 1.0.0 -e dev
apictl get mcp-server-revisions -n ChoreoConnect -v 1.0.0 -r admin -e dev
apictl get mcp-server-revisions -n ChoreoConnect -v 1.0.0 -q deployed:true -e dev
NOTE: All the 3 flags (--name (-n), --version (-v) and --environment (-e)) are mandatory.
```

### Options

```
  -e, --environment string   Environment to be searched
      --format string        Pretty-print revisions using Go Templates. Use "{{ jsonPretty . }}" to list all fields
  -h, --help                 help for mcp-server-revisions
  -n, --name string          Name of the MCP Server to get the revision
  -r, --provider string      Provider of the MCP Server
  -q, --query strings        Query pattern
  -v, --version string       Version of the MCP Server to get the revision
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl get](apictl_get.md)	 - Get APIs/MCPServers/APIProducts/Applications or revisions of a specific API/MCPServers/APIProduct in an environment or Get the Correlation Log Configurations or Get the log level of each API/MCPServers in an environment or Get the environments

