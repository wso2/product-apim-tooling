## apimcli list apis

Display a list of APIs in an environment

### Synopsis



Display a list of APIs in the environment specified by the flag --environment, -e

apimcli apis list -e dev
apimcli apis list -e dev -q version:1.0.0
apimcli apis list -e prod -q provider:admin
apimcli apis list -e staging -u admin -p admin


```
apimcli list apis [flags]
```

### Options

```
  -e, --environment string   Environment to be searched (default "default")
  -h, --help                 help for apis
  -p, --password string      Password
  -u, --username string      Username
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO
* [apimcli list](apimcli_list.md)	 - List APIs/Applications in an environment or List the environments

