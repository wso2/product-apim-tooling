# CLI for Importing and Exporting APIs
## For WSO2 API Manager 2.1.x

Command Line tool for importing and exporting APIs between different API Environemnts

## Getting Started

- ### Setting up the development environment
    1. Install [Go 1.8.x](https://golang.org/dl)
    2. Setup `$GOROOT` and `$GOPATH` correctly : [Tutorial](https://www.goinggo.net/2016/05/installing-go-and-your-workspace.html) 
    3. [Install Glide](https://github.com/Masterminds/glide#install) - The Dependency manager for Go
    4. Clone the [repository](https://github.com/wso2/product-apim-tooling) to your `$GOPATH` such that the project location is `$GOPATH/src/github.com/wso2/product-apim-tooling` (This is mandatory for Go's dependencies to identify each other)
    5. `cd` into `product-apim-tooling/import-export-cli` 
    6. Execute `glide install` to download all the dependencies

- ### Building
    `cd` into `product-apim-tooling/import-export-cli`
    
    Execute `./build.sh -t apimcli.go -v 1.0.0 -f` to build for all platforms.
    
    Created packages will be available at `build/target` directory
      
- ### Running
    Select a generated archive suitable for your platform and extract it to a desired location.
    
    Then execute `apimcli` to start the application.
    
    Execute `apimcli --help` for further instructions.

- ### Adding Environments
    Add environments by either manually editing `main_config.yaml` or using the command
    `apimcli set`.
    
    Execute `set --help` for detailed instructions
    
- ### Command Autocompletion (For Bash Only)
    Copy the file `apimcli_bash_completion.sh` to `/etc/bash_completion.d/` and source it with
    `source /etc/bash_completion.d/apimcli_bash_completion.sh` to enable bash auto-completion.

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
                NOTE: user will be prompted to enter credentials if they are not provided with these flags
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
                NOTE: user will be prompted to enter credentials if they are not provided with these flags
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
                NOTE: user will be prompted to enter credentials if they are not provided with these flags
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
                --apim, -a (API Manager endpoint)
                --registration, -r (Registration Endpoint)
                --token, -t (Token Endpoint)
            Examples:
                apimcli add-env -n dev \
                --apim https://localhost:9443 \ 
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
           Enable verbose logs (Provides more information on execution)
      --insecure, -k
          Allow connections to SSL sites without certs
      --help, -h
          Display information and example usage of a command
```
