## apimctl list apps

Display a list of Applications in an environment specific to an owner

### Synopsis


Display a list of Applications of the user in the environment specified by the flag --environment, -e

```
apimctl list apps [flags]
```

### Examples

```
apimctl list apps -e dev
apimctl list apps -e dev -o sampleUser
apimctl list apps -e prod -o sampleUser
apimctl list apps -e staging -o sampleUser
```

### Options

```
  -e, --environment string   Environment to be searched
      --format string        Pretty-print outputusing Go templates. Use {{jsonPretty .}} to list all fields
  -h, --help                 help for apps
  -o, --owner string         Owner of the Application
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO
* [apimctl list](apimctl_list.md)	 - List APIs/Applications in an environment or List the environments

