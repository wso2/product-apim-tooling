## apictl vcs rollback

Rollback the environment to the last working state in case of an error

### Synopsis

Rollback the environment to the last working state in case of an error

```
apictl vcs rollback [flags]
```

### Examples

```
apictl rollback  -e dev
```

### Options

```
  -e, --environment string   Name of the environment to check the project(s) status
  -h, --help                 help for rollback
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl vcs](apictl_vcs.md)	 - Checks status and deploys projects

