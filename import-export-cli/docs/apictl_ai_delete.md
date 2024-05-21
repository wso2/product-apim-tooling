## apictl ai delete

Purge APIs and API Products of a tenant from one environment from the vector database.

### Synopsis

Purge APIs and API Products of a tenant from one environment specified by flag (--environment, -e)

```
apictl ai delete [flags]
```

### Examples

```
apictl ai delete artifacts --token 2fdca1b6-6a28-4aea-add6-77c97033bdb9 --endpoint https://dev-tools.wso2.com/apim-ai-service -e production 
NOTE:The flag (--environment (-e)) is mandatory
```

### Options

```
  -h, --help   help for delete
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl ai](apictl_ai.md)	 - AI related commands.
* [apictl ai delete artifacts](apictl_ai_delete_artifacts.md)	 - Purge APIs and API Products of a tenant from one environment from a vector database.

