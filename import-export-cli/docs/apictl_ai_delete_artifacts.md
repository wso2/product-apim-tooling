## apictl ai delete artifacts

Purge APIs and API Products of a tenant from one environment from a vector database.

### Synopsis

Purge APIs and API Products of a tenant from one environment from a vector database.

```
apictl ai delete artifacts (--endpoint <endpoint-url> --token <on-prem-key-of-the-organization> --environment <environment-from-which-artifacts-should-be-purged>) [flags]
```

### Examples

```
apictl ai delete artifacts -e production
NOTE:The flag (--environment (-e)) is mandatory
```

### Options

```
      --endpoint string      endpoint of the marketplace assistant service
  -e, --environment string   Environment from which the APIs should be Purged
  -h, --help                 help for artifacts
      --token string         on-prem-key of the organization
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl ai delete](apictl_ai_delete.md)	 - Purge APIs and API Products of a tenant from one environment from the vector database.

