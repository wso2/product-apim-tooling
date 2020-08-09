## apictl set

Set configuration parameters

### Synopsis

Set configuration parameters. Use at least one of the following flags
* --http-request-timeout <time-in-milli-seconds>
* --export-directory <path-to-directory-where-apis-should-be-saved>
* --mode <mode-of-apictl>
* --vcs-deletion-enabled <enable-or-disable-project-deletion-via-vcs>
* --vcs-config-path <path-to-custom-vcs-config-file>

```
apictl set [flags]
```

### Examples

```
apictl set --http-request-timeout 3600 --export-directory /home/user/exported-apis
apictl set --http-request-timeout 5000 --export-directory C:\Documents\exported
apictl set --http-request-timeout 5000
apictl set --mode kubernetes
apictl set --mode default
apictl set --vcs-deletion-enabled=true
apictl set --vcs-config-path /home/user/custom/vcs-config.yaml
```

### Options

```
      --export-directory string    Path to directory where APIs should be saved (default "/home/wso2user/.wso2apictl/exported")
  -h, --help                       help for set
      --http-request-timeout int   Timeout for HTTP Client (default 10000)
  -m, --mode string                If mode is set to "k8s", apictl is capable of executing Kubectl commands. For example "apictl get pods" -> "kubectl get pods". To go back to the default mode, set the mode to "default" (default "default")
      --vcs-config-path string     Path to the VCS Configuration yaml file which keeps the VCS meta data
      --vcs-deletion-enabled       Specifies whether project deletion is allowed during deployment.
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl](apictl.md)	 - CLI for Importing and Exporting APIs and Applications

