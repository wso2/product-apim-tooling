## apictl mi delete user

Delete a user from the Micro Integrator

### Synopsis

Delete a user with the name specified by the command line argument [user-name] from a Micro Integrator in the environment specified by the flag --environment, -e

```
apictl mi delete user [user-name] [flags]
```

### Examples

```
To delete a user
  apictl mi delete user capp-tester -e dev
NOTE: The flag (--environment (-e)) is mandatory
```

### Options

```
  -e, --environment string   Environment of the micro integrator from which a user should be deleted
  -h, --help                 help for user
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl mi delete](apictl_mi_delete.md)	 - Delete users from a Micro Integrator instance

