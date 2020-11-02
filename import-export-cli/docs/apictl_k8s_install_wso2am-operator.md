## apictl k8s install wso2am-operator

Install WSO2AM Operator

### Synopsis

Install WSO2AM Operator in the configured K8s cluster

```
apictl k8s install wso2am-operator [flags]
```

### Examples

```
apictl k8s install wso2am-operator
apictl k8s install wso2am-operator -f path/to/operator/configs
apictl k8s install wso2am-operator -f path/to/operator/config/file.yaml
```

### Options

```
  -f, --from-file string   Path to wso2am-operator directory
  -h, --help               help for wso2am-operator
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl k8s install](apictl_k8s_install.md)	 - Install an operator in the configured K8s cluster

