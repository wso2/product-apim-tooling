## apictl get apps

Display a list of Applications in an environment specific to an owner

### Synopsis

Display a list of Applications of the user in the environment specified by the flag --environment, -e

```
apictl get apps [flags]
```

### Examples

```
apictl get apps -e dev 
apictl get apps -e dev -o sampleUser
apictl get apps -e prod -o sampleUser
apictl get apps -e staging -o sampleUser
apictl get apps -e dev -l 40
NOTE: The flag (--environment (-e)) is mandatory
```

### Options

```
  -e, --environment string   Environment to be searched
      --format string        Pretty-print outputusing Go templates. Use "{{jsonPretty .}}" to list all fields
  -h, --help                 help for apps
  -l, --limit string         Maximum number of applications to return (default "25")
  -o, --owner string         Owner of the Application
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl get](apictl_get.md)	 - Get APIs/MCPServers/APIProducts/Applications or revisions of a specific API/MCPServers/APIProduct in an environment or Get the Correlation Log Configurations or Get the log level of each API/MCPServers in an environment or Get the environments

