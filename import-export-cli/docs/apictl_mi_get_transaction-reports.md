## apictl mi get transaction-reports

Generate transaction count summary report

### Synopsis

Generate the transaction count summary report at the given location for the given period of time.
If a location not provided, generate the report in current directory.
If an end date not provided, generate the report with values upto current date of the Micro Integrator in the environment specified by the flag --environment, -e

```
apictl mi get transaction-reports [start] [end] [flags]
```

### Examples

```
Example:
To generate transaction count report consisting data within a specified time period at a specified location
  apictl mi get transaction-reports 2020-05 2020-06 --path </dir_path> -e dev
To generate transaction count report with data from a given month upto the current month at a specified location
  apictl mi get transaction-reports 2020-01 -p </dir_path> -e dev
To generate transaction count report at the current location with data between 2020-01 and 2020-05
  apictl mi get transaction-reports 2020-01 2020-05 -e dev
NOTE: The [start] argument and the flag (--environment (-e)) is mandatory
```

### Options

```
  -e, --environment string   Environment to be searched
  -h, --help                 help for transaction-reports
  -p, --path string          destination file location
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl mi get](apictl_mi_get.md)	 - Get information about artifacts deployed in a Micro Integrator instance

