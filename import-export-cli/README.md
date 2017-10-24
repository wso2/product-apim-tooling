# CLI for Importing and Exporting APIs
## For WSO2 API Manager 3.0.0

Command Line tool for importing and exporting APIs between different API Environemnts

## Getting Started

- ### Building
    Execute ```build.sh``` to build for necessary platforms. Further instructions will be available upon executing 
    build.sh
      
- ### Running
    Extract the compressed archive generated to a desired location
    Then execute ```{extracted-location}/bin/wso2apim``` to start the application
    Execute ```{extracted-location}/bin/wso2apim --help``` for further instructions

- ### Adding Environments
    Add environments by either manually editing ```{extracted-location}/main_config.yaml``` or using the command
    ```wso2apim set``` command.
    Type ```wso2apim set --help``` for detailed instructions

<hr/>
<br/>

## Usage 
```bash
     wso2apim [command]
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
            wso2apim export-api -n TestAPI -v 1.0.1 -e staging
            wso2apim export-api -n TestAPI -v 1.0.1 -e staging -u admin -p 123456
            wso2apim export-api -n TestAPI -v 1.0.1 -e staging -u admin
            wso2apim export-api -n TestAPI -v 1.0.1 -e staging -p 123456
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
            wso2apim import-api -n TestAPI.zip -e dev
            wso2apim import-api -n TestAPI.zip -e dev -u admin -p 123456
            wso2apim import-api -n TestAPI.zip -e dev -u admin
            wso2apim import-api -n TestAPI.zip -e dev -p 123456 
            wso2apim import-api -n TestAPI -e dev
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
            wso2apim list -e dev
            wso2apim list -e staging 
            wso2apim list -e staging -u admin -p 123456
            wso2apim list -e staging -u admin
            wso2apim list -e staging -p 123456
```
* #### add-env
```bash
        Flags:
            Required:
                --name, -n (Name of the environment)
                --apim (API Manager endpoint)
                --registration (Registration Endpoint)
                --token (Token Endpoint)
            Examples:
                wso2apim add-env -n dev \
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
                wso2apim remove-env -n dev
                
```

* #### reset-user
```bash
        Flags
            --environment, -e
        Examples:
            wso2apim reset-user -e dev
```
* #### version
```bash
        wso2apim version 
``` 

* #### set
```bash
        Flags
            --httpRequestTimeout
            --exportDirectory
        Examples:
            wso2apim set --httpRequestTimeout 10000
            wso2apim set --exportDirectory /home/user/exported 
```
        
#### Global Flags
```bash
    --verbose
    --insecure, -k
    --help, -h
```
