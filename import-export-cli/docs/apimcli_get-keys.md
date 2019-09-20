## apimcli get-keys

Generate access token 

### Synopsis


Get API subscribed to a default application and generate access token

```
apimcli get-keys (--name <name-of-the-api> --version <version-of-the-api> --environment <environment-from-which-the-api-should-be-exported> --provider <API provider name>) [flags]
```

### Examples

```
apimcli get-keys -n TwitterAPI -v 1.0.0 -e dev --provider admin
apimcli get-keys -n FacebookAPI -v 2.1.0 -e production --provider admin
NOTE: all three flags (--name (-n), --version (-v), --provider (-r)) are mandatory
```

### Options

```
  -e, --environment string   Environment to which the API should be exported
  -h, --help                 help for export-api
  -n, --name string          Name of the API to be exported
      --provider string      Provider of the API
  -v, --version string       Version of the API to be exported
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO
* [apimcli](apimcli.md)	 - CLI for Importing and Exporting APIs and Applications

