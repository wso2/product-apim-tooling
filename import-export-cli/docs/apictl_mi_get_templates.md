## apictl mi get templates

Get information about templates deployed in a Micro Integrator

### Synopsis

Get information about the template specified by command line arguments [template-type] and [template-name]
If not specified, list all the templates in the environment specified by the flag --environment, -e

```
apictl mi get templates [template-type] [template-name] [flags]
```

### Examples

```
To list all the templates
apictl mi get templates -e dev
To get details about a specific template type
apictl mi get templates TemplateType
To get details about a specific template
apictl mi get templates TemplateType TemplateName -e dev
NOTE: The flag (--environment (-e)) is mandatory
```

### Options

```
  -e, --environment string   Environment to be searched
      --format string        Pretty-print templates using Go Templates. Use "{{ jsonPretty . }}" to list all fields
  -h, --help                 help for templates
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl mi get](apictl_mi_get.md)	 - Get information about artifacts deployed in a Micro Integrator instance

