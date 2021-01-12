## apictl mi activate endpoint

Activate a endpoint deployed in a Micro Integrator

### Synopsis

Activate the endpoint specified by the command line argument [endpoint-name] deployed in a Micro Integrator in the environment specified by the flag --environment, -e

```
apictl mi activate endpoint [endpoint-name] [flags]
```

### Examples

```
To activate a endpoint
  apictl mi activate endpoint TestEP -e dev
NOTE: The flag (--environment (-e)) is mandatory
```

### Options

```
  -e, --environment string   Environment of the micro integrator in which the endpoint should be activated
  -h, --help                 help for endpoint
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl mi activate](apictl_mi_activate.md)	 - Activate artifacts deployed in a Micro Integrator instance

