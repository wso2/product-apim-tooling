## apimctl list

List APIs/Applications in an environment or List the environments

### Synopsis


Display a list containing all the APIs available in the environment specified by flag (--environment, -e)/
Display a list of Applications of a specific user in the environment specified by flag (--environment, -e)
OR
List all the environments

```
apimctl list [flags]
```

### Examples

```
apimctl list envs
apimctl list apis -e dev
```

### Options

```
  -h, --help   help for list
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO
* [apimctl](apimctl.md)	 - CLI for Importing and Exporting APIs and Applications
* [apimctl list apis](apimctl_list_apis.md)	 - Display a list of APIs in an environment
* [apimctl list apps](apimctl_list_apps.md)	 - Display a list of Applications in an environment specific to an owner
* [apimctl list envs](apimctl_list_envs.md)	 - Display the list of environments

