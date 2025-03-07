## apictl add env

Add Environment to Config file

### Synopsis

Add new environment and its related endpoints to the config file

```
apictl add env [environment] [flags]
```

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
--token https://gw.com:9443/oauth2/token \
--mi https://localhost:9164

apictl add env prod \
--apim  https://apim.com:9443 \
--ai-service https://dev-tools.wso2.com/apim-ai-service/v2 \
--ai-token-service https://api.asgardeo.io/t/wso2devtools/oauth2/token

apictl add env test \
--registration https://idp.com:9443 \
--publisher https://apim.com:9443 \
--devportal  https://apps.com:9443 \
--admin  https://apim.com:9443 \
--ai-service https://dev-tools.wso2.com/apim-ai-service/v2 \
--ai-token-service https://api.asgardeo.io/t/wso2devtools/oauth2/token

apictl add env dev \
--apim https://apim.com:9443 \
--registration https://idp.com:9443 \
--token https://gw.com:9443/oauth2/token

You can either provide only the flag --apim , or all the other 4 flags (--registration --publisher --devportal --admin) without providing --apim flag.
If you are omitting any of --registration --publisher --devportal --admin flags, you need to specify --apim flag with the API Manager endpoint. In both of the
cases --token flag is optional and use it to specify the gateway token endpoint. This will be used for "apictl get-keys" operation.
To add a micro integrator instance to an environment you can use the --mi flag.
To add ai related service endpoints to an environment you can use the --ai-service and --ai-token-service flags.
```

### Options

```
      --admin string              Admin endpoint for the environment
      --ai-service string         AI service endpoint for the environment
      --ai-token-service string   AI token service endpoint for the environment
      --apim string               API Manager endpoint for the environment
      --devportal string          DevPortal endpoint for the environment
  -h, --help                      help for env
      --mi string                 Micro Integrator Management endpoint for the environment
      --publisher string          Publisher endpoint for the environment
      --registration string       Registration endpoint for the environment
      --token string              Token endpoint for the environment
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl add](apictl_add.md)	 - Add Environment to Config file

