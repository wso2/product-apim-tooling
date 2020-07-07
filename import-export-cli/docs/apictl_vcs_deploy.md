## apictl vcs deploy

Deploys project changes to the specified environment

### Synopsis

Deploys project changes to the specified environment

```
apictl vcs deploy [flags]
```

### Examples

```
apictl deploy  -e dev
```

### Options

```
  -e, --environment string   Name of the environment to deploy the project(s)
  -h, --help                 help for deploy
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl vcs](apictl_vcs.md)	 - Update an projects in an environment by calling the version control system (git)

