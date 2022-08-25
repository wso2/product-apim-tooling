## apictl import policy api

Import an API Policy

### Synopsis

Import an API Policy to an environment

```
apictl import policy api --file <path-to-api-policy> --environment <environment> [flags]
```

### Examples

```
apictl import policy api -f add_header_v1.zip -e dev
apictl import policy api -f AddHeader -e production
NOTE: Both the flags (--file (-f) and --environment (-e)) are mandatory
```

### Options

```
  -e, --environment string   Environment from the which the API Policy should be imported
  -f, --file string          File path of the API Policy to be imported
  -h, --help                 help for api policy
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl import policy](apictl_import_policy.md)	 - Import a Policy

