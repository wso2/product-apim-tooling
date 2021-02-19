## apictl mg deploy

Deploy an API (apictl project) in Microgateway

### Synopsis

Deploy an API (apictl project) in Microgateway by specifying the adapter host url.

```
apictl mg deploy [flags]
```

### Examples

```
apictl mg deploy api -c https://localhost:9095 -f petstore -u admin -p admin

Note: The flags --host (-c), and --username (-u) are mandatory. The password can be included via the flag --password (-p) or entered at the prompt.
```

### Options

```
  -h, --help   help for deploy
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl mg](apictl_mg.md)	 - Handle Microgateway related operations
* [apictl mg deploy api](apictl_mg_deploy_api.md)	 - Deploy an API (apictl project) in Microgateway

