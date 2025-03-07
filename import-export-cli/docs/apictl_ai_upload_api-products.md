## apictl ai upload api-products

Upload API Products of a tenant from one environment to a vector database.

### Synopsis

Upload public API Products of a tenant from one environment specified by flag (--environment, -e)

```
apictl ai upload api-products (--key <base64-encoded-client_id-and-client_secret> --environment <environment-from-which-artifacts-should-be-uploaded>) [flags]
```

### Examples

```
apictl upload api-products --key Zk9DaTR2Tko1OVBwSHVjQzJDQVlmWXVBRGRNYTphNEZ3SGxxMGlDSUtWczJNUElJRG5lcFpuWU1h -e production --all
apictl upload api-products --key Zk9DaTR2Tko1OVBwSHVjQzJDQVlmWXVBRGRNYTphNEZ3SGxxMGlDSUtWczJNUElJRG5lcFpuWU1h -e production
apictl upload api-products --key Zk9DaTR2Tko1OVBwSHVjQzJDQVlmWXVBRGRNYTphNEZ3SGxxMGlDSUtWczJNUElJRG5lcFpuWU1h -e production
NOTE:The flags (--key and --environment (-e)) are mandatory
```

### Options

```
      --all                  Upload both apis and api products
  -e, --environment string   Environment from which the APIs should be uploaded
  -h, --help                 help for api-products
      --key string           Base64 encoded client_id and client_secret pair
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl ai upload](apictl_ai_upload.md)	 - Upload APIs and API Products of a tenant from one environment to a vector database to provide context to the marketplace assistant.

