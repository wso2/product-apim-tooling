## CLI for Importing and Exporting APIs
### WSO2 API Manager 3.0.0

Config file `env_endpoints_all.yaml` should be completed before starting the tool.

#### Usage 
     wso2apim [command]
#### Commands
    export-api
        Flags
            Required
                --name, -n
                --version, -v
                --environment, -e
            Optional
                --username, -u
                --password, -p
        Examples:
            wso2apim export-api -n TestAPI -v 1.0.1 -e staging
            wso2apim export-api -n TestAPI -v 1.0.1 -e staging -u admin -p 123456
            wso2apim export-api -n TestAPI -v 1.0.1 -e staging -u admin
            wso2apim export-api -n TestAPI -v 1.0.1 -e staging -p 123456 
    import-api
        Flags
            Required
                --name, -n
                --environment, -e
            Optional
                --username, -u 
                --password, -p 
        Examples:
            wso2apim import-api -n TestAPI.zip -e dev
            wso2apim import-api -n TestAPI.zip -e dev -u admin -p 123456
            wso2apim import-api -n TestAPI.zip -e dev -u admin
            wso2apim import-api -n TestAPI.zip -e dev -p 123456 
            wso2apim import-api -n TestAPI -e dev
    list
        Flags
            Required
                --environment, -e
            Optional
                --username, -u 
                --password, -p 
        Examples:
            wso2apim list -e dev
            wso2ppim list -e staging 
            wso2ppim list -e staging -u admin -p 123456
            wso2ppim list -e staging -u admin
            wso2ppim list -e staging -p 123456
    reset-user
        Flags
            --environment, -e
        Example:
            wso2paim reset-user -e dev
    version
        wso2apim version 
#### Global Flags
    --verbose
    --help, -h
