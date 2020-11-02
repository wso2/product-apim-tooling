## apictl k8s delete

Delete resources related to kubernetes

### Synopsis

Delete resources by filenames, stdin, resources and names, or by resources and label selector in kubernetes mode

```
apictl k8s delete [flags]
```

### Examples

```
apictl delete api petstore
apictl delete api -l name=myLabel
```

### Options

```
  -h, --help   help for delete
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl k8s](apictl_k8s.md)	 - Kubernetes mode based commands
* [apictl k8s delete apictl](apictl_k8s_delete_apictl.md)	 - Delete API resources

