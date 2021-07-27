## apimcli import-api

Import API

### Synopsis


Import an API to an environment
Examples:
apimcli import-api -f qa/TwitterAPI.zip -e dev
apimcli import-api -f staging/FacebookAPI.zip -e production -u admin -p admin


```
apimcli import-api (--file <api-zip-file> --environment <environment-to-which-the-api-should-be-imported>) [flags]
```

### Options

```
  -e, --environment string   Environment from the which the API should be imported (default "default")
  -f, --file string          Name of the API to be imported
  -h, --help                 help for import-api
  -p, --password string      Password
      --preserve-provider    Preserve existing provider of API after exporting (default true)
  -u, --username string      Username
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO
* [apimcli](apimcli.md)	 - CLI for Importing and Exporting APIs and Applications

