## apictl k8s delete apictl

Delete API resources

### Synopsis

Delete API resources by API name or label selector in kubernetes mode

```
apictl k8s delete apictl delete api (<name-of-the-api> or -l name=<name-of-the-label>) [flags]
```

### Examples

```
apictl delete api petstore
  apictl delete api -l name=myLabel
```

### Options

```
  -h, --help   help for apictl
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl k8s delete](apictl_k8s_delete.md)	 - Delete resources related to kubernetes

