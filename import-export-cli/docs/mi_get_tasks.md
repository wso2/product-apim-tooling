## mi get tasks

Get information about tasks deployed in a Micro Integrator

### Synopsis

Get information about the tasks specified by command line argument [task-name]
If not specified, list all the tasks deployed in a Micro Integrator in the environment specified by the flag --environment, -e

```
mi get tasks [task-name] [flags]
```

### Examples

```
To list all the tasks
   mi get tasks -e dev
To get details about a specific tasks
   mi get tasks SampleTask -e dev
NOTE: The flag (--environment (-e)) is mandatory
```

### Options

```
  -e, --environment string   Environment to be searched
      --format string        Pretty-print using Go Templates. Use "{{ jsonPretty . }}" to list all fields
  -h, --help                 help for tasks
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [mi get](mi_get.md)	 - Get information about artifacts deployed in a Micro Integrator instance

