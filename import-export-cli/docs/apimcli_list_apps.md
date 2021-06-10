## apimcli list apps

Display a list of Applications in an environment specific to an owner

### Synopsis


Display a list of Applications of the user in the environment specified by the flag --environment, -e

```
apimcli list apps [flags]
```

### Examples

```
apimcli list apps -e dev
apimcli list apps -e dev -o sampleUser
apimcli list apps -e prod -o sampleUser
apimcli list apps -e staging -o sampleUser
apimcli list apps -e dev -l 40
```

### Options

```
  -e, --environment string   Environment to be searched
      --format string        Pretty-print outputusing Go templates. Use {{jsonPretty .}} to list all fields
  -h, --help                 help for apps
  -l, --limit string         Maximum number of applications to return
  -o, --owner string         Owner of the Application
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO
* [apimcli list](apimcli_list.md)	 - List APIs/Applications in an environment or List the environments

