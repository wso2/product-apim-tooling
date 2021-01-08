## apictl mi activate

Activate artifacts deployed in a Micro Integrator instance

### Synopsis

Activate artifacts deployed in a Micro Integrator instance in the environment specified by the flag (--environment, -e)

```
apictl mi activate [flags]
```

### Examples

```
apictl mi activate endpoint TestEP -e dev
```

### Options

```
  -h, --help   help for activate
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl mi](apictl_mi.md)	 - Micro Integrator related commands
* [apictl mi activate endpoint](apictl_mi_activate_endpoint.md)	 - Activate a endpoint deployed in a Micro Integrator
* [apictl mi activate message-processor](apictl_mi_activate_message-processor.md)	 - Activate a message processor deployed in a Micro Integrator
* [apictl mi activate proxy-service](apictl_mi_activate_proxy-service.md)	 - Activate a proxy service deployed in a Micro Integrator

