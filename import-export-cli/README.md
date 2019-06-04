# CLI for Importing and Exporting APIs and Applications
## For WSO2 API Manager 2.6.0

Command Line tool for importing and exporting APIs/Applications between different API Environemnts

## Getting Started

- ### Setting up the development environment
    1. Install [Go 1.12.x](https://golang.org/dl)
    2. Fork the [repository](https://github.com/wso2/product-apim-tooling)
    3. Clone your fork into any directory
    5. `cd` into cloned directory and then cd into `product-apim-tooling/import-export-cli`
    6. Execute `go mod vendor` or `go mod download` to download all the dependencies

- ### Building
    `cd` into `product-apim-tooling/import-export-cli`
    
    Execute `./build.sh -t apimcli.go -v 1.0.0 -f` to build for all platforms.
    
    Created packages will be available at `build/target` directory

- ### Generating docs
    After changing commands run following to generate documents and shell completions
    `go run tools/gen.go`
    Commit changes to version control
      
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

## Command reference 

A reference for all commands can be found in [here](docs/apimcli.md)
