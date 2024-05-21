## apictl ai

AI related commands.

### Synopsis

Perform AI related operations such as uploading APIs and API Products of a tenant from one environment to a vector database to provide context to the marketplace assistant.

```
apictl ai [flags]
```

### Examples

```
apictl ai upload apis --token 2fdca1b6-6a28-4aea-add6-77c97033bdb9 --endpoint https://dev-tools.wso2.com/apim-ai-service -e production
apictl ai upload api-products --token 2fdca1b6-6a28-4aea-add6-77c97033bdb9 --endpoint https://dev-tools.wso2.com/apim-ai-service -e production 
apictl ai upload apis --token 2fdca1b6-6a28-4aea-add6-77c97033bdb9 --endpoint https://dev-tools.wso2.com/apim-ai-service -e production --all
apictl ai delete apis --token 2fdca1b6-6a28-4aea-add6-77c97033bdb9 -e -endpoint https://dev-tools.wso2.com/apim-ai-service production --all
NOTE:The flag (--environment (-e)) is mandatory
```

### Options

```
  -h, --help   help for ai
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl](apictl.md)	 - CLI for Importing and Exporting APIs and Applications and Managing WSO2 Micro Integrator
* [apictl ai delete](apictl_ai_delete.md)	 - Purge APIs and API Products of a tenant from one environment from the vector database.
* [apictl ai upload](apictl_ai_upload.md)	 - Upload APIs and API Products of a tenant from one environment to a vector database to provide context to the marketplace assistant.

