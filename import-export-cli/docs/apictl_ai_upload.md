## apictl ai upload

Upload APIs and API Products of a tenant from one environment to a vector database to provide context to the marketplace assistant.

### Synopsis

Upload APIs and API Products of a tenant from one environment specified by flag (--environment, -e)

```
apictl ai upload [flags]
```

### Examples

```
apictl ai upload apis --key Zk9DaTR2Tko1OVBwSHVjQzJDQVlmWXVBRGRNYTphNEZ3SGxxMGlDSUtWczJNUElJRG5lcFpuWU1h -e production --all
apictl ai upload api-products --key Zk9DaTR2Tko1OVBwSHVjQzJDQVlmWXVBRGRNYTphNEZ3SGxxMGlDSUtWczJNUElJRG5lcFpuWU1h -e production --all
NOTE:The flags (--key and --environment (-e)) are mandatory
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

