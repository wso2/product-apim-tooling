## apictl mi get transaction-counts

Retrieve transaction count

### Synopsis

Retrieve transaction count based on the given year and month.
If year and month not provided, retrieve the count for the current year and month of Micro Integrator in the environment specified by the flag --environment, -e

```
apictl mi get transaction-counts [year] [month] [flags]
```

### Examples

```
To get the transaction count for the current month
  apictl mi get transaction-counts -e dev
To get the transaction count for a specific month
  apictl mi get transaction-counts 2020 06 -e dev
NOTE: The flag (--environment (-e)) is mandatory
```

### Options

```
  -e, --environment string   Environment to be searched
      --format string        Pretty-print using Go Templates. Use "{{ jsonPretty . }}" to list all fields
  -h, --help                 help for transaction-counts
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl mi get](apictl_mi_get.md)	 - Get information about artifacts deployed in a Micro Integrator instance

