## apictl get-keys

Generate access token to invoke the API

### Synopsis

Generate JWT token to invoke the API by subscribing to a default application for testing purposes

```
apictl get-keys [flags]
```

### Examples

```
apictl get-keys -n TwitterAPI -v 1.0.0 -e dev --provider admin
NOTE: Both the flags (--name (-n) and --environment (-e)) are mandatory
```

### Options

```
  -e, --environment string   Key generation environment
  -h, --help                 help for get-keys
  -n, --name string          API to generate keys
  -r, --provider string      Provider of the API
  -v, --version string       Version of the API
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl](apictl.md)	 - CLI for Importing and Exporting APIs and Applications

