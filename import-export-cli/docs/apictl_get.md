## apictl get

Get APIs/APIProducts/Applications in an environment or Get the environments

### Synopsis

Display a list containing all the APIs available in the environment specified by flag (--environment, -e)/
Display a list containing all the API Products available in the environment specified by flag (--environment, -e)/
Display a list of Applications of a specific user in the environment specified by flag (--environment, -e)
OR
List all the environments

```
apictl get [flags]
```

### Examples

```
apictl get envs
apictl get apis -e dev
apictl get api-products -e dev
apictl get apps -e dev
```

### Options

```
  -h, --help   help for get
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl](apictl.md)	 - CLI for Importing and Exporting APIs and Applications and Managing WSO2 Micro Integrator
* [apictl get api-products](apictl_get_api-products.md)	 - Display a list of API Products in an environment
* [apictl get apis](apictl_get_apis.md)	 - Display a list of APIs in an environment
* [apictl get apps](apictl_get_apps.md)	 - Display a list of Applications in an environment specific to an owner
* [apictl get envs](apictl_get_envs.md)	 - Display the list of environments
* [apictl get keys](apictl_get_keys.md)	 - Generate access token to invoke the API or API Product
* [apictl get revisions](apictl_get_revisions.md)	 - Display a list of Revisions for the API

