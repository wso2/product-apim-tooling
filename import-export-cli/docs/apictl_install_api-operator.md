## apictl install api-operator

Install API Operator

### Synopsis

Install API Operator in the configured K8s cluster

```
apictl install api-operator [flags]
```

### Examples

```
apictl api-operator
apictl api-operator -f path/to/operator/configs
apictl api-operator -f path/to/operator/config/file.yaml
```

### Options

```
  -f, --from-file string   Path to API Operator directory
  -h, --help               help for api-operator
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl install](apictl_install.md)	 - Install an operator

