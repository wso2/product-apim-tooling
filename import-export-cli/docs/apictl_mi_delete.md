## apictl mi delete

Delete users from a Micro Integrator instance

### Synopsis

Delete users from a Micro Integrator instance in the environment specified by the flag (--environment, -e)

```
apictl mi delete [flags]
```

### Examples

```
apictl mi delete user capp-tester -e dev
```

### Options

```
  -h, --help   help for delete
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl mi](apictl_mi.md)	 - Micro Integrator related commands
* [apictl mi delete role](apictl_mi_delete_role.md)	 - Delete a role from the Micro Integrator
* [apictl mi delete user](apictl_mi_delete_user.md)	 - Delete a user from the Micro Integrator

