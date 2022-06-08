## apictl import policy rate-limiting

Import Throttling Policy

### Synopsis

Import a Throttling Policy to an environment

```
apictl import policy rate-limiting --file <path-to-api> --environment <environment> [flags]
```

### Examples

```
apictl import rate-limiting -f qa/customadvanced -e dev
apictl import rate-limiting -f Env1/Exported/sub1 -e production
apictl import rate-limiting -f ~/CustomPolicy -e production -u
apictl import rate-limiting -f ~/mythottlepolicy -e production --update
NOTE: Both the flags (--file (-f) and --environment (-e)) are mandatory
```

### Options

```
  -e, --environment string   Environment from the which the Throttling Policy should be imported
  -f, --file string          File path of the Throttling Policy to be imported
  -h, --help                 help for rate-limiting
  -u, --update               Update an existing Throttling Policy or create a new Throttling Policy
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl import policy](apictl_import_policy.md)	 - Import a Policy

