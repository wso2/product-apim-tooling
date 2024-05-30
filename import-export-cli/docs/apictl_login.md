## apictl login

Login to an API Manager

### Synopsis

Login to an API Manager using credentials or set token for authentication

```
apictl login [environment] [flags]
```

### Examples

```
apictl login dev -u admin -p admin
apictl login dev -u admin
cat ~/.mypassword | apictl login dev -u admin
apictl login dev --token e79bda48-3406-3178-acce-f6e4dbdcbb12
```

### Options

```
  -h, --help              help for login
  -p, --password string   Password for login
      --password-stdin    Get password from stdin
      --token string      Personal access token
  -u, --username string   Username for login
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl](apictl.md)	 - CLI for Importing and Exporting APIs and Applications and Managing WSO2 Micro Integrator

