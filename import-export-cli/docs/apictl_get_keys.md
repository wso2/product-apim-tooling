## apictl get keys

Generate access token to invoke the API or API Product or MCP Server

### Synopsis

Generate JWT token to invoke the API or API Product or MCP Server by subscribing to a default application for testing purposes

```
apictl get keys [flags]
```

### Examples

```
apictl get keys -n TwitterAPI -v 1.0.0 -e dev --provider admin
NOTE: Both the flags (--name (-n) and --environment (-e)) are mandatory.
You can override the default token endpoint using --token (-t) optional flag providing a new token endpoint
```

### Options

```
  -e, --environment string   Key generation environment
  -h, --help                 help for keys
  -n, --name string          API or API Product or MCP Server to generate keys
  -r, --provider string      Provider of the API or API Product or MCP Server
  -t, --token string         Token endpoint URL of Environment
  -v, --version string       Version of the API or API Product or MCP Server
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl get](apictl_get.md)	 - Get APIs/MCPServers/APIProducts/Applications or revisions of a specific API/MCPServers/APIProduct in an environment or Get the Correlation Log Configurations or Get the log level of each API/MCPServers in an environment or Get the environments

