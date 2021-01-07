## apictl mi add

Add new users or loggers to a Micro Integrator instance

### Synopsis

Add new users or loggers to a Micro Integrator instance in the environment specified by the flag (--environment, -e)

```
apictl mi add [flags]
```

### Examples

```
apictl mi add user capp-developer -e dev
apictl mi add log-level synapse-api org.apache.synapse.rest.API DEBUG -e dev
```

### Options

```
  -h, --help   help for add
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl mi](apictl_mi.md)	 - Micro Integrator related commands
* [apictl mi add log-level](apictl_mi_add_log-level.md)	 - Add new Logger to a Micro Integrator
* [apictl mi add user](apictl_mi_add_user.md)	 - Add new user to a Micro Integrator

