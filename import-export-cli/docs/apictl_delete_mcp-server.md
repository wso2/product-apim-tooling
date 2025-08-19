## apictl delete mcp-server

Delete MCP Server

### Synopsis

Delete an MCP Server from an environment

```
apictl delete mcp-server (--name <name-of-the-mcp-server> --version <version-of-the-mcp-server> --provider <provider-of-the-mcp-server> --environment <environment-from-which-the-mcp-server-should-be-deleted>) [flags]
```

### Examples

```
apictl delete mcp-server -n ChoreoConnect -v 1.0.0 -r admin -e dev
apictl delete mcp-server -n ChoreoConnect -v 2.1.0 -e production
NOTE: The 3 flags (--name (-n), --version (-v), and --environment (-e)) are mandatory.
```

### Options

```
  -e, --environment string   Environment from which the MCP Server should be deleted
  -h, --help                 help for mcp-server
  -n, --name string          Name of the MCP Server to be deleted
  -r, --provider string      Provider of the MCP Server to be deleted
  -v, --version string       Version of the MCP Server to be deleted
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl delete](apictl_delete.md)	 - Delete an API/MCPServer/APIProduct/Application in an environment

