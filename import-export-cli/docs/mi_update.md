## mi update

Update log level of Loggers in a Micro Integrator instance

### Synopsis

Update log level of Loggers in a Micro Integrator instance in the environment specified by the flag (--environment, -e)

```
mi update [flags]
```

### Examples

```
 mi update log-level org-apache-coyote DEBUG -e dev
```

### Options

```
  -h, --help   help for update
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [mi](mi.md)	 - Micro Integrator related commands
* [mi update hashicorp-secret](mi_update_hashicorp-secret.md)	 - Update the secret ID of HashiCorp configuration in a Micro Integrator
* [mi update log-level](mi_update_log-level.md)	 - Update log level of a Logger in a Micro Integrator
* [mi update user](mi_update_user.md)	 - Update roles of a user in a Micro Integrator

