## apictl mi get message-processors

Get information about message processors deployed in a Micro Integrator

### Synopsis

Get information about the message processors specified by command line argument [messageprocessor-name]
If not specified, list all the message processors deployed in a Micro Integrator in the environment specified by the flag --environment, -e

```
apictl mi get message-processors [messageprocessor-name] [flags]
```

### Examples

```
To list all the message processors
  apictl mi get message-processors -e dev
To get details about a specific message processors
  apictl mi get message-processors TestMessageProcessor -e dev
NOTE: The flag (--environment (-e)) is mandatory
```

### Options

```
  -e, --environment string   Environment to be searched
      --format string        Pretty-print message processors using Go Templates. Use "{{ jsonPretty . }}" to list all fields
  -h, --help                 help for message-processors
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl mi get](apictl_mi_get.md)	 - Get information about artifacts deployed in a Micro Integrator instance

