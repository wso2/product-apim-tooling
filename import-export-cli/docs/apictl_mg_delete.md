## apictl mg delete

Delete an API in Microgateway

### Synopsis

Delete an API in Microgateway by specifying name, version, host, username and optionally vhost

```
apictl mg delete [flags]
```

### Examples

```
apictl mg api -n petstore -v 0.0.1 --host https://localhost:9095 -u admin -t www.pets.com

Note: The flags --name (-n), --version (-v), --host (-c), and --username (-u) are mandatory. The password can be included via the flag --password (-p) or entered at the prompt.
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

