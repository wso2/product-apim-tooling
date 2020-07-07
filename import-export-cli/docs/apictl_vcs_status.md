## apictl vcs status

Gets the status report of project changes of the specified environment

### Synopsis

Gets the status report of project changes of the specified environment

```
apictl vcs status [flags]
```

### Examples

```
apictl status  -e dev
```

### Options

```
  -e, --environment string   Name of the environment to check the project(s) status
  -h, --help                 help for status
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl vcs](apictl_vcs.md)	 - Update an projects in an environment by calling the version control system (git)

