## apictl ai upload apis

Upload APIs of a tenant from one environment to a vector database.

### Synopsis

Upload APIs of a tenant from one environment to a vector database to provide context to the marketplace assistant.

```
apictl ai upload apis (--environment <environment-from-which-artifacts-should-be-uploaded> --all) [flags]
```

### Examples

```
apictl upload apis -e production --all
apictl upload apis -e production
apictl upload apis -e production
NOTE:The flag (--environment (-e)) is mandatory
```

### Options

```
      --all                  Upload both apis and api products
  -e, --environment string   Environment from which the APIs should be uploaded
  -h, --help                 help for apis
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl ai upload](apictl_ai_upload.md)	 - Upload APIs and API Products of a tenant from one environment to a vector database to provide context to the marketplace assistant.

