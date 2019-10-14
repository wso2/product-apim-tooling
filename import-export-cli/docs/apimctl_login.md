## apimctl login

Login to an API Manager

### Synopsis


Login to an API Manager using credentials

```
apimctl login [environment] [flags]
```

### Examples

```
apimctl login dev -u admin -p admin
apimctl login dev -u admin
cat ~/.mypassword | apimctl login dev -u admin
```

### Options

```
  -h, --help              help for login
  -p, --password string   Password for login
      --password-stdin    Get password from stdin
  -u, --username string   Username for login
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO
* [apimctl](apimctl.md)	 - CLI for Importing and Exporting APIs and Applications

