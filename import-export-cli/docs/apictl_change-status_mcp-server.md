## apictl change-status mcp-server

Change Status of an MCP Server

### Synopsis

Change the lifecycle status of an MCP Server in an environment

```
apictl change-status mcp-server (--action <action-of-the-mcpserver-state-change> --name <name-of-the-mcpserver> --version <version-of-the-mcpserver> --provider <provider-of-the-mcpserver> --environment <environment-from-which-the-mcpserver-state-should-be-changed>) [flags]
```

### Examples

```
apictl change-status mcp-server -a Publish -n MyMCPServer -v 1.0.0 -r admin -e dev
apictl change-status mcp-server -a Publish -n MyMCPServer -v 2.1.0 -e production
NOTE: The 4 flags (--action (-a), --name (-n), --version (-v), and --environment (-e)) are mandatory.
```

### Options

```
  -a, --action string        Action to be taken to change the status of the MCP Server
  -e, --environment string   Environment of which the MCP Server state should be changed
  -h, --help                 help for mcp-server
  -n, --name string          Name of the MCP Server to be state changed
  -r, --provider string      Provider of the MCP Server
  -v, --version string       Version of the MCP Server to be state changed
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl change-status](apictl_change-status.md)	 - Change Status of an API, MCP Server or Product

