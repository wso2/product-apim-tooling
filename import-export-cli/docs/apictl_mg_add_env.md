## apictl mg add env

Add Environment to Config file

### Synopsis

Add new environment and its related endpoints to the config file

```
apictl mg add env [flags]
```

### Examples

```
apictl mg add env prod --host  https://localhost:9443 

NOTE: The flag --host (-c) is mandatory and it has to specify the microgateway adapter url.
```

### Options

```
  -h, --help          help for env
  -c, --host string   The adapter host url with port
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl mg add](apictl_mg_add.md)	 - Add Environment to Config file

