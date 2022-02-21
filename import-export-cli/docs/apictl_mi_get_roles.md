## apictl mi get roles

Get information about roles

### Synopsis

Get information about the roles in primary and secondary user stores.
List all roles of the Micro Integrator in the environment specified by the flag --environment, -e

```
apictl mi get roles [role-name] [flags]
```

### Examples

```
To list all the roles
  apictl mi get roles -e dev
To get details about a role by providing the role name
  apictl mi get roles [role-name] -e dev
To get details about a role in a secondary user store
  apictl mi get roles [role-name] -d [domain] -e dev
NOTE: The flag (--environment (-e)) is mandatory
```

### Options

```
  -d, --domain string        Filter roles by domain
  -e, --environment string   Environment to be searched
      --format string        Pretty-print using Go Templates. Use "{{ jsonPretty . }}" to list all fields
  -h, --help                 help for roles
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl mi get](apictl_mi_get.md)	 - Get information about artifacts deployed in a Micro Integrator instance

