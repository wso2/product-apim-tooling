## apictl mg deploy

Deploy an API (apictl project) in Microgateway

### Synopsis

Deploy an API (apictl project) in Microgateway by specifying the adapter host url.

```
apictl mg deploy [flags]
```

### Examples

```
apictl mg deploy -h https://localhost:9095 -f petstore -u admin -p admin

Note: The flags --host (-c), and --username (-u) are mandatory. The password can be included via the flag --password (-p) or entered at the prompt.
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

