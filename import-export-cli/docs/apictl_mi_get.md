## apictl mi get

Get information about artifacts deployed in a Micro Integrator instance

### Synopsis

Get information about artifacts deployed in a Micro Integrator instance in the environment specified by the flag (--environment, -e)

```
apictl mi get [flags]
```

### Examples

```
apictl mi get apis -e dev
apictl mi get endpoints -e dev
```

### Options

```
  -h, --help   help for get
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl mi](apictl_mi.md)	 - Micro Integrator related commands
* [apictl mi get apis](apictl_mi_get_apis.md)	 - Get information about apis deployed in a Micro Integrator
* [apictl mi get composite-apps](apictl_mi_get_composite-apps.md)	 - Get information about composite apps deployed in a Micro Integrator
* [apictl mi get connectors](apictl_mi_get_connectors.md)	 - Get information about connectors deployed in a Micro Integrator
* [apictl mi get data-services](apictl_mi_get_data-services.md)	 - Get information about data services deployed in a Micro Integrator
* [apictl mi get endpoints](apictl_mi_get_endpoints.md)	 - Get information about endpoints deployed in a Micro Integrator
* [apictl mi get inbound-endpoints](apictl_mi_get_inbound-endpoints.md)	 - Get information about inbound endpoints deployed in a Micro Integrator
* [apictl mi get local-entries](apictl_mi_get_local-entries.md)	 - Get information about local entries deployed in a Micro Integrator
* [apictl mi get log-levels](apictl_mi_get_log-levels.md)	 - Get information about a Logger configured in a Micro Integrator
* [apictl mi get logs](apictl_mi_get_logs.md)	 - List all the available log files
* [apictl mi get message-processors](apictl_mi_get_message-processors.md)	 - Get information about message processors deployed in a Micro Integrator
* [apictl mi get message-stores](apictl_mi_get_message-stores.md)	 - Get information about message stores deployed in a Micro Integrator
* [apictl mi get proxy-services](apictl_mi_get_proxy-services.md)	 - Get information about proxy services deployed in a Micro Integrator
* [apictl mi get sequences](apictl_mi_get_sequences.md)	 - Get information about sequences deployed in a Micro Integrator
* [apictl mi get tasks](apictl_mi_get_tasks.md)	 - Get information about tasks deployed in a Micro Integrator
* [apictl mi get templates](apictl_mi_get_templates.md)	 - Get information about templates deployed in a Micro Integrator
* [apictl mi get transaction-counts](apictl_mi_get_transaction-counts.md)	 - Retrieve transaction count
* [apictl mi get transaction-reports](apictl_mi_get_transaction-reports.md)	 - Generate transaction count summary report
* [apictl mi get users](apictl_mi_get_users.md)	 - Get information about users

