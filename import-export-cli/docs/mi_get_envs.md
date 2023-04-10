## mi get envs

Display the list of environments

### Synopsis

Display a list of environments defined in 'main_config.yaml' file

```
mi get envs [flags]
```

### Examples

```
 list envs
```

### Options

```
      --format string   Pretty-print environments using go templates (default "table {{.Name}}\t{{.MiManagementEndpoint}}")
  -h, --help            help for envs
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [mi get](mi_get.md)	 - Get information about artifacts deployed in a Micro Integrator instance

