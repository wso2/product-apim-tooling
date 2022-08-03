## apictl delete policy

Delete Policy

### Synopsis

Delete a Policy from an environment

```
apictl delete policy [policy type] (--environment <environment-from-which-the-application-should-be-deleted>) [flags]
```

### Examples

```
apictl delete policy api -n addHeader -e dev
NOTE: Both the flags (--name (-n), and --environment (-e)) are mandatory.
```

### Options

```
  -e, --environment string   Environment from which the Policy should be deleted
  -h, --help                 help for policy
  -n, --name string          Name of the Policy to be deleted
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl delete](apictl_delete.md)	 - Delete an API/APIProduct/Application/Policy in an environment
* [apictl delete policy api](apictl_delete_policy_api.md) - Delete an API Policy in an environment

