## apictl mi deactivate

Deactivate artifacts deployed in a Micro Integrator instance

### Synopsis

Deactivate artifacts deployed in a Micro Integrator instance in the environment specified by the flag (--environment, -e)

```
apictl mi deactivate [flags]
```

### Examples

```
apictl mi deactivate endpoint TestEP -e dev
```

### Options

```
  -h, --help   help for deactivate
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl mi](apictl_mi.md)	 - Micro Integrator related commands
* [apictl mi deactivate endpoint](apictl_mi_deactivate_endpoint.md)	 - Deactivate a endpoint deployed in a Micro Integrator
* [apictl mi deactivate message-processor](apictl_mi_deactivate_message-processor.md)	 - Deactivate a message processor deployed in a Micro Integrator
* [apictl mi deactivate proxy-service](apictl_mi_deactivate_proxy-service.md)	 - Deactivate a proxy service deployed in a Micro Integrator

