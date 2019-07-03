## apimcli export-app

Export App

### Synopsis


Export an Application from a specified  environment

```
apimcli export-app (--name <name-of-the-application> --owner <owner-of-the-application> --environment <environment-from-which-the-app-should-be-exported>) [flags]
```

### Examples

```
apimcli export-app -n SampleApp -o admin -e dev
apimcli export-app -n SampleApp -o admin -e prod
NOTE: Flag --name (-n) and --owner (-o) are mandatory
```

### Options

```
  -e, --environment string   Environment to which the Application should be exported
  -h, --help                 help for export-app
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
* [apimcli](apimcli.md)	 - CLI for Importing and Exporting APIs and Applications

