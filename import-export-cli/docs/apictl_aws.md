## apictl aws

AWS Api-gateway related commands

### Synopsis

AWS Api-gateway related commands such as init.

```
apictl aws [flags]
```

### Examples

```
apictl aws init -n PetStore -s Demo

apictl aws init --name PetStore --stage Demo

apictl aws init --name Shopping --stage Live

NOTE: Both flags --name (-n) and --stage (-s) are mandatory as both values are needed to get the openAPI from AWS API Gateway.
Make sure the API name and Stage name are correct.
Also make sure you have AWS CLI installed and configured before executing the aws init command.
Vist https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-welcome.html for more info
```

### Options

```
  -h, --help   help for aws
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [apictl](apictl.md)	 - CLI for Importing and Exporting APIs and Applications and Managing WSO2 Micro Integrator
* [apictl aws init](apictl_aws_init.md)	 - Initialize an API project for an AWS API

