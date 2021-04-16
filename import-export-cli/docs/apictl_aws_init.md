## apictl aws init

Initialize an API project for an AWS API

### Synopsis

Downloading the OpenAPI specification of an API from the AWS API Gateway to initialize a WSO2 API project

```
apictl aws init [flags]
```

### Examples

```
apictl aws init -n Petstore -s Demo
apictl aws init --name Petstore --stage Demo
apictl aws init --name Shopping --stage Live

NOTE: Both the flags --name (-n) and --stage (-s) are mandatory as both values are needed to get the openAPI from AWS API Gateway.
Make sure the API name and the Stage name are correct.
Also make sure you have AWS CLI installed and configured before executing the aws init command.
(Vist https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-welcome.html for more information)
```

### Options

```
  -f, --force          Force create project
  -h, --help           help for init
  -n, --name string    Name of the API to get from AWS Api Gateway
  -s, --stage string   Stage name of the API to get from AWS Api Gateway
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl aws](apictl_aws.md)	 - AWS Api-gateway related commands

