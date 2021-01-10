## apictl mi update log-level

Update log level of a Logger in a Micro Integrator

### Synopsis

Update the log level of a Logger named [logger-name] to [log-level] specified by the command line arguments in a Micro Integrator in the environment specified by the flag --environment, -e

```
apictl mi update log-level [logger-name] [log-level] [flags]
```

### Examples

```
To update the log level
  apictl mi update log-level org-apache-coyote DEBUG -e dev
NOTE: The flag (--environment (-e)) is mandatory
```

### Options

```
  -e, --environment string   Environment of the micro integrator of which the logger should be updated
  -h, --help                 help for log-level
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl mi update](apictl_mi_update.md)	 - Update log level of Loggers in a Micro Integrator instance

