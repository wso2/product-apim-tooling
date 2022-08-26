## apictl delete policy

Delete a Policy

### Synopsis

Delete a Policy in an environment

```
apictl delete policy [flags]
```

### Examples

```
apictl delete policy api -n addHeader -e prod
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

* [apictl delete](apictl_delete.md)	 - Delete an API/APIProduct/Application in an environment
* [apictl delete policy api](apictl_delete_policy_api.md)	 - Delete an API Policy
* [apictl delete policy rate-limiting](apictl_delete_policy_rate-limiting.md)	 - Delete Throttling Policy

