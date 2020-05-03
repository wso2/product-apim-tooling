## apictl change-api-status

Change API Status

### Synopsis

Change the lifecycle status of an API in an environment

```
apictl change-api-status (--action <action-of-the-api-state-change> --name <name-of-the-api> --version <version-of-the-api> --provider <provider-of-the-api> --environment <environment-from-which-the-api-state-should-be-changed>) [flags]
```

### Examples

```
apictl change-api-status -a Publish -n TwitterAPI -v 1.0.0 -r admin -e dev
apictl change-api-status -a Publish -n FacebookAPI -v 2.1.0 -e production
NOTE: The 4 flags (--action (-a), --name (-n), --version (-v), and --environment (-e)) are mandatory.
```

### Options

```
  -a, --action string        Action to be taken to change the status of the API
  -e, --environment string   Environment of which the API state should be changed
  -h, --help                 help for change-api-status
  -n, --name string          Name of the API to be state changed
  -r, --provider string      Provider of the API
  -v, --version string       Version of the API to be state changed
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl](apictl.md)	 - CLI for Importing and Exporting APIs and Applications

