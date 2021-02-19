## apictl mg undeploy

Undeploy an API in Microgateway

### Synopsis

Undeploy an API in Microgateway by specifying name, version, host, username and optionally vhost

```
apictl mg undeploy [flags]
```

### Examples

```
apictl mg undeploy api --host https://localhost:9095 -n petstore -v 0.0.1 -u admin
apictl mg undeploy api -n petstore -v 0.0.1 -c https://localhost:9095 -u admin --vhost www.pets.com 
apictl mg undeploy api -n SwaggerPetstore -v 0.0.1 --host https://localhost:9095 -u admin -p admin

Note: The flags --name (-n), --version (-v), --host (-c), and --username (-u) are mandatory. The password can be included via the flag --password (-p) or entered at the prompt.
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

