## apictl export app

Export App

### Synopsis

Export an Application from a specified  environment

```
apictl export app (--name <name-of-the-application> --owner <owner-of-the-application> --environment <environment-from-which-the-app-should-be-exported>) [flags]
```

### Examples

```
apictl export app -n SampleApp -o admin -e dev
apictl export app -n SampleApp -o admin -e prod
NOTE: All the 3 flags (--name (-n), --owner (-o) and --environment (-e)) are mandatory
```

### Options

```
  -e, --environment string   Environment to which the Application should be exported
      --format string        File format of exported archive (json or yaml) (default "YAML")
  -h, --help                 help for app
  -n, --name string          Name of the Application to be exported
  -o, --owner string         Owner of the Application to be exported
      --with-keys            Export keys for the application 
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl export](apictl_export.md)	 - Export an API/MCPServer/API Product/Application/Policy in an environment

