## apictl mi add user

Add new user to a Micro Integrator

### Synopsis

Add a new user with the name specified by the command line argument [user-name] to a Micro Integrator in the environment specified by the flag --environment, -e

```
apictl mi add user [user-name] [flags]
```

### Examples

```
To add a new user
  apictl mi add user capp-tester -e dev
NOTE: The flag (--environment (-e)) is mandatory
```

### Options

```
  -e, --environment string   Environment of the micro integrator to which a new user should be added
  -h, --help                 help for user
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl mi add](apictl_mi_add.md)	 - Add new users or loggers to a Micro Integrator instance

