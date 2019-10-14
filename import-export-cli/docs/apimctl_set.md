## apimctl set

Set configuration

### Synopsis


Set configuration parameters. Use at least one of the following flags
* --http-request-timeout <time-in-milli-seconds>
* --export-directory <path-to-directory-where-apis-should-be-saved>
* --mode <mode-of-apimctl>

```
apimctl set [flags]
```

### Examples

```
apimctl set --http-request-timeout 3600 --export-directory /home/user/exported-apis
apimctl set --http-request-timeout 5000 --export-directory C:\Documents\exported
apimctl set --http-request-timeout 5000
apimctl set --token-type JWT
apimctl set --token-type OAUTH
apimctl set --mode kubernetes
apimctl set --mode default
```

### Options

```
      --export-directory string    Path to directory where APIs should be saved (default "/home/dinusha/.wso2apimctl/exported")
  -h, --help                       help for set
      --http-request-timeout int   Timeout for HTTP Client (default 10000)
  -m, --mode string                mode of apimctl
  -t, --token-type string          Type of the token to be generated
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO
* [apimctl](apimctl.md)	 - CLI for Importing and Exporting APIs and Applications

