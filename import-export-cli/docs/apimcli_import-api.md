## apimcli import-api

Import API

### Synopsis


Import an API to an environment

```
apimcli import-api --file <Path to API> --environment <Environment to be imported> [flags]
```

### Examples

```
apimcli import-api -f qa/TwitterAPI.zip -e dev
apimcli import-api -f staging/FacebookAPI.zip -e production -u admin -p admin
apimcli import-api -f ~/myapi -e production -u admin -p admin --update
apimcli import-api -f ~/myapi -e production -u admin -p admin --update --inject
```

### Options

```
  -e, --environment string   Environment from the which the API should be imported
  -f, --file string          Name of the API to be imported
  -h, --help                 help for import-api
      --inject               Inject variables definedin params file to the given API.
      --params string        Provide a API Manager params file (default "api_params.yaml")
  -p, --password string      Password
      --preserve-provider    Preserve existing provider of API after exporting (default true)
      --update               Update API if exists. Otherwise it will create API
  -u, --username string      Username
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO
* [apimcli](apimcli.md)	 - CLI for Importing and Exporting APIs and Applications

