## apictl import api

Import API

### Synopsis

Import an API to an environment

```
apictl import api --file <path-to-api> --environment <environment> [flags]
```

### Examples

```
apictl import api -f qa/TwitterAPI.zip -e dev
apictl import api -f staging/FacebookAPI.zip -e production
apictl import api -f ~/myapi -e production --update --rotate-revision
apictl import api -f ~/myapi -e production --update
NOTE: Both the flags (--file (-f) and --environment (-e)) are mandatory
```

### Options

```
  -e, --environment string   Environment from the which the API should be imported
  -f, --file string          Name of the API to be imported
  -h, --help                 help for api
      --params string        Provide an API Manager params file or a directory generated using "gen deployment-dir" command
      --preserve-provider    Preserve existing provider of API after importing (default true)
      --rotate-revision      Rotate the revisions with each update
      --skip-cleanup         Leave all temporary files created during import process
      --skip-deployments     Update only the working copy and skip deployment steps in import
      --update               Update an existing API or create a new API
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl import](apictl_import.md)	 - Import an API/API Product/Application to an environment

