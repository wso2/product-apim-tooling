## apictl import api-product

Import API Product

### Synopsis

Import an API Product to an environment

```
apictl import api-product (--file <path-to-api-product> --environment <environment-to-which-the-api-product-should-be-imported>) [flags]
```

### Examples

```
apictl import api-product -f qa/LeasingAPIProduct.zip -e dev
apictl import api-product -f staging/CreditAPIProduct.zip -e production --update-api-product
apictl import api-product -f ~/myapiproduct -e production
apictl import api-product -f ~/myapiproduct -e production --update-api-product --update-apis
NOTE: Both the flags (--file (-f) and --environment (-e)) are mandatory
```

### Options

```
  -e, --environment string   Environment from the which the API Product should be imported
  -f, --file string          Name of the API Product to be imported
  -h, --help                 help for api-product
      --import-apis          Import dependent APIs associated with the API Product
      --preserve-provider    Preserve existing provider of API Product after importing (default true)
      --skipCleanup          Leave all temporary files created during import process
      --update-api-product   Update an existing API Product or create a new API Product
      --update-apis          Update existing dependent APIs associated with the API Product
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl import](apictl_import.md)	 - Import an API Product to an environment

