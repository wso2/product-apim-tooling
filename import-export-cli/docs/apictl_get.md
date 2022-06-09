## apictl get

Get APIs/APIProducts/Applications or revisions of a specific API/APIProduct in an environment or Get the log level of each API in an environment or Get the environments

### Synopsis

Display a list containing all the APIs available in the environment specified by flag (--environment, -e)/
Display a list containing all the API Products available in the environment specified by flag (--environment, -e)/
Display a list of Applications of a specific user in the environment specified by flag (--environment, -e)/
Display a list of API revisions of a specific API in the environment specified by flag (--environment, -e)/
Display a list of API Product revisions of a specific API Product in the environment specified by flag (--environment, -e)/
Get a generated JWT token to invoke an API or API Product by subscribing to a default application for testing purposes in the environment specified by flag (--environment, -e)/
Get the log level of each API in the environment specified by flag (--environment, -e)
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
apictl get api-revisions -n PizzaAPI -v 1.0.0 -e dev
apictl get api-product-revisions -n PizzaProduct -v 1.0.0 -e dev
apictl get keys -n TwitterAPI -v 1.0.0 -e dev
apictl get api-logging -e dev --tenant-domain carbon.super
apictl get api-logging --api-id bf36ca3a-0332-49ba-abce-e9992228ae06 -e dev --tenant-domain carbon.super
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
* [apictl get api-logging](apictl_get_api-logging.md)	 - Display a list of API loggers in an environment
* [apictl get api-product-revisions](apictl_get_api-product-revisions.md)	 - Display a list of Revisions for the API Products
* [apictl get api-products](apictl_get_api-products.md)	 - Display a list of API Products in an environment
* [apictl get api-revisions](apictl_get_api-revisions.md)	 - Display a list of Revisions for the API
* [apictl get apis](apictl_get_apis.md)	 - Display a list of APIs in an environment
* [apictl get apps](apictl_get_apps.md)	 - Display a list of Applications in an environment specific to an owner
* [apictl get envs](apictl_get_envs.md)	 - Display the list of environments
* [apictl get keys](apictl_get_keys.md)	 - Generate access token to invoke the API or API Product
* [apictl get policies](apictl_get_policies.md)	 - Get Policy list

