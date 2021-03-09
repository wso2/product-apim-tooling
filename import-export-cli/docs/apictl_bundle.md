## apictl bundle

Archive any source project artifact to zip format

### Synopsis

Archive API, Application or API Product projects to a zip format. Bundle name will have project name and version

```
apictl bundle [flags]
```

### Examples

```
apictl bundle -s /home/prod/APIs/API1-1.0.0 -d /home/prod/Projects/
apictl bundle -s /home/prod/APIs/API1-1.0.0 
NOTE: The flag (--source (-s)) is mandatory.
```

### Options

```
  -d, --destination string   Path of the directory where the bundle should be generated
  -h, --help                 help for bundle
  -s, --source string        Path of the source directory to bundle
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl](apictl.md)	 - CLI for Importing and Exporting APIs and Applications and Managing WSO2 Micro Integrator

