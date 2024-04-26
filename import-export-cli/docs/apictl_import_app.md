## apictl import app

Import App

### Synopsis

Import an Application to an environment

```
apictl import app (--file <app-zip-file> --environment <environment-to-which-the-app-should-be-imported>) [flags]
```

### Examples

```
apictl import app -f qa/apps/sampleApp.zip -e dev
apictl import app -f staging/apps/sampleApp.zip -e prod -o testUser
apictl import app -f qa/apps/sampleApp.zip --preserve-owner --skip-subscriptions -e prod
NOTE: Both the flags (--file (-f) and --environment (-e)) are mandatory
```

### Options

```
  -e, --environment string   Environment from the which the Application should be imported
  -f, --file string          Name of the ZIP file of the Application to be imported
  -h, --help                 help for app
      --ignore-tier          Ignores validation of the tier
  -o, --owner string         Name of the target owner of the Application as desired by the Importer
      --preserve-owner       Preserves app owner
      --skip-cleanup         Leave all temporary files created during import process
      --skip-keys            Skip importing keys of the Application
  -s, --skip-subscriptions   Skip subscriptions of the Application
      --update               Update the Application if it is already imported
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl import](apictl_import.md)	 - Import an API/API Product/Application to an environment

