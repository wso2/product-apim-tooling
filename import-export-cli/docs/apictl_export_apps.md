## apictl export apps

Export Applications

### Synopsis

Export Applications of a given tenant from a specified environment

```
apictl export apps (--environment <environment-from-which-the-app-should-be-exported> --format <export-format> --force) [flags]
```

### Examples

```
apictl export apps -e dev --force
apictl export apps -e prod
NOTE: The flag (--environment (-e)) is mandatory
```

### Options

```
  -e, --environment string   Environment from which the Applications should be exported
      --force                Clean all the previously exported Apps of the given target tenant, in the given environment if any, and to export Apps from beginning
      --format string        File format of exported archive (json or yaml) (default "YAML")
  -h, --help                 help for apps
      --with-keys            Export keys for the applications
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl export](apictl_export.md)	 - Export an API/API Product/Application/Policy in an environment

