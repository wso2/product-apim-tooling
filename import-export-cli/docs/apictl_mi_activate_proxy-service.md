## apictl mi activate proxy-service

Activate a proxy service deployed in a Micro Integrator

### Synopsis

Activate the proxy service specified by the command line argument [proxy-name] deployed in a Micro Integrator in the environment specified by the flag --environment, -e

```
apictl mi activate proxy-service [proxy-name] [flags]
```

### Examples

```
To activate a proxy service
  apictl mi activate proxy-service SampleProxy -e dev
NOTE: The flag (--environment (-e)) is mandatory
```

### Options

```
  -e, --environment string   Environment of the micro integrator in which the proxy service should be activated
  -h, --help                 help for proxy-service
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl mi activate](apictl_mi_activate.md)	 - Activate artifacts deployed in a Micro Integrator instance

