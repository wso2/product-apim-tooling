## apictl mg add

Add Environment to Config file

### Synopsis

Add new environment and its related endpoints to the config file

### Examples

```
apictl mg add env prod --host  https://localhost:9443

NOTE: The flag --host (-c) is mandatory and it has to specify the microgateway adapter url.
```

### Options

```
  -h, --help   help for add
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl mg](apictl_mg.md)	 - Handle Microgateway related operations
* [apictl mg add env](apictl_mg_add_env.md)	 - Add Environment to Config file

