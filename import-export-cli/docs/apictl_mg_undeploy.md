## apictl mg undeploy

Undeploy an API in Microgateway

### Synopsis

Undeploy an API in Microgateway by specifying name, version, environment, username and optionally vhost

```
apictl mg undeploy [flags]
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
  -h, --help   help for undeploy
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl mg](apictl_mg.md)	 - Handle Microgateway related operations
* [apictl mg undeploy api](apictl_mg_undeploy_api.md)	 - Undeploy an API in Microgateway

