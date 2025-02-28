## apictl set

Set configuration parameters, per API log levels or correlation component configurations

### Synopsis

Set configuration parameters. You can use one of the following flags
* --http-request-timeout <time-in-milli-seconds>
* --tls-renegotiation-mode <never|once|freely>
* --export-directory <path-to-directory-where-apis-should-be-saved>
* --vcs-deletion-enabled <enable-or-disable-project-deletion-via-vcs>
* --vcs-config-path <path-to-custom-vcs-config-file>
* --vcs-deployment-repo-path <path-to-deployment-repo-for-vcs>
* --vcs-source-repo-path <path-to-source-repo-for-vcs>
* --ai-thread-count <number-of-threads>
* --ai-token <on-prem-key-of-ai-features>

```
apictl set [flags]
```

### Examples

```
apictl set --http-request-timeout 3600 --export-directory /home/user/exported-apis
apictl set --http-request-timeout 5000 --export-directory C:\Documents\exported
apictl set --http-request-timeout 5000
apictl set --tls-renegotiation-mode freely
apictl set --vcs-deletion-enabled=true
apictl set --vcs-config-path /home/user/custom/vcs-config.yaml
apictl set --vcs-deployment-repo-path /home/user/custom/deployment
apictl set --vcs-source-repo-path /home/user/custom/source
apictl set api-logging --api-id bf36ca3a-0332-49ba-abce-e9992228ae06 --log-level full -e dev --tenant-domain carbon.super
apictl set correlation-logging --component-name http --enable true -e dev
apictl set --ai-thread-count 5
apictl set --ai-token ad232sda-asa2a-assdsd-sds43
```

### Options

```
      --ai-thread-count int               No of threads to be used by Marketplace Assistant for parallel processing
      --ai-token string                   Token (On prem key) of AI features
      --export-directory string           Path to directory where APIs should be saved (default "/Users/wso2user/.wso2apictl/exported")
  -h, --help                              help for set
      --http-request-timeout int          Timeout for HTTP Client (default 10000)
      --tls-renegotiation-mode string     Supported TLS renegotiation mode (default "never")
      --vcs-config-path string            Path to the VCS Configuration yaml file which keeps the VCS meta data
      --vcs-deletion-enabled              Specifies whether project deletion is allowed during deployment.
      --vcs-deployment-repo-path string   Path to the deoployment repository to be considered during VCS deploy
      --vcs-source-repo-path string       Path to the source repository to be considered during VCS deploy
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl](apictl.md)	 - CLI for Importing and Exporting APIs and Applications and Managing WSO2 Micro Integrator
* [apictl set api-logging](apictl_set_api-logging.md)	 - Set the log level for an API in an environment
* [apictl set correlation-logging](apictl_set_correlation-logging.md)	 - Set the correlation configs for a correlation logging component in an environment

