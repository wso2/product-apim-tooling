## apictl vcs init

Initializes a GIT repository with API Controller

### Synopsis

Initializes a GIT repository with API Controller (apictl). Before start using a GIT repository 
for 'vcs' commands, the GIT repository should be initialized once via 'vcs init'. This will create a file 'vcs.yaml'
in the root location of the GIT repository, which is used by API Controller  to uniquely identify the GIT repository. 
'vcs.yaml' should be committed to the GIT repository.

```
apictl vcs init [flags]
```

### Examples

```
apictl vcs init
```

### Options

```
  -f, --force   Forcefully reinitialize and replace vcs.yaml if already exists in the repository root.
  -h, --help    help for init
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl vcs](apictl_vcs.md)	 - Checks status and deploys projects

