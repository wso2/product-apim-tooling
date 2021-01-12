## apictl mi deactivate endpoint

Deactivate a endpoint deployed in a Micro Integrator

### Synopsis

Deactivate the endpoint specified by the command line argument [endpoint-name] deployed in a Micro Integrator in the environment specified by the flag --environment, -e

```
apictl mi deactivate endpoint [endpoint-name] [flags]
```

### Examples

```
To deactivate a endpoint
  apictl mi deactivate endpoint TestEP -e dev
NOTE: The flag (--environment (-e)) is mandatory
```

### Options

```
  -e, --environment string   Environment of the micro integrator in which the endpoint should be deactivated
  -h, --help                 help for endpoint
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl mi deactivate](apictl_mi_deactivate.md)	 - Deactivate artifacts deployed in a Micro Integrator instance

