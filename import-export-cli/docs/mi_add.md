## mi add

Add new users or loggers to a Micro Integrator instance

### Synopsis

Add new users or loggers to a Micro Integrator instance in the environment specified by the flag (--environment, -e)

```
mi add [flags]
```

### Examples

```
 mi add user capp-developer -e dev
 mi add log-level synapse-api org.apache.synapse.rest.API DEBUG -e dev
```

### Options

```
  -h, --help   help for add
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [mi](mi.md)	 - Micro Integrator related commands
* [mi add env](mi_add_env.md)	 - Add Environment to Config file
* [mi add log-level](mi_add_log-level.md)	 - Add new Logger to a Micro Integrator
* [mi add role](mi_add_role.md)	 - Add new role to a Micro Integrator
* [mi add user](mi_add_user.md)	 - Add new user to a Micro Integrator

