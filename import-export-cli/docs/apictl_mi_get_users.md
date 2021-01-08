## apictl mi get users

Get information about users

### Synopsis

Get information about the users filtered by username pattern and role.
If not provided list all users of the Micro Integrator in the environment specified by the flag --environment, -e

```
apictl mi get users [user-name] [flags]
```

### Examples

```
Example:
To list all the users
  apictl mi get users -e dev
To get the list of users with specific role
  apictl mi get users -r [role-name] -e dev
To get the list of users with a username matching with the wild card Ex: "*mi*" matches with "admin"
  apictl mi get users -p [pattern] -e dev
To get details about a user by providing the user-id
  apictl mi get users [user-id] -e dev
NOTE: The flag (--environment (-e)) is mandatory
```

### Options

```
  -e, --environment string   Environment to be searched
      --format string        Pretty-print using Go Templates. Use "{{ jsonPretty . }}" to list all fields
  -h, --help                 help for users
  -p, --pattern string       Filter users by regex
  -r, --role string          Filter users by role
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl mi get](apictl_mi_get.md)	 - Get information about artifacts deployed in a Micro Integrator instance

