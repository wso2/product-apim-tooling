## apictl mi secret create

Encrypt secrets

### Synopsis

Create secrets based on given arguments

```
apictl mi secret create [flags]
```

### Examples

```
To encrypt secret and get output on console
  apictl mi secret create
To encrypt secret and get output as a .properties file (stored in the security folder in apictl executable directory)
  apictl mi secret create -o file
To encrypt secret and get output as a .yaml file (stored in the security folder in apictl executable directory)
  apictl mi secret create -o k8
To bulk encrypt secrets defined in a properties file
  apictl mi secret create -f <file_path>
To bulk encrypt secrets defined in a properties file and get a .yaml file (stored in the security folder in apictl executable directory)
  apictl mi secret create -o k8 -f <file_path>
```

### Options

```
  -c, --cipher string      Encryption algorithm
  -f, --from-file string   Path to the properties file which contain secrets to be encrypted
  -h, --help               help for create
  -o, --output string      Get the output in yaml(k8) or properties(file) format. By default the output is printed to the console (default "console")
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl mi secret](apictl_mi_secret.md)	 - Manage sensitive information

