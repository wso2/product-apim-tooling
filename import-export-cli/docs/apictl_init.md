## apictl init

Initialize a new project in given path

### Synopsis

Initialize a new project in given path. If a OpenAPI specification provided API will be populated with details from it

```
apictl init [project path] [flags]
```

### Examples

```
apictl init myapi --oas petstore.yaml
apictl init Petstore --oas https://petstore.swagger.io/v2/swagger.json
apictl init Petstore --oas https://petstore.swagger.io/v2/swagger.json --initial-state=PUBLISHED
apictl init MyAwesomeAPI --oas ./swagger.yaml -d definition.yaml
```

### Options

```
  -d, --definition string      Provide a YAML definition of API
  -f, --force                  Force create project (default false)
  -h, --help                   help for init
      --initial-state string   Provide the initial state of the API; Valid states: [CREATED PUBLISHED]
      --oas string             Provide an OpenAPI specification file for the API
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl](apictl.md)	 - CLI for Importing and Exporting APIs and Applications

