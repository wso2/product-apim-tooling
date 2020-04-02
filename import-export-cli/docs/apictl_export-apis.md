## apictl export-apis

Export APIs for migration

### Synopsis

Export all the APIs of a tenant from one environment, to be imported into another environment

```
apictl export-apis (--environment <environment-from-which-artifacts-should-be-exported> -t <Tenant-domain-of-the-resources-to-be-exported> --force) [flags]
```

### Examples

```
apictl export-apis -e production -t wso2.org --force
apictl export-apis -e production
NOTE: The flag (--environment (-e)) is mandatory
```

### Options

```
  -e, --environment string   Environment to which the API should be exported
      --force                Clean all the previously exported APIs of the given target tenant, in the given environment if any, and to export APIs from beginning
      --format string        File format of exported archives(json or yaml)
  -h, --help                 help for export-apis
      --preserveStatus       Preserve API status when exporting. Otherwise API will be exported in CREATED status (default true)
  -t, --tenant string        Tenant domain of the resources to be exported
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl](apictl.md)	 - CLI for Importing and Exporting APIs and Applications

