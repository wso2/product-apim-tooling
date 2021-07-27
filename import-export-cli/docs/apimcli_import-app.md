## apimcli import-app

Import App

### Synopsis


Import an Application to an environment
Examples:
apimcli import-app -f qa/apps/sampleApp.zip -e dev
apimcli Import App -f staging/apps/sampleApp.zip -e prod -o testUser -u admin -p admin
apimcli import-app -f qa/apps/sampleApp.zip --preserveOwner --skipSubscriptions --skipKeys -e prod


```
apimcli import-app (--file <app-zip-file> --environment <environment-to-which-the-app-should-be-imported>) [flags]
```

### Options

```
  -e, --environment string   Environment from the which the Application should be imported (default "default")
  -f, --file string          Name of the Application to be imported
  -h, --help                 help for import-app
  -o, --owner string         Name of the target owner of the Application as desired by the Importer
  -p, --password string      Password
  -r, --preserveOwner        Preserves app owner
      --skipKeys             Skip importing keys of application
  -s, --skipSubscriptions    Skip subscriptions of the Application
      --update               Update application or create new
  -u, --username string      Username
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO
* [apimcli](apimcli.md)	 - CLI for Importing and Exporting APIs and Applications

