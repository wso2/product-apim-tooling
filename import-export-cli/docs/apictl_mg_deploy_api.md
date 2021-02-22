## apictl mg deploy api

Deploy an API (apictl project) in Microgateway

### Synopsis

Deploy an API (apictl project) in Microgateway by specifying the adapter host url.

```
apictl mg deploy api [flags]
```

### Examples

```
apictl mg deploy api -c https://localhost:9095 -f petstore -u admin -p admin

Note: The flags --host (-c), and --deployAPIUsername (-u) are mandatory. The password can be included via the flag --password (-p) or entered at the prompt.
```

### Options

```
  -f, --file string       Filepath of the apictl project to be deployed
  -h, --help              help for api
  -c, --host string       Host url for the control plane with port
  -o, --override          Whether to deploy an API irrespective of its existance. Overrides when exists.
  -p, --password string   Password of the user (Can be provided at the prompt)
      --skip-cleanup      Whether to keep all temporary files created during deploy process
  -u, --username string   Username of the user
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl mg deploy](apictl_mg_deploy.md)	 - Deploy an API (apictl project) in Microgateway

