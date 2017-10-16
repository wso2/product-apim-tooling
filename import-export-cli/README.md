# CLI for Importing and Exporting APIs
## WSO2 API Manager 3.0.0

#### Usage 
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
        Example:
            wso2paim reset-user -e dev
```
* #### version
```bash
        wso2apim version 
``` 
        
#### Global Flags
```bash
    --verbose
    --help, -h
```
