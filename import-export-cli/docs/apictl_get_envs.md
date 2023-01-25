## apictl get envs

Display the list of environments

### Synopsis

Display a list of environments defined in 'main_config.yaml' file

```
apictl get envs [flags]
```

### Examples

```
apictl list envs
```

### Options

```
      --format string   Pretty-print environments using go templates (default "table {{.Name}}\t{{.ApiManagerEndpoint}}\t{{.RegistrationEndpoint}}\t{{.TokenEndpoint}}\t{{.PublisherEndpoint}}\t{{.ApplicationEndpoint}}\t{{.AdminEndpoint}}\t{{.MiManagementEndpoint}}")
  -h, --help            help for envs
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl get](apictl_get.md)	 - Get APIs/APIProducts/Applications or revisions of a specific API/APIProduct in an environment or Get the Correlation Log Configurations or Get the log level of each API in an environment or Get the environments

