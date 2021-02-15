## apictl mg deploy

Deploy apictl project.

### Synopsis

Deploy the apictl project in Microgateway

```
apictl mg deploy --host [control plane url] --file [file name] --username [username] --password [password] [flags]
```

### Examples

```
apictl mg deploy -h https://localhost:9095 -f petstore -u admin -p admin
cat ~/.mypassword | apictl mg  deploy -h https://localhost:9095 -f petstore -u admin
```

### Options

```
  -f, --file string       Filepath of the apictl project to be deployed
  -h, --help              help for deploy
  -c, --host string       Host url for the control plane with port
  -o, --overwrite         Whether to update an existing API
  -p, --password string   Password of the user (Can be provided at the prompt)
      --skipCleanup       Whether to keep all temporary files created during deploy process
  -u, --username string   Username of the user
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl mg](apictl_mg.md)	 - Handle Microgateway related operations

