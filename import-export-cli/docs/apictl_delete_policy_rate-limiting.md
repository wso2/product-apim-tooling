## apictl delete policy rate-limiting

Delete Throttling Policy

### Synopsis

Export Throttling Policy from an environment

```
apictl delete policy rate-limiting (--name <name-of-the-rate-limiting-policy> --type <type-of-the-rate-limiting-policy> --environment <environment-from-which-the-rate-limiting-policy-should-be-deleted>) [flags]
```

### Examples

```
apictl delete policy rate-limiting -n addHeader --type advanced -e dev
NOTE: All the 2 flags (--name (-n), --type and --environment (-e)) are mandatory.
```

### Options

```
  -e, --environment string   Environment from which the Throttling Policy should be deleted
  -h, --help                 help for Throttling Policy
  -n, --name string          Name of the Throttling Policy to be deleted
  --type string              Type of the Throttling Policy to be deleted
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl delete policy](apictl_delete_policy.md) - Delete a Policy

