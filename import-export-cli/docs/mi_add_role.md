## mi add role

Add new role to a Micro Integrator

### Synopsis

Add a new role with the name specified by the command line argument [role-name] to a Micro Integrator in the environment specified by the flag --environment, -e

```
mi add role [role-name] [flags]
```

### Examples

```
To add a new role
   mi add role [role-name] -e dev
NOTE: The flag (--environment (-e)) is mandatory
```

### Options

```
  -e, --environment string   Environment of the micro integrator to which a new user should be added
  -h, --help                 help for role
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [mi add](mi_add.md)	 - Add new users or loggers to a Micro Integrator instance

