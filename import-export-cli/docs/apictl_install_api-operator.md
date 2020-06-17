## apictl install api-operator

Install API Operator

### Synopsis

Install API Operator in the configured K8s cluster

```
apictl install api-operator [flags]
```

### Examples

```
apictl install api-operator
apictl install api-operator -f path/to/operator/configs
apictl install api-operator -f path/to/operator/config/file.yaml
```

### Options

```
  -f, --from-file string       Path to API Operator directory
  -h, --help                   help for api-operator
  -c, --key-file string        Credentials file
  -p, --password string        Password of the given user
      --password-stdin         Prompt for password of the given user in the stdin (default false)
  -R, --registry-type string   Registry type: DOCKER_HUB | AMAZON_ECR |GCR | HTTP
  -r, --repository string      Repository name or URI
  -u, --username string        Username of the repository
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl install](apictl_install.md)	 - Install an operator

