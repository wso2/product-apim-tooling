## apictl list envs

Display the list of environments

### Synopsis

Display a list of environments defined in 'main_config.yaml' file

```
apictl list envs [flags]
```

### Examples

```
apictl list envs
```

### Options

```
      --format string   Pretty-print environments using go templates (default "table {{.Name}}\t{{.ApiManagerEndpoint}}\t{{.RegistrationEndpoint}}\t{{.TokenEndpoint}}\t{{.PublisherEndpoint}}\t{{.ApplicationEndpoint}}\t{{.AdminEndpoint}}")
  -h, --help            help for envs
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl list](apictl_list.md)	 - List APIs/Applications in an environment or List the environments

