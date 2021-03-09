## apictl mg get

List APIs in Microgateway

### Synopsis

Display a list of all the APIs in a Microgateway Adapter environment or a set of APIs with a limit set or filtered by apiType

```
apictl mg get [flags]
```

### Examples

```
apictl mg get apis --environment dev
apictl mg get apis -q type:http --environment dev -l 100
apictl mg get apis -q type:ws --environment dev

Note: The flags --environment (-e) is mandatory. 
The user needs to be logged in to use this command.
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

