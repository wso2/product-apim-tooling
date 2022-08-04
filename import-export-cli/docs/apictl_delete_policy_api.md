## apictl delete policy api

Delete API Policy

### Synopsis

Export API Policy from an environment

```
apictl delete policy api (--name <name-of-the-api-policy> --environment <environment-from-which-the-api-policy-should-be-deleted>) [flags]
```

### Examples

```
apictl delete policy api -n addHeader -e dev
NOTE: All the 2 flags (--name (-n) and --environment (-e)) are mandatory.
```

### Options

```
  -e, --environment string   Environment from which the API Policy should be deleted
  -h, --help                 help for API Policy
  -n, --name string          Name of the API Policy to be deleted
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl delete policy](apictl_delete_policy.md) - Delete a Policy

