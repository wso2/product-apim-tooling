## apictl set mcp-server-logging

Set the log level for an MCP Server in an environment

### Synopsis

Set the log level for an MCP Server in the environment specified

```
apictl set mcp-server-logging [flags]
```

### Examples

```
apictl set mcp-server-logging --mcp-server-id bf36ca3a-0332-49ba-abce-e9992228ae06 --log-level full -e dev --tenant-domain carbon.super
apictl set mcp-server-logging --mcp-server-id bf36ca3a-0332-49ba-abce-e9992228ae06 --log-level off -e dev --tenant-domain carbon.super
```

### Options

```
  -e, --environment string     Environment of the MCP Server which the log level should be set
  -h, --help                   help for mcp-server-logging
      --log-level string       Log Level
  -i, --mcp-server-id string   MCP Server ID
      --tenant-domain string   Tenant Domain
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl set](apictl_set.md)	 - Set configuration parameters, per API log levels, MCP Server log levels or correlation component configurations

