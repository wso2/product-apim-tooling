## apictl get policies

Get Policy list

### Synopsis

Get a list of Policies in an environment

```
apictl get policies [flags]
```

### Examples

```
apictl get policies rate-limiting -e production -q type:sub
```

### Options

```
  -h, --help   help for policies
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl get](apictl_get.md)	 - Get APIs/MCPServers/APIProducts/Applications or revisions of a specific API/MCPServers/APIProduct in an environment or Get the Correlation Log Configurations or Get the log level of each API/MCPServers in an environment or Get the environments
* [apictl get policies api](apictl_get_policies_api.md)	 - Display a list of API Policies
* [apictl get policies rate-limiting](apictl_get_policies_rate-limiting.md)	 - Display a list of APIs in an environment

