## mi update user

Update roles of a user in a Micro Integrator

### Synopsis

Update the roles of a user named [user-name] specified by the command line arguments in a Micro Integrator in the environment specified by the flag --environment, -e

```
mi update user [user-name] [flags]
```

### Examples

```
To update the roles
   mi update user [user-name] -e dev
NOTE: The flag (--environment (-e)) is mandatory
```

### Options

```
  -e, --environment string   Environment of the Micro Integrator of which the user's roles should be updated
  -h, --help                 help for user
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [mi update](mi_update.md)	 - Update log level of Loggers in a Micro Integrator instance

