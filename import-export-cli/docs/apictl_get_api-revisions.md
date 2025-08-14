## apictl get api-revisions

Display a list of Revisions for the API

### Synopsis

Display a list of Revisions available for the API in the environment specified

```
apictl get api-revisions [flags]
```

### Examples

```
apictl get api-revisions -n PizzaAPI -v 1.0.0 -e dev
apictl get api-revisions -n TwitterAPI -v 1.0.0 -r admin -e dev
apictl get api-revisions -n PizzaShackAPI -v 1.0.0 -q deployed:true -e dev
NOTE: All the 3 flags (--name (-n), --version (-v) and --environment (-e)) are mandatory.
```

### Options

```
  -e, --environment string   Environment to be searched
      --format string        Pretty-print revisions using Go Templates. Use "{{ jsonPretty . }}" to list all fields
  -h, --help                 help for api-revisions
  -n, --name string          Name of the API to get the revision
  -r, --provider string      Provider of the API
  -q, --query strings        Query pattern
  -v, --version string       Version of the API to get the revision
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl get](apictl_get.md)	 - Get APIs/MCPServers/APIProducts/Applications or revisions of a specific API/MCPServers/APIProduct in an environment or Get the Correlation Log Configurations or Get the log level of each API/MCPServers in an environment or Get the environments

