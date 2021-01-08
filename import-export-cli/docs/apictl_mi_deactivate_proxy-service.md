## apictl mi deactivate proxy-service

Deactivate a proxy service deployed in a Micro Integrator

### Synopsis

Deactivate the proxy service specified by the command line argument [proxy-name] deployed in a Micro Integrator in the environment specified by the flag --environment, -e

```
apictl mi deactivate proxy-service [proxy-name] [flags]
```

### Examples

```
To deactivate a proxy service
  apictl mi deactivate proxy-service SampleProxy -e dev
NOTE: The flag (--environment (-e)) is mandatory
```

### Options

```
  -e, --environment string   Environment of the micro integrator in which the proxy service should be deactivated
  -h, --help                 help for proxy-service
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl mi deactivate](apictl_mi_deactivate.md)	 - Deactivate artifacts deployed in a Micro Integrator instance

