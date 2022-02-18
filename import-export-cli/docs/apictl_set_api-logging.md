## apictl set api-logging

Set the log level for an API in an environment

### Synopsis

Set the log level for an API in the environment specified

```
apictl set api-logging [flags]
```

### Examples

```
apictl set api-logging --api-id bf36ca3a-0332-49ba-abce-e9992228ae06 --log-level full -e dev --tenant-domain carbon.super
apictl set api-logging --api-id bf36ca3a-0332-49ba-abce-e9992228ae06 --log-level off -e dev --tenant-domain carbon.super
```

### Options

```
  -i, --api-id string          API ID
  -e, --environment string     Environment of the API which the log level should be set
  -h, --help                   help for api-logging
      --log-level string       Log Level
      --tenant-domain string   Tenant Domain
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl set](apictl_set.md)	 - Set configuration parameters or per API log levels

