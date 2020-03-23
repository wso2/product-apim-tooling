## apictl remove-env

Remove Environment from Config file

### Synopsis

Remove Environment and its related endpoints from the config file

```
apictl remove-env [flags]
```

### Examples

```
apictl remove-env -e production
NOTE: The flag (--environment (-e)) is mandatory
```

### Options

```
  -e, --environment string   Name of the environment to be removed
  -h, --help                 help for remove-env
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl](apictl.md)	 - CLI for Importing and Exporting APIs and Applications

