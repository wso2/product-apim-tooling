## apictl add

Add Environment to Config file

### Synopsis

Add new environment and its related endpoints to the config file

### Examples

```
apictl add env production \
--apim  https://localhost:9443 

apictl add env dev \
--mi  https://localhost:9164

apictl add env prod \
--apim  https://apim.com:9443 \
--mi https://localhost:9164

apictl add env test \
--registration https://idp.com:9443 \
--publisher https://apim.com:9443 \
--devportal  https://apps.com:9443 \
--admin  https://apim.com:9443 \
--token https://gw.com:8243/token \
--mi https://localhost:9164

apictl add env dev \
--apim https://apim.com:9443 \
--registration https://idp.com:9443 \
--token https://gw.com:8243/token

NOTE: The flag --environment (-e) is mandatory.
You can either provide only the flag --apim , or all the other 4 flags (--registration --publisher --devportal --admin) without providing --apim flag.
If you are omitting any of --registration --publisher --devportal --admin flags, you need to specify --apim flag with the API Manager endpoint. In both of the
cases --token flag is optional and use it to specify the gateway token endpoint. This will be used for "apictl get-keys" operation.
To add a micro integrator instance to an environment you can use the --mi flag.
```

### Options

```
  -h, --help   help for add
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl](apictl.md)	 - CLI for Importing and Exporting APIs and Applications and Managing WSO2 Micro Integrator
* [apictl add env](apictl_add_env.md)	 - Add Environment to Config file

