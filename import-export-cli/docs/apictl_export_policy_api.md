## apictl export policy api

Export an API Policy

### Synopsis

Export an API Policy from an environment

```
apictl export policy api (--name <name-of-the-api-policy> --environment <environment-from-which-the-api-policy-should-be-exported>) [flags]
```

### Examples

```
apictl export policy api -n addHeader -e dev
NOTE: All the 2 flags (--name (-n) and --environment (-e)) are mandatory.
```

### Options

```
  -e, --environment string   Environment of the API Policy to be exported
      --format string        Type of the Policy definition file exported
  -h, --help                 help for api
  -n, --name string          Name of the API Policy to be exported
  -v, --version string       Version of the API Policy to be exported
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl export policy](apictl_export_policy.md)	 - Export/Import a Policy

