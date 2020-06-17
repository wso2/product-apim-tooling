## apictl login

Login to an API Manager

### Synopsis

Login to an API Manager using credentials

```
apictl login [environment] [flags]
```

### Examples

```
apictl login dev -u admin -p admin
apictl login dev -u admin
cat ~/.mypassword | apictl login dev -u admin
```

### Options

```
  -h, --help              help for login
  -p, --password string   Password for login
      --password-stdin    Get password from stdin (default false)
  -u, --username string   Username for login
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl](apictl.md)	 - CLI for Importing and Exporting APIs and Applications

