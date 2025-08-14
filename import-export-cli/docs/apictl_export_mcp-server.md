## apictl export mcp-server

Export MCP Server

### Synopsis

Export an MCP Server from an environment

```
apictl export mcp-server (--name <name-of-the-mcp-server> --version <version-of-the-mcp-server> --provider <provider-of-the-mcp-server> --environment <environment-from-which-the-mcp-server-should-be-exported>) [flags]
```

### Examples

```
apictl export mcp-server -n ChoreoConnect -v 1.0.0 -r admin -e dev
apictl export mcp-server -n ChoreoConnect -v 2.1.0 --rev 6 -r admin -e production
apictl export mcp-server -n ChoreoConnect -v 2.1.0 --rev 2 -r admin -e production
NOTE: All the 3 flags (--name (-n), --version (-v) and --environment (-e)) are mandatory. If --rev is not provided, working copy of the MCP Server
without deployment environments will be exported.
```

### Options

```
  -e, --environment string     Environment to which the MCP Server should be exported
      --format string          File format of exported archive(json or yaml) (default "YAML")
  -h, --help                   help for mcp-server
      --latest                 Export the latest revision of the MCP Server
  -n, --name string            Name of the MCP Server to be exported
      --preserve-credentials   Preserve endpoint credentials when exporting. Otherwise credentials will not be exported
      --preserve-status        Preserve MCP Server status when exporting. Otherwise MCP Server will be exported in CREATED status (default true)
  -r, --provider string        Provider of the MCP Server
      --rev string             Revision number of the MCP Server to be exported
  -v, --version string         Version of the MCP Server to be exported
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl export](apictl_export.md)	 - Export an API/MCPServer/API Product/Application/Policy in an environment

