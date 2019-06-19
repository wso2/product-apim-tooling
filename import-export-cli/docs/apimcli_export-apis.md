## apimcli export-apis

Export APIs for migration

### Synopsis


Export all the APIs of a tenant from an APIM 2.6.0 environment environment, to be imported into 3.0.0 environment

```
apimcli export-apis [--environment <environment-from-which-artifacts-should-be-exported>] -u <user_name> -p <password> [-t <Tenant-domain-of-the-resources-to-be-exported>] [--force] [flags]
```

### Examples

```
apimcli export-apis -e production-2.6.0 -u wso2admin@wso2.org -p 12345 -t wso2.org -k --force
		apimcli export-apis -e production-2.6.0 -u admin -p admin -k
```

### Options

```
  -e, --environment string   Environment to which the API should be exported
      --force                Clean all the previously exported APIs of the given target tenant, in the given environment if any, and to export APIs from beginning
  -h, --help                 help for export-apis
  -t, --tenant string        Tenant domain of the resources to be exported
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO
* [apimcli](apimcli.md)	 - CLI for Importing and Exporting APIs and Applications

