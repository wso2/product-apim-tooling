## apictl undeploy api

Undeploy API

### Synopsis

Undeploy an API revision from gateway environments

```
apictl undeploy api (--name <name-of-the-api> --version <version-of-the-api> --provider <provider-of-the-api> --rev <revision-number-of-the-api> --gateway-env <gateway-environment> --environment <environment-from-which-the-api-should-be-undeployed>) [flags]
```

### Examples

```
apictl undeploy api -n TwitterAPI -v 1.0.0 -rev 2 -e dev
apictl undeploy api -n FacebookAPI -v 2.1.0 --rev 6 -g Label1 Label2 Label3 -e production
apictl undeploy api -n FacebookAPI -v 2.1.0 -r alice --rev 2 -g Label1 -e production
NOTE: All the 4 flags (--name (-n), --version (-v), --rev, --environment (-e)) are mandatory. 
If the flag (--gateway-env (-g)) is not provided, revision will be undeployed from all deployed gateway environments.
```

### Options

```
  -e, --environment string   Environment of which the API should be undeployed
  -g, --gateway-env string   Gateway environment which the revision has to be undeployed
  -h, --help                 help for api
  -n, --name string          Name of the API to be exported
  -r, --provider string      Provider of the API
      --rev string           Revision number of the API to undeploy
  -v, --version string       Version of the API to be exported
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl undeploy](apictl_undeploy.md)	 - Undeploy an API/API Product revision from a gateway environment

