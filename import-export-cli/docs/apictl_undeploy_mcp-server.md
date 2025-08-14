## apictl undeploy mcp-server

Undeploy MCP Server

### Synopsis

Undeploy an MCP Server revision from gateway environments

```
apictl undeploy mcp-server (--name <name-of-the-mcpserver> --version <version-of-the-mcpserver> --provider <provider-of-the-mcpserver> --rev <revision-number-of-the-mcpserver> --gateway-env <gateway-environment> --environment <environment-from-which-the-mcpserver-should-be-undeployed>) [flags]
```

### Examples

```
apictl undeploy mcp-server -n MyMCPServer -v 1.0.0 --rev 2 -e dev
apictl undeploy mcp-server -n MyMCPServer -v 2.1.0 --rev 6 -g Label1 -g Label2 -g Label3 -e production
apictl undeploy mcp-server -n MyMCPServer -v 2.1.0 -r alice --rev 2 -g Label1 -e production
NOTE: All the 4 flags (--name (-n), --version (-v), --rev, --environment (-e)) are mandatory.
If the flag (--gateway-env (-g)) is not provided, revision will be undeployed from all the deployed gateway environments.
```

### Options

```
  -e, --environment string    Environment of which the MCP Server should be undeployed
  -g, --gateway-env strings   Gateway environment which the revision has to be undeployed
  -h, --help                  help for mcp-server
  -n, --name string           Name of the MCP Server to be undeployed
  -r, --provider string       Provider of the MCP Server
      --rev string            Revision number of the MCP Server to undeploy
  -v, --version string        Version of the MCP Server to be undeployed
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl undeploy](apictl_undeploy.md)	 - Undeploy an API/MCP Server/API Product revision from a gateway environment

