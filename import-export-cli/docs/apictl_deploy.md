## apictl deploy

Deploy an API/API Product in a gateway environment

### Synopsis

Deploy an API/API Product available in the environment specified by flag (--environment, -e)
to the gateway specified by flag (--gateway, -g)

```
apictl deploy [flags]
```

### Examples

```
apictl deploy api -n TwitterAPI -v 1.0.0 -r admin --rev 1 -g Label1 -e dev
apictl deploy api -n FacebookAPI -v 2.1.0 --rev 6 -g Label1 Label2 Label3 -e production
apictl export api-product -n FacebookAPI -v 2.1.0 --rev 2 -r admin -g Label1 -e production --hide-on-devportal
```

### Options

```
  -h, --help   help for deploy
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl](apictl.md)	 - CLI for Importing and Exporting APIs and Applications and Managing WSO2 Micro Integrator
* [apictl deploy api](apictl_deploy_api.md)	 - Deploy API

