## apictl mg deploy api

Deploy an API (apictl project) in Microgateway

### Synopsis

Deploy an API (apictl project) in Microgateway by specifying the microgateway adapter environment.

```
apictl mg deploy api [flags]
```

### Examples

```
apictl mg deploy api -e dev -f petstore

Note: The flags --environment (-e), --file (-f) are mandatory. The user needs to be logged in to use this command.
```

### Options

```
  -e, --environment string   Microgateway adapter environment to add the API
  -f, --file string          Filepath of the apictl project to be deployed
  -h, --help                 help for api
  -o, --override             Whether to deploy an API irrespective of its existance. Overrides when exists.
      --skip-cleanup         Whether to keep all temporary files created during deploy process
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl mg deploy](apictl_mg_deploy.md)	 - Deploy an API (apictl project) in Microgateway

