## apictl mg deploy

Deploy an API (apictl project) in Microgateway

### Synopsis

Deploy an API (apictl project) in Microgateway by specifying the microgateway adapter environment.

```
apictl mg deploy [flags]
```

### Examples

```
apictl mg deploy api -e dev -f petstore

Note: The flags --environment (-e), --file (-f) are mandatory. The user needs to be logged in to use this command.
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

