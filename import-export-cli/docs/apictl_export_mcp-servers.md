## apictl export mcp-servers

Export MCP Servers for migration

### Synopsis

Export all the MCP Servers of a tenant from one environment, to be imported into another environment

```
apictl export mcp-servers (--environment <environment-from-which-artifacts-should-be-exported> --format <export-format> --preserve-status --force) [flags]
```

### Examples

```
apictl export mcp-servers -e production --force
apictl export mcp-servers -e production
NOTE: The flag (--environment (-e)) is mandatory
```

### Options

```
      --all                    Export working copy and all revisions for the MCP Servers in the environments 
  -e, --environment string     Environment from which the MCP Servers should be exported
      --force                  Clean all the previously exported MCP Servers of the given target tenant, in the given environment if any, and to export MCP Servers from beginning
      --format string          File format of exported archives(json or yaml) (default "YAML")
  -h, --help                   help for mcp-servers
      --preserve-credentials   Preserve endpoint credentials when exporting. Otherwise credentials will not be exported
      --preserve-status        Preserve MCP Server status when exporting. Otherwise MCP Server will be exported in CREATED status (default true)
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl export](apictl_export.md)	 - Export an API/MCPServer/API Product/Application/Policy in an environment

