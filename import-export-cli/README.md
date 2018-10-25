# CLI for Importing and Exporting APIs and Applications
## For WSO2 API Manager 2.6.0

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
    > Under environment, `api_manager_endpoint`, `registration_endpoint`, `token_endpoint` fields are mandatory. 
    
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
		--insecure, -k
          	Allow connections to SSL sites without certs
            Optional:
                --username, -u
                --password, -p
                NOTE: user will be prompted to enter credentials if they are not provided with these flags
        Examples:
            apimcli export-api -n TestAPI -v 1.0.1 -r admin -e staging -k
            apimcli export-api -n TestAPI -v 1.0.1 -r admin -e staging -u admin -p 123456 -k
            apimcli export-api -n TestAPI -v 1.0.1 -r admin -e staging -u admin -k
            apimcli export-api -n TestAPI -v 1.0.1 -r admin -e staging -p 123456 -k
```

* #### import-api
    
```bash
        Flags:
            Required:
                --file, -f
                --environment, -e
		--insecure, -k
          	Allow connections to SSL sites without certs
            Optional:
                --username, -u 
                --password, -p
		--preserve-provider, -r 
                NOTE: user will be prompted to enter credentials if they are not provided with these flags
        Examples:
            apimcli import-api -f dev/TestAPI_1.0.0.zip -e dev -k
            apimcli import-api -f qa/TestAPI_2.0.0.zip -e dev -u admin -p 123456 -k
            apimcli import-api -f staging/TestAPI_1.1.zip -e dev -u admin -k
            apimcli import-api -f production/TestAPI_3.0.1.zip -e dev -p 123456 -k
            apimcli import-api -f TestAPI -e dev -k
            apimcli import-api -f TestAPI -e dev -r=false -k
```
* #### list apis
```bash
        Flags:
            Required:
                --environment, -e
		--insecure, -k
          	Allow connections to SSL sites without certs
            Optional:
                --username, -u 
                --password, -p 
                NOTE: user will be prompted to enter credentials if they are not provided with these flags
                --query, -q
        Examples:
            apimcli list apis -e dev -k
            apimcli list apis -e prod -q version:1.0.0 -k
            apimcli list apis -e prod -q provider:admin -k
            apimcli list apis -e staging -k
            apimcli list apis -e staging -u admin -p 123456 -k
            apimcli list apis -e staging -p 123456 -k
```

* #### export-app
```bash
        Flags
            Required:
                 --name, -n          
                 --owner, -o         
                 --environment, -e
		 --insecure, -k
          	Allow connections to SSL sites without certs
            Optional
                 --username, -u
                 --password, -p   
                 NOTE: user will be prompted to enter credentials if they are not provided with these flags
        Examples:        
            apimcli export-app -n SampleApp -o admin -e dev -k
            apimcli export-app -n SampleApp -o admin -e prod -k        
```
* #### import-app
```bash    
        Flags
            Required
                  --file, -f          
                  --environment, -e   
		  --insecure, -k
          	Allow connections to SSL sites without certs
            Optional
                  --skipSubscriptions, -s    
                  --owner, -o        
                  --preserveOwner, -r        
                  --username, -u      
                  --password, -p     
        Examples:      
            apimcli import-app -f qa/apps/sampleApp.zip -e dev -k
            apimcli import-app -f staging/apps/sampleApp.zip -e prod -o testUser -u admin -p admin -k
            apimcli import-app -f qa/apps/sampleApp.zip --preserveOwner --skipSubscriptions -e staging -k               
```
* #### list apps
```bash
        Flags
            Required
                  --environment, -e          
                  --owner, -o
		  --insecure, -k
          	  Allow connections to SSL sites without certs
            Optional
                  --username, -u             
                  --password, -p             
        Examples:
            apimcli list apps -e dev -o admin -k
            apimcli list apps -e staging -o sampleUser -u admin -p 123456 -k                       

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
                --registration (Registration endpoint)
                --apim (API Manager endpoint) 
                --token (Token endpoint) 
		--import-export (Endpoint for environment) 
                --admin (Admin REST API endpoint) 
                --api_list (API listing REST API endpoint)
                --app_list (Application listing REST API endpoint)
            
        Examples:  
           apimcli add-env -n dev \
                      --registration https://localhost:9443/client-registration/v0.14/register \
                      --apim https://localhost:9443 \
                      --token https://localhost:8243/token \
                      --import-export https://localhost:9443/api-import-export-2.6.0-v0 \
                      --admin https://localhost:9443/api/am/admin/v0.14 \
                      --api_list https://localhost:9443/api/am/publisher/v0.14/apis \
                      --app_list https://localhost:9443/api/am/store/v0.14/applications
		      
	   apimcli add-env -n prod \
                      --registration https://localhost:9444/client-registration/v0.14/register \
                      --apim https://localhost:9444 \
                      --token https://localhost:8244/token \
                      --import-export https://localhost:9444/api-import-export-2.6.0-v0 \
                      --admin https://localhost:9444/api/am/admin/v0.14 \
                      --api_list https://localhost:9444/api/am/publisher/v0.14/apis \
                      --app_list https://localhost:9444/api/am/store/v0.14/applications
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
