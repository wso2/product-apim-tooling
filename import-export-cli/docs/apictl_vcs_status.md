## apictl vcs status

Shows the list of projects that are ready to deploy

### Synopsis

Shows the list of projects that are ready to deploy to the specified environment by --environment(-e)
NOTE: --environment (-e) flag is mandatory

```
apictl vcs status [flags]
```

### Examples

```
apictl vcs status -e dev
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

* [apictl vcs](apictl_vcs.md)	 - Checks status and deploys projects

