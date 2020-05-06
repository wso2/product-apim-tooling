## apictl set

Set configuration

### Synopsis

Set configuration parameters. Use at least one of the following flags
* --http-request-timeout <time-in-milli-seconds>
* --export-directory <path-to-directory-where-apis-should-be-saved>
* --mode <mode-of-apictl>

```
apictl set [flags]
```

### Examples

```
apictl set --http-request-timeout 3600 --export-directory /home/user/exported-apis
apictl set --http-request-timeout 5000 --export-directory C:\Documents\exported
apictl set --http-request-timeout 5000
apictl set --token-type JWT
apictl set --token-type OAUTH
apictl set --mode kubernetes
apictl set --mode default
```

### Options

```
      --export-directory string    Path to directory where APIs should be saved (default "/home/wasura/.wso2apictl/exported")
  -h, --help                       help for set
      --http-request-timeout int   Timeout for HTTP Client (default 10000)
  -m, --mode string                mode of apictl
  -t, --token-type string          Type of the token to be generated
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl](apictl.md)	 - CLI for Importing and Exporting APIs and Applications

