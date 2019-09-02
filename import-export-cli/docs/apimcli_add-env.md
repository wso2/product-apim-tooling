## apimcli add-env

Add Environment to Config file

### Synopsis


Add new environment and its related endpoints to the config file

```
apimcli add-env [flags]
```

### Examples

```
apimcli add-env -n production \
--registration https://localhost:9443/client-registration/v0.14/register \
--apim  https://localhost:9443 \
--token https://localhost:8243/token

apimcli add-env -n test \
--registration https://localhost:9443/client-registration/v0.14/register \
--import-export https://localhost:9443/api-import-export-2.6.0-v0 \
--list https://localhsot:9443/api/am/publisher/v0.14/apis \
--apim  https://localhost:9443 \
--token https://localhost:8243/token

apimcli add-env -n dev --apim https://localhost:9443 \
--token	https://localhost:8243/token \
--registration http://localhost:9763/client-registration/v0.14/register
```

### Options

```
      --admin string           Admin endpoint for the environment
      --api_list string        API List endpoint for the environment
      --apim string            API Manager endpoint for the environment
      --app_list string        Application List endpoint for the environment
  -h, --help                   help for add-env
      --import-export string   API Import Export endpoint for the environment
  -n, --name string            Name of the environment to be added
      --registration string    Registration endpoint for the environment
      --token string           Token endpoint for the environment
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO
* [apimcli](apimcli.md)	 - CLI for Importing and Exporting APIs and Applications

