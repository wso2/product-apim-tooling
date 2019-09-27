## apimcli import-api

Import API

### Synopsis


Import an API to an environment

```
apimcli import-api --file <PATH_TO_API> --environment <ENVIRONMENT> [flags]
```

### Examples

```
apimcli import-api -f qa/TwitterAPI.zip -e dev
apimcli import-api -f staging/FacebookAPI.zip -e production
apimcli import-api -f ~/myapi -e production --update
apimcli import-api -f ~/myapi -e production --update --inject
```

### Options

```
  -e, --environment string   Environment from the which the API should be imported
  -f, --file string          Name of the API to be imported
  -h, --help                 help for import-api
      --params string        Provide a API Manager params file (default "api_params.yaml")
      --preserve-provider    Preserve existing provider of API after importing (default true)
      --skipCleanup          Leave all temporary files created during import process
      --update               Update an existing API or create a new API
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO
* [apimcli](apimcli.md)	 - CLI for Importing and Exporting APIs and Applications

