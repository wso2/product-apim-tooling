## apictl delete policy api

Delete an API Policy

### Synopsis

Delete an API Policy from an environment

```
apictl delete policy api (--name <name-of-the-api-policy> --environment <environment-from-which-the-policy-should-be-deleted>) [flags]
```

### Examples

```
apictl delete policy api -n addHeader -e dev
 NOTE: The 2 flags (--name (-n) and --environment (-e)) are mandatory.
```

### Options

```
  -e, --environment string   Environment from which the API Policy should be deleted
  -h, --help                 help for api
  -n, --name string          Name of the API Policy to be deleted
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl delete policy](apictl_delete_policy.md)	 - Delete a Policy

