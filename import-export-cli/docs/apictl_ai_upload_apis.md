## apictl ai upload apis

Upload APIs of a tenant from one environment to a vector database.

### Synopsis

Upload APIs of a tenant from one environment to a vector database to provide context to the marketplace assistant.

```
apictl ai upload apis (--endpoint <endpoint-url> --token <on-prem-key-of-the-organization> --environment <environment-from-which-artifacts-should-be-uploaded> --all) [flags]
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
      --endpoint string      endpoint of the marketplace assistant service
  -e, --environment string   Environment from which the APIs should be uploaded
  -h, --help                 help for apis
      --token string         on-prem-key of the organization
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl ai upload](apictl_ai_upload.md)	 - Upload APIs and API Products of a tenant from one environment to a vector database to provide context to the marketplace assistant.

