## apictl delete api

Delete API

### Synopsis

Delete an API from an environment in default mode and delete API resources by API name or label selector in kubernetes mode

```
apictl delete api (--name <name-of-the-api> --version <version-of-the-api> --provider <provider-of-the-api> --environment <environment-from-which-the-api-should-be-deleted>) [Flags]
Kubernetes Mode:
  apictl delete api (<name-of-the-api> or -l name=<name-of-the-label>) [flags]
```

### Examples

```
Default Mode:
  apictl delete api -n TwitterAPI -v 1.0.0 -r admin -e dev
  apictl delete api -n FacebookAPI -v 2.1.0 -e production
NOTE: The 3 flags (--name (-n), --version (-v), and --environment (-e)) are mandatory.
Kubernetes Mode:
  apictl delete api petstore
  apictl delete api -l name=myLabel
```

### Options

```
  -e, --environment string   Environment from which the API should be deleted
  -h, --help                 help for api
  -n, --name string          Name of the API to be deleted
  -r, --provider string      Provider of the API to be deleted
  -v, --version string       Version of the API to be deleted
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl delete](apictl_delete.md)	 - Delete an API/APIProduct/Application in an environment

