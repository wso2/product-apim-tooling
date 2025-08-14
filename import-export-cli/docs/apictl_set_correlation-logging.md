## apictl set correlation-logging

Set the correlation configs for a correlation logging component in an environment

### Synopsis

Set the correlation configs for a correlation logging component in the environment specified
NOTE: The flags (--component-name (-i), --enable and --environment (-e)) are mandatory.

```
apictl set correlation-logging [flags]
```

### Examples

```
apictl set correlation-logging --component-name http --enable true -e dev
apictl set correlation-logging --component-name jdbc --enable true --denied-threads MessageDeliveryTaskThreadPool,HumanTaskServer,BPELServer -e dev
```

### Options

```
  -i, --component-name string   Component Name
      --denied-threads string   Denied Threads
      --enable string           Enable - true or false
  -e, --environment string      Environment where the correlation component configuration should be set
  -h, --help                    help for correlation-logging
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl set](apictl_set.md)	 - Set configuration parameters, per API log levels, MCP Server log levels or correlation component configurations

