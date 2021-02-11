## apictl deploy api

Deploy API

### Synopsis

Deploy an API to the given gateway environment

```
apictl deploy api (--name <name-of-the-api> --version <version-of-the-api> --provider <provider-of-the-api> --rev <revision_number> --gateway <gateway-environment> --environment <environment-from-which-the-api-should-be-deployed>) [flags]
```

### Examples

```
apictl deploy api -n TwitterAPI -v 1.0.0 -r admin --rev 1 -g Label1 -e dev
apictl deploy api -n FacebookAPI -v 2.1.0 --rev 6 -g Label1 Label2 Label3 -e production
apictl deploy api -n FacebookAPI -v 2.1.0 --rev 2 -r admin -g Label1 -e production --hide-on-devportal
NOTE: All the 5 flags (--name (-n), --version (-v) , --rev and --gateway (-g) --environment (-e)) are mandatory.
```

### Options

```
  -e, --environment string   Environment to which the API should be deployed
  -g, --gateway string       Gateways which the revision has to be deployed
  -h, --help                 help for api
      --hide-on-devportal    Hide the gateway environment on devportal
  -n, --name string          Name of the API to be deployed
  -r, --provider string      Provider of the API
      --rev string           Revision number of the API to be deployed
  -v, --version string       Version of the API to be deployed
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl deploy](apictl_deploy.md)	 - Deploy an API/API Product in a gateway environment

