# CLI for Importing and Exporting APIs
## For WSO2 API Manager 2.1.x

Command Line tool for importing and exporting APIs between different API Environemnts

## Usage 
```bash
     apimcli [command]
```

### Commands
   * #### export-api
```bash
        Flags:
            Required:
                --name, -n
                --version, -v
                --environment, -e
            Optional:
                --username, -u
                --password, -p
        Examples:
            apimcli export-api -n TestAPI -v 1.0.1 -e staging
            apimcli export-api -n TestAPI -v 1.0.1 -e staging -u admin -p 123456
            apimcli export-api -n TestAPI -v 1.0.1 -e staging -u admin
            apimcli export-api -n TestAPI -v 1.0.1 -e staging -p 123456
```

* #### import-api
    
```bash
        Flags:
            Required:
                --file, -f
                --environment, -e
            Optional:
                --username, -u 
                --password, -p 
        Examples:
            apimcli import-api -f dev/TestAPI_1.0.0.zip -e dev
            apimcli import-api -f qa/TestAPI_2.0.0.zip -e dev -u admin -p 123456
            apimcli import-api -f staging/TestAPI_1.1.zip -e dev -u admin
            apimcli import-api -f production/TestAPI_3.0.1.zip -e dev -p 123456 
            apimcli import-api -f TestAPI -e dev
```
* #### list apis
```bash
        Flags:
            Required:
                --environment, -e
            Optional:
                --username, -u 
                --password, -p 
                --query, -q
        Examples:
            apimcli list apis -e dev
            apimcli list apis -e prod -q version:1.0.0 
            apimcli list apis -e prod -q provider:admin 
            apimcli list apis -e staging 
            apimcli list apis -e staging -u admin -p 123456
            apimcli list apis -e staging -p 123456
```

*  #### list envs
```bash
        Flags:
            None
        Example:
            apimcli list envs
```

* #### add-env
```bash
        Flags:
            Required:
                --name, -n (Name of the environment)
                --publisher, -p (Publisher endpoint)
                --registration, -r (Registration Endpoint)
                --token, -t (Token Endpoint)
            Examples:
                apimcli add-env -n dev \
                --apim https://localhost:9292/api/am/publisher/v1.0 \
                --registration https://localhost:9443/identity/connect/register \
                --token https: https://localhost:9443/oauth2/token
```
* #### remove-env
```bash
        Flags:
            Required:
                --name, -n (Name of the environment)
            Examples:
                apimcli remove-env -n dev
                
```

* #### reset-user
```bash
        Flags
            --environment, -e
        Examples:
            apimcli reset-user -e dev
```
* #### version
```bash
        apimcli version 
``` 

* #### set
```bash
        Flags
            --httpRequestTimeout
            --exportDirectory
        Examples:
            apimcli set --httpRequestTimeout 10000
            apimcli set --exportDirectory /home/user/exported 
```
        
#### Global Flags
```bash
    --verbose
    --insecure, -k
    --help, -h
```
