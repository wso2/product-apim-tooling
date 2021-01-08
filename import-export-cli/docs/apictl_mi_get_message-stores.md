## apictl mi get message-stores

Get information about message stores deployed in a Micro Integrator

### Synopsis

Get information about the message stores specified by command line argument [messagestore-name]
If not specified, list all the message stores deployed in a Micro Integrator in the environment specified by the flag --environment, -e

```
apictl mi get message-stores [messagestore-name] [flags]
```

### Examples

```
To list all the message stores
  apictl mi get message-stores -e dev
To get details about a specific message stores
  apictl mi get message-stores TestMessageStore -e dev
NOTE: The flag (--environment (-e)) is mandatory
```

### Options

```
  -e, --environment string   Environment to be searched
      --format string        Pretty-print using Go Templates. Use "{{ jsonPretty . }}" to list all fields
  -h, --help                 help for message-stores
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl mi get](apictl_mi_get.md)	 - Get information about artifacts deployed in a Micro Integrator instance

