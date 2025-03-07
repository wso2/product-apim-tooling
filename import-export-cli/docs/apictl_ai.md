## apictl ai

AI related commands.

### Synopsis

Perform AI related operations such as uploading APIs and API Products of a tenant from one environment to a vector database to provide context to the marketplace assistant.

```
apictl ai [flags]
```

### Examples

```
apictl ai upload apis --key Zk9DaTR2Tko1OVBwSHVjQzJDQVlmWXVBRGRNYTphNEZ3SGxxMGlDSUtWczJNUElJRG5lcFpuWU1h -e production
apictl ai upload api-products --key Zk9DaTR2Tko1OVBwSHVjQzJDQVlmWXVBRGRNYTphNEZ3SGxxMGlDSUtWczJNUElJRG5lcFpuWU1h -e production
apictl ai upload apis --key Zk9DaTR2Tko1OVBwSHVjQzJDQVlmWXVBRGRNYTphNEZ3SGxxMGlDSUtWczJNUElJRG5lcFpuWU1h -e production --all
apictl ai delete apis --key Zk9DaTR2Tko1OVBwSHVjQzJDQVlmWXVBRGRNYTphNEZ3SGxxMGlDSUtWczJNUElJRG5lcFpuWU1h -e production --all
NOTE:The flags (--key and --environment (-e)) are mandatory
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

