## apictl delete policy rate-limiting

Delete Throttling Policy

### Synopsis

Delete a throttling policy from an environment

```
apictl delete policy rate-limiting (--name <name-of-the-throttling-policy> --environment <environment-from-which-the-policy-should-be-deleted>)--type <type-of-the-throttling-policy> [flags]
```

### Examples

```
apictl delete policy rate-limiting -n Gold -e dev --type sub 
apictl delete policy rate-limiting -n AppPolicy -e prod --type app
apictl delete policy rate-limiting -n TestPolicy -e dev --type advanced 
apictl delete policy rate-limiting -n CustomPolicy -e prod --type custom 
NOTE: All the 2 flags (--name (-n) and --environment (-e)) are mandatory.
```

### Options

```
  -e, --environment string   Environment from which the Throttling Policy should be deleted
  -h, --help                 help for rate-limiting
  -n, --name string          Name of the Throttling Policy to be deleted
  -t, --type string          Type of the Throttling Policies to be exported (sub,app,custom,advanced)
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl delete policy](apictl_delete_policy.md)	 - Delete a Policy

