## mi update hashicorp-secret

Update the secret ID of HashiCorp configuration in a Micro Integrator

### Synopsis

Update the secret ID of the HashiCorp configuration in a Micro Integrator in the environment specified by the flag --environment, -e

```
mi update hashicorp-secret [secret-id] [flags]
```

### Examples

```
To update the secret ID
   mi update hashicorp-secret new_secret_id -e dev
NOTE: The flag (--environment (-e)) is mandatory
```

### Options

```
  -e, --environment string   Environment of the micro integrator of which the HashiCorp secret ID should be updated
  -h, --help                 help for hashicorp-secret
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [mi update](mi_update.md)	 - Update log level of Loggers in a Micro Integrator instance

