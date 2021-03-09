## apictl mg get apis

List APIs in Microgateway

### Synopsis

Display a list of all the APIs in a Microgateway Adapter environment or a set of APIs with a limit set or filtered by apiType

```
apictl mg get apis [flags]
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
  -e, --environment string   Microgateway adapter environment to list APIs from
  -h, --help                 help for apis
  -l, --limit string         Maximum number of APIs to return
  -q, --query string         Query to filter the APIs
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl mg get](apictl_mg_get.md)	 - List APIs in Microgateway

