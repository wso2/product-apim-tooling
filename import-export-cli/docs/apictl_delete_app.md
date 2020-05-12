## apictl delete app

Delete App

### Synopsis

Delete an Application from an environment

```
apictl delete app (--name <name-of-the-application> --owner <owner-of-the-application> --environment <environment-from-which-the-application-should-be-deleted>) [flags]
```

### Examples

```
apictl delete app -n TestApplication -o admin -e dev
apictl delete app -n SampleApplication -e production
NOTE: Both the flags (--name (-n), and --environment (-e)) are mandatory and the flag --owner (-o) is optional.
```

### Options

```
  -e, --environment string   Environment from which the Application should be deleted
  -h, --help                 help for app
  -n, --name string          Name of the Application to be deleted
  -o, --owner string         Owner of the Application to be deleted
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl delete](apictl_delete.md)	 - Delete an API/Application in an environment

