## CLI for Importing and Exporting APIs
### WSO2 API Manager 3.0.0

Complete the config file `env_endpoints_all.yaml` before starting the tool.

#### Usage 
     wso2apim [command]
#### Commands
    exportAPI
        Flags
            --name, -n
            --version, -v
            --environment, -e
        Example:
            wso2apim exportAPI -n TestAPI -v 1.0.1 -e staging
    importAPI
        Similar to exportAPI
        Example:
            wso2apim importAPI -n TestAPI -v 1.0.0 -e dev
    version
        wso2apim version 
#### Global Flags
    --verbose
    --help, -h
