## apictl import mcp-server

Import MCP Server

### Synopsis

Import an MCP Server to an environment

```
apictl import mcp-server --file <path-to-mcp-server> --environment <environment> [flags]
```

### Examples

```
apictl import mcp-server -f qa/ChoreoConnect.zip -e dev
apictl import mcp-server -f staging/ChoreoConnect.zip -e production
apictl import mcp-server -f ~/my-mcp-server -e production --update --rotate-revision
apictl import mcp-server -f ~/my-mcp-server -e production --update
NOTE: Both the flags (--file (-f) and --environment (-e)) are mandatory
```

### Options

```
      --dry-run              Get verification of the governance compliance of the MCP Server without importing it
  -e, --environment string   Environment from the which the MCP Server should be imported
  -f, --file string          Name of the MCP Server to be imported
      --format string        Output format of violation results in dry-run mode. Supported formats: [table, json, list]. If not provided, the default format is table.
  -h, --help                 help for mcp-server
      --params string        Provide an API Manager params file or a directory generated using "gen deployment-dir" command
      --preserve-provider    Preserve existing provider of MCP Server after importing (default true)
      --rotate-revision      Rotate the revisions with each update
      --skip-cleanup         Leave all temporary files created during import process
      --skip-deployments     Update only the working copy and skip deployment steps in import
      --update               Update an existing MCP Server or create a new MCP Server
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl import](apictl_import.md)	 - Import an API/MCP Server/API Product/Application to an environment

