## apictl mg delete api

Delete an API in Microgateway

### Synopsis

Delete an API in Microgateway by specifying name, version, host, username and optionally vhost

```
apictl mg delete api [flags]
```

### Examples

```
apictl mg delete api --host https://localhost:9095 -n petstore -v 0.0.1 -u admin
  apictl mg delete api -n petstore -v 0.0.1 -c https://localhost:9095 -u admin --vhost www.pets.com 
  apictl mg delete api -n SwaggerPetstore -v 0.0.1 --host https://localhost:9095 -u admin -p admin

Note: The flags --name (-n), --version (-v), --host (-c), and --username (-u) are mandatory. The password can be included via the flag --password (-p) or entered at the prompt.
```

### Options

```
  -h, --help              help for api
  -c, --host string       The adapter host url with port
  -n, --name string       API name
  -p, --password string   Password of the user (Can be provided at the prompt)
  -u, --username string   Username with delete permissions
  -v, --version string    API version
  -t, --vhost string      Virtual host the API needs to be deleted from
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl mg delete](apictl_mg_delete.md)	 - Delete an API in Microgateway

