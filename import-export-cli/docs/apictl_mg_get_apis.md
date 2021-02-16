## apictl mg get apis

List APIs in Microgateway

### Synopsis

Display a list of all the APIs in Microgateway or a set of APIs with a limit set or filtered by apiType

```
apictl mg get apis [flags]
```

### Examples

```
apictl mg get apis --host https://localhost:9095 -u admin
apictl mg get apis -t http --host https://localhost:9095 -u admin -l 100
apictl mg get apis -t ws --host https://localhost:9095 -u admin

Note: The flags --host (-c), --username (-u) are mandatory. The password can be included via the flag --password (-p) or entered at the prompt.
```

### Options

```
  -h, --help              help for apis
  -c, --host string       The adapter host url with port
  -l, --limit string      Maximum number of APIs to return
  -p, --password string   Password of the user (Can be provided at the prompt)
  -t, --type string       API type to filter the APIs
  -u, --username string   The username
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl mg get](apictl_mg_get.md)	 - List APIs in Microgateway

