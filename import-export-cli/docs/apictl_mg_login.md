## apictl mg login

Login to a Microgateway Adapter environment

### Synopsis

Login to a Microgateway Adapter environment using username and password

```
apictl mg login [environment] [flags]
```

### Examples

```
apictl mg login dev -u admin -p admin
apictl mg login dev -u admin
cat ~/.mypassword | apictl mg login dev -u admin --password-stdin
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

* [apictl mg](apictl_mg.md)	 - Handle Microgateway related operations

