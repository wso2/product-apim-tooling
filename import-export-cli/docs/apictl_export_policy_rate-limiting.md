## apictl export policy rate-limiting

Export Throttling Policies

### Synopsis

Export ThrottlingPolicies from an environment

```
apictl export policy rate-limiting (--type <type-of-the-throttling-policy> --environment <environment-from-which-the-throttling-policies-should-be-exported>) [flags]
```

### Examples

```
apictl export policy rate-limiting -n Gold -e dev --type sub 
apictl export policy rate-limiting -n AppPolicy -e prod --type app --format JSON
apictl export policy rate-limiting -n TestPolicy -e dev --type advanced 
apictl export policy rate-limiting -n CustomPolicy -e prod --type custom 
NOTE: All the 2 flags (--name (-n) and --environment (-e)) are mandatory.
```

### Options

```
  -e, --environment string   Environment to which the Throttling Policies should be exported
      --format string        File format of exported archive(JSON or YAML) (default "YAML")
  -h, --help                 help for rate-limiting
  -n, --name string          Name of the Throttling Policy to be exported
  -t, --type string          Type of the Throttling Policies to be exported (sub,app,custom,advanced)
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl export policy](apictl_export_policy.md)	 - Export/Import a Policy

