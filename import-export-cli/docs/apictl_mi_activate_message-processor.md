## apictl mi activate message-processor

Activate a message processor deployed in a Micro Integrator

### Synopsis

Activate the message processor specified by the command line argument [messageprocessor-name] deployed in a Micro Integrator in the environment specified by the flag --environment, -e

```
apictl mi activate message-processor [messageprocessor-name] [flags]
```

### Examples

```
To activate a message processor
  apictl mi activate message-processor TestMessageProcessor -e dev
NOTE: The flag (--environment (-e)) is mandatory
```

### Options

```
  -e, --environment string   Environment of the micro integrator in which the message processor should be activated
  -h, --help                 help for message-processor
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl mi activate](apictl_mi_activate.md)	 - Activate artifacts deployed in a Micro Integrator instance

