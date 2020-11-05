## apictl vcs deploy

Deploys projects to the specified environment

### Synopsis

Deploys projects to the specified environment specified by --environment(-e). 
Only the changed projects compared to the revision at the last successful deployment will be deployed. 
If any project(s) got failed during the deployment, by default, the operation will rollback the environment to the last successful state. If this needs to be avoided, use --skipRollback=true
NOTE: --environment (-e) flag is mandatory

```
apictl vcs deploy [flags]
```

### Examples

```
apictl vcs deploy -e dev
apictl vcs deploy -e dev --skipRollback=true
```

### Options

```
  -e, --environment string   Name of the environment to deploy the project(s)
  -h, --help                 help for deploy
      --skipRollback         Specifies whether rolling back to the last successful revision during an error situation should be skipped
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl vcs](apictl_vcs.md)	 - Checks status and deploys projects

