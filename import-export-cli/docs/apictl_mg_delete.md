## apictl mg delete

Delete an API in Microgateway

### Synopsis

Delete an API by specifying name, version, host, username 
and optionally vhost by specifying the flags (--name (-n), --version (-v), --host (-c), 
--username (-u), and optionally --vhost (-t). Note: The password can be included 
via the flag --password (-p) or entered at the prompt.

```
apictl mg delete [flags]
```

### Examples

```
apictl mg api -n petstore -v 0.0.1 --host https://localhost:9095 -u admin -t www.pets.com
```

### Options

```
  -h, --help   help for delete
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl mg](apictl_mg.md)	 - Handle Microgateway related operations
* [apictl mg delete api](apictl_mg_delete_api.md)	 - Delete an API in Microgateway

