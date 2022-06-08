## apictl export policy

Export/Import a Policy

### Synopsis

Export/Import a Policy in an environment or Import a Policy to an environment

```
apictl export policy [flags]
```

### Examples

```
apictl export policy rate-limiting -n Silver -e prod --type subscription
```

### Options

```
  -h, --help   help for policy
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl export](apictl_export.md)	 - Export an API/API Product/Application/Policy in an environment
* [apictl export policy rate-limiting](apictl_export_policy_rate-limiting.md)	 - Export Throttling Policies

