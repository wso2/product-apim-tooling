## apictl list

List APIs/Applications in an environment or List the environments

### Synopsis

Display a list containing all the APIs available in the environment specified by flag (--environment, -e)/
Display a list of Applications of a specific user in the environment specified by flag (--environment, -e)
OR
List all the environments

```
apictl list [flags]
```

### Examples

```
apictl list envs
apictl list apis -e dev
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

* [apictl](apictl.md)	 - CLI for Importing and Exporting APIs and Applications
* [apictl list apis](apictl_list_apis.md)	 - Display a list of APIs in an environment
* [apictl list apps](apictl_list_apps.md)	 - Display a list of Applications in an environment specific to an owner
* [apictl list envs](apictl_list_envs.md)	 - Display the list of environments

