## apictl mi update

Update log level of Loggers in a Micro Integrator instance

### Synopsis

Update log level of Loggers in a Micro Integrator instance in the environment specified by the flag (--environment, -e)

```
apictl mi update [flags]
```

### Examples

```
apictl mi update log-level org-apache-coyote DEBUG -e dev
```

### Options

```
  -h, --help   help for update
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl mi](apictl_mi.md)	 - Micro Integrator related commands
* [apictl mi update hashicorp-secret](apictl_mi_update_hashicorp-secret.md)	 - Update the secret ID of HashiCorp configuration in a Micro Integrator
* [apictl mi update log-level](apictl_mi_update_log-level.md)	 - Update log level of a Logger in a Micro Integrator

