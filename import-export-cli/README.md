# CLI for Importing and Exporting APIs
## For WSO2 API Manager 3.0.0

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
    Extract the compressed archive generated to a desired location.
    
    Then execute `import-export-cli/apimcli` to start the application.
    
    Execute `import-export-cli/apimcli --help` for further instructions.

- ### Adding Environments
    Add environments by either manually editing `import-export-cli/bin/main_config.yaml` or using the command
    `apimcli set` command.
    
    Type `apimcli set --help` for detailed instructions
    
  ### Command Autocompletion (For Bash Only)
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
            apimcli import-api -f qa/TestAPI_1.2.1.zip -e dev -u admin -p 123456
            apimcli import-api -f staging/TestAPI_2.1.3.zip -e dev -u admin
            apimcli import-api -f production/TestAPI_3.1.0.zip -e dev -p 123456 
            apimcli import-api -f TestAPI_1.2.1.zip -e dev
```
* #### list apis
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
