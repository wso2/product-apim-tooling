## apimctl init

Initialize a new project in given path

### Synopsis


Initialize a new project in given path. If a OpenAPI specification provided API will be populated with details from it

```
apimctl init [project path] [flags]
```

### Examples

```
apimctl init myapi --oas petstore.yaml
apimctl init Petstore --oas https://petstore.swagger.io/v2/swagger.json
apimctl init MyAwesomeAPI --oas ./swagger.yaml -d definition.yaml
```

### Options

```
  -d, --definition string   Provide a YAML definition of API
  -f, --force               Force create project
  -h, --help                help for init
      --oas string          Provide an OpenAPI specification file for the API
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO
* [apimctl](apimctl.md)	 - CLI for Importing and Exporting APIs and Applications

