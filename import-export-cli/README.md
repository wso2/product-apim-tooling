# CLI for Importing and Exporting APIs
## For WSO2 API Manager 3.0.0

Command Line tool for importing and exporting APIs between different API Environemnts

## Getting Started

- ### Building
    Execute ```build.sh``` to build for necessary platforms. Further instructions will be available upon executing 
    build.sh
      
- ### Running
    Extract the compressed archive generated to a desired location
    Then execute ```import-export-cli/bin/apimcli``` to start the application
    Execute ```import-export-cli/bin/apimcli --help``` for further instructions

- ### Adding Environments
    Add environments by either manually editing ```import-export-cli/bin/main_config.yaml``` or using the command
    ```apimcli set``` command.
    Type ```apimcli set --help``` for detailed instructions

<hr/>
<br/>

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
                --name, -n
                --environment, -e
            Optional:
                --username, -u 
                --password, -p 
        Examples:
            apimcli import-api -n TestAPI.zip -e dev
            apimcli import-api -n TestAPI.zip -e dev -u admin -p 123456
            apimcli import-api -n TestAPI.zip -e dev -u admin
            apimcli import-api -n TestAPI.zip -e dev -p 123456 
            apimcli import-api -n TestAPI -e dev
```
* #### list
```bash
        Flags:
            Required:
                --environment, -e
            Optional:
                --username, -u 
                --password, -p 
        Examples:
            apimcli list -e dev
            apimcli list -e staging 
            apimcli list -e staging -u admin -p 123456
            apimcli list -e staging -u admin
            apimcli list -e staging -p 123456
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
