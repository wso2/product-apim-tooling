## apictl mg deploy

Deploy apictl project.

### Synopsis

Deploy the apictl project in Microgateway

```
apictl mg deploy --host [control plane url] --file [file name] --username [username] --password [password] [flags]
```

### Examples

```
apictl mg deploy -h https://localhost:9095 -f qa/TwitterAPI.zip -u admin -p admin
cat ~/.mypassword | apictlmg  deploy -h https://localhost:9095 -f qa/TwitterAPI.zip -u admin
```

### Options

```
  -f, --file string       Provide the filepath of the apictl project to be imported
  -h, --help              help for deploy
  -c, --host string       Provide the host url for the control plane with port
  -p, --password string   Provide the password
      --skipCleanup       Leave all temporary files created during import process
  -u, --username string   Provide the username
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl mg](apictl_mg.md)	 - Handle Microgateway related operations

