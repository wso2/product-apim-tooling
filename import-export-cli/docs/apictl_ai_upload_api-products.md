## apictl ai upload api-products

Upload API Products of a tenant from one environment to a vector database.

### Synopsis

Upload public API Products of a tenant from one environment specified by flag (--environment, -e)

```
apictl ai upload api-products (--endpoint <endpoint-url> --token <on-prem-key-of-the-organization> --environment <environment-from-which-artifacts-should-be-uploaded>) [flags]
```

### Examples

```
apictl upload api-products -e production --all
apictl upload api-products -e production
apictl upload api-products -e production
NOTE:The flag (--environment (-e)) is mandatory
```

### Options

```
      --all                  Upload both apis and api products
      --endpoint string      endpoint of the marketplace assistant service
  -e, --environment string   Environment from which the APIs should be uploaded
  -h, --help                 help for api-products
      --token string         on-prem-key of the organization
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl ai upload](apictl_ai_upload.md)	 - Upload APIs and API Products of a tenant from one environment to a vector database to provide context to the marketplace assistant.

