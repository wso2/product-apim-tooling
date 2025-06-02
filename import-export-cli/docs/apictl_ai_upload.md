## apictl ai upload

Upload APIs and API Products of a tenant from one environment to a vector database to provide context to the marketplace assistant.

### Synopsis

Upload APIs and API Products of a tenant from one environment specified by flag (--environment, -e)

```
apictl ai upload [flags]
```

### Examples

```
apictl ai upload apis -e production --all
apictl ai upload api-products -e production --all
NOTE:The flag (--environment (-e)) is mandatory
```

### Options

```
  -h, --help   help for upload
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl ai](apictl_ai.md)	 - AI related commands.
* [apictl ai upload api-products](apictl_ai_upload_api-products.md)	 - Upload API Products of a tenant from one environment to a vector database.
* [apictl ai upload apis](apictl_ai_upload_apis.md)	 - Upload APIs of a tenant from one environment to a vector database.

