## apictl ai delete

Purge APIs and API Products of a tenant from one environment from the vector database.

### Synopsis

Purge APIs and API Products of a tenant from one environment specified by flag (--environment, -e)

```
apictl ai delete [flags]
```

### Examples

```
apictl ai delete artifacts --key Zk9DaTR2Tko1OVBwSHVjQzJDQVlmWXVBRGRNYTphNEZ3SGxxMGlDSUtWczJNUElJRG5lcFpuWU1h -e production
NOTE:The flags (--key and --environment (-e)) are mandatory
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

