## apictl mg delete api

Delete an API in Microgateway

### Synopsis

Delete an API by specifying name, version, host, username 
and optionally vhost by specifying the flags (--name (-n), --version (-v), --host (-c), 
--username (-u), and optionally --vhost (-t). Note: The password can be included 
via the flag --password (-p) or entered at the prompt.

```
apictl mg delete api [flags]
```

### Examples

```
apictl mg api--host https://localhost:9095 -u admin
  apictl mg api -n petstore -v 0.0.1 -c https://localhost:9095 -u admin -t www.pets.com 
  apictl mg api -n "petstore VIP" -v 0.0.1 --host https://localhost:9095 -u admin -p admin
```

### Options

```
  -h, --help              help for api
  -c, --host string       The adapter host url with port
  -n, --name string       API name
  -p, --password string   Password of the user
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

