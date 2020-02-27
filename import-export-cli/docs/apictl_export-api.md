## apictl export-api

Export API

### Synopsis

Export APIs from an environment

```
apictl export-api (--name <name-of-the-api> --version <version-of-the-api> --environment <environment-from-which-the-api-should-be-exported>) [flags]
```

### Examples

```
apictl export-api -n TwitterAPI -v 1.0.0 -e dev --provider admin
apictl export-api -n FacebookAPI -v 2.1.0 -e production --provider admin
NOTE: all three flags (--name (-n), --version (-v), --provider (-r)) are mandatory
```

### Options

```
  -e, --environment string   Environment to which the API should be exported
      --format string        File format of exported archive(json or yaml)
  -h, --help                 help for export-api
  -n, --name string          Name of the API to be exported
      --preserveStatus       Preserve API status when exporting. Otherwise API will be exported in CREATED status (default true)
  -r, --provider string      Provider of the API
  -v, --version string       Version of the API to be exported
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl](apictl.md)	 - CLI for Importing and Exporting APIs and Applications

