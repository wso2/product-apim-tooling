# CLI for Importing and Exporting APIs and Applications
## For WSO2 API Manager 2.2.0

Command Line tool for importing and exporting APIs/Applications between different API Environemnts

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
    Select a generated archive suitable for your platform (Mac, Windows, Linux) and extract it to a desired location and`cd` into it.
    Then execute `apimcli` to start the application.
    > NOTE: Execute `./apimcli` if the working directory is the same where the executable resides
    
    > Add the location of extracted folder to your system's $PATH variable to access the executable from anywhere 
    
    Execute `apimcli --help` for further instructions.

- ### Adding Environments
    Add environments by either manually editing `$HOME/.wso2apimcli/main_config.yaml` or using the command
    `apimcli add-env`.
    > NOTE: Directory structure for configuration files (`$HOME/.wso2apimcli`) will be created upon execution of `apimcli`
    
    Execute `apimcli add-env --help` for detailed instructions
    > Under environment, `api_manager_endpoint`, `registration_endpoint`, `token_endpoint` fields are mandatory. Others are optional
    
- ### Command Autocompletion (For Bash Only)
    Copy the file `apimcli_bash_completion.sh` to `/etc/bash_completion.d/` and source it with
    `source /etc/bash_completion.d/apimcli_bash_completion.sh` to enable bash auto-completion.

***

## Usage 
```bash
     apimcli [command]
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

### Commands
   * #### export-api
```bash
        Flags:
            Required:
                --name, -n
                --version, -v
                --provider, -r
                --environment, -e
            Optional:
                --username, -u
                --password, -p
                NOTE: user will be prompted to enter credentials if they are not provided with these flags
        Examples:
            apimcli export-api -n TestAPI -v 1.0.1 -r admin -e staging
            apimcli export-api -n TestAPI -v 1.0.1 -r admin -e staging -u admin -p 123456
            apimcli export-api -n TestAPI -v 1.0.1 -r admin -e staging -u admin
            apimcli export-api -n TestAPI -v 1.0.1 -r admin -e staging -p 123456
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
		--preserve-provider, -r 
                NOTE: user will be prompted to enter credentials if they are not provided with these flags
        Examples:
            apimcli import-api -f dev/TestAPI_1.0.0.zip -e dev
            apimcli import-api -f qa/TestAPI_2.0.0.zip -e dev -u admin -p 123456
            apimcli import-api -f staging/TestAPI_1.1.zip -e dev -u admin
            apimcli import-api -f production/TestAPI_3.0.1.zip -e dev -p 123456 
            apimcli import-api -f TestAPI -e dev
            apimcli import-api -f TestAPI -e dev -r=false
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

* #### export-app
```bash
        Flags
            Required:
                 --name, -n          
                 --owner, -o         
                 --environment, -e
            Optional
                 --username, -u
                 --password, -p   
                 NOTE: user will be prompted to enter credentials if they are not provided with these flags
        Examples:        
            apimcli export-app -n SampleApp -o admin -e dev
            apimcli export-app -n SampleApp -o admin -e prod         
```
* #### import-app
```bash    
        Flags
            Required
                  --file, -f          
                  --environment, -e   
            Optional
                  --skipSubscriptions, -s    
                  --owner, -o        
                  --preserveOwner, -r        
                  --username, -u      
                  --password, -p     
        Examples:      
            apimcli import-app -f qa/apps/sampleApp.zip -e dev
            apimcli Import-app -f staging/apps/sampleApp.zip -e prod -o testUser -u admin -p admin
            apimcli import-app -f qa/apps/sampleApp.zip --preserveOwner --skipSubscriptions -e staging               
```
* #### list apps
```bash
        Flags
            Required
                  --environment, -e          
                  --owner, -o         
            Optional
                  --username, -u             
                  --password, -p             
        Examples:
            apimcli list apps -e dev -o admin 
            apimcli list apps -e staging -o sampleUser -u admin -p 123456                         

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
                --apim (API Manager endpoint)
                --registration (Registration Endpoint)
                --token (Token Endpoint)
            Optional:
                --import-export (API Import Export Endpoint for environment)
                --list (API List endpoint for environment)
            
        Examples:
           apimcli add-env -n dev \
                --apim https://localhost:9443 \ 
                --registration https://localhost:9443/identity/connect/register \
                --token https: https://localhost:9443/oauth2/token
                
           apimcli add-env -n prod \
                --apim https://localhost:9443 \ 
                --registration https://localhost:9443/identity/connect/register \
                --import-export https://localhost:9443/api-import-export-2.2.0-v2 \
                --list https://localhsot:9443/api/am/publisher/v0.12/apis \
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
            --http-request-timeout
            --export-directory
        Examples:
            apimcli set --http-request-timeout 10000
            apimcli set --export-directory /home/user/exported 
```
