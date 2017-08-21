## CLI for Importing and Exporting APIs
### WSO2 API Manager 3.0.0

Complete the config file `env_endpoints_all.yaml` before starting the tool.

#### Usage 
     wso2apim [command]
#### Commands
    export-api
        Flags
            --name, -n
            --version, -v
            --environment, -e
        Example:
            wso2apim export-api -n TestAPI -v 1.0.1 -e staging
    import-api
        Flags
            --name, -n (filename)
            --environment, -e
        Example:
            wso2apim import-api -n TestAPI.zip -e dev
    list
        Flags
            None
        Example:
            wso2apim list -e dev
            wso2ppim list -e staging 
    version
        wso2apim version 
#### Global Flags
    --verbose
    --help, -h
