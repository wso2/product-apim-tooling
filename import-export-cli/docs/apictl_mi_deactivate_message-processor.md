## apictl mi deactivate message-processor

Deactivate a message processor deployed in a Micro Integrator

### Synopsis

Deactivate the message processor specified by the command line argument [messageprocessor-name] deployed in a Micro Integrator in the environment specified by the flag --environment, -e

```
apictl mi deactivate message-processor [messageprocessor-name] [flags]
```

### Examples

```
To deactivate a message processor
  apictl mi deactivate message-processor TestMessageProcessor -e dev
NOTE: The flag (--environment (-e)) is mandatory
```

### Options

```
  -e, --environment string   Environment of the micro integrator in which the message processor should be deactivated
  -h, --help                 help for message-processor
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl mi deactivate](apictl_mi_deactivate.md)	 - Deactivate artifacts deployed in a Micro Integrator instance

