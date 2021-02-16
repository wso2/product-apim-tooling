## apictl mg get

List APIs in Microgateway

### Synopsis

Display a list of all the APIs in Microgateway or a set of APIs with a limit set or filtered by apiType

```
apictl mg get [flags]
```

### Examples

```
apictl mg get apis -t http --host https://localhost:9095 -u admin -l 100

	Note: The flags --host (-c), --username (-u) are mandatory. The password can be included via the flag --password (-p) or entered at the prompt.
```

### Options

```
  -h, --help   help for get
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl mg](apictl_mg.md)	 - Handle Microgateway related operations
* [apictl mg get apis](apictl_mg_get_apis.md)	 - List APIs in Microgateway

