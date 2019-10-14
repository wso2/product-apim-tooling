## apimctl import-app

Import App

### Synopsis


Import an Application to an environment

```
apimctl import-app (--file <app-zip-file> --environment <environment-to-which-the-app-should-be-imported>) [flags]
```

### Examples

```
apimctl import-app -f qa/apps/sampleApp.zip -e dev
apimctl Import App -f staging/apps/sampleApp.zip -e prod -o testUser
apimctl import-app -f qa/apps/sampleApp.zip --preserveOwner --skipSubscriptions -e prod
```

### Options

```
  -e, --environment string   Environment from the which the Application should be imported
  -f, --file string          Name of the Application to be imported
  -h, --help                 help for import-app
  -o, --owner string         Name of the target owner of the Application as desired by the Importer
      --preserveOwner        Preserves app owner
      --skipKeys             Skip importing keys of application
  -s, --skipSubscriptions    Skip subscriptions of the Application
      --update               Update application or create new
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO
* [apimctl](apimctl.md)	 - CLI for Importing and Exporting APIs and Applications

