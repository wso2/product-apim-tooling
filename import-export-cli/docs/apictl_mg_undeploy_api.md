## apictl mg undeploy api

Undeploy an API in Microgateway

### Synopsis

Undeploy an API in Microgateway by specifying name, version, environment, username and optionally vhost

```
apictl mg undeploy api [flags]
```

### Examples

```
apictl mg undeploy api --environment dev -n petstore -v 0.0.1
   apictl mg undeploy api -n petstore -v 0.0.1 -e dev --vhost www.pets.com 
   apictl mg undeploy api -n SwaggerPetstore -v 0.0.1 --environment dev

Note: The flags --name (-n), --version (-v), --environment (-e) are mandatory. The user needs to be logged in to use this command.
```

### Options

```
  -e, --environment string   Microgateway adapter environment to be undeployed from
  -h, --help                 help for api
  -n, --name string          API name
  -v, --version string       API version
  -t, --vhost string         Virtual host the API needs to be undeployed from
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl mg undeploy](apictl_mg_undeploy.md)	 - Undeploy an API in Microgateway

