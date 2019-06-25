## apimcli login

Login to an API Manager

### Synopsis


Login to an API Manager using credentials

```
apimcli login [environment] [flags]
```

### Examples

```
apimcli login dev -u admin -p admin
apimcli login dev -u admin
cat ~/.mypassword | apimcli login dev -u admin
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
* [apimcli](apimcli.md)	 - CLI for Importing and Exporting APIs and Applications

