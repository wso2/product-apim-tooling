# CLI for Importing and Exporting APIs and Applications and Managing WSO2 Micro Integrator
## For WSO2 API Manager 4.0.0

Command Line tool for importing and exporting APIs/Applications/API Products in between different API environments of WSO2 API Manager and managing WSO2 Micro Integrator

## Getting Started

- ### Setting up the development environment
    1. Install [Go 1.14](https://golang.org/dl)
    2. Fork the [repository](https://github.com/wso2/product-apim-tooling)
    3. Clone your fork into any directory
    5. `cd` into cloned directory and then cd into `product-apim-tooling/import-export-cli`
    6. Execute `go mod vendor` or `go mod download` to download all the dependencies
    7. Run `go generate ./...` to pack resources
    8. Run `go build apictl.go` to build and test program

- ### Building
    `cd` into `product-apim-tooling/import-export-cli`
    
    Execute `./build.sh -t apictl.go -v 4.0.0 -f` to build for all platforms.
    
    Created packages will be available at `build/target` directory

- ### Resources
    All the resources(Markdown files, templates etc) need to be stored inside box/resources directory
    You can create sub directories as you like.
    
    Run `go generate ./...` to bake your files into go files.
    Then refer to them using `box.Get(filename)`. Always use relative path to resources directory with a leading slash(/)
    
    For example: For a file in `box/resources/init/file` you can retrieve it using `box.Get("/init/file)`
    
    **DO NOT COMMIT blob.go to VCS**    

- ### Generating docs
    After changing commands run following to generate documents and shell completions
    `go run tools/gen.go`
    Commit changes to version control
      
- ### Running
    Select a generated archive suitable for your platform (Mac, Windows, Linux) and extract it to a desired location and`cd` into it.
    Then execute `apictl` to start the application.
    > NOTE: Execute `./apictl` if the working directory is the same where the executable resides
    
    > Add the location of extracted folder to your system's $PATH variable to access the executable from anywhere 
    
    Execute `apictl --help` for further instructions.

- ### Adding Environments
    Add environments by either manually editing `$HOME/.wso2apictl/main_config.yaml` or using the command
    `apictl add-env`.
    > NOTE: Directory structure for configuration files (`$HOME/.wso2apictl`) will be created upon execution of `apictl`
    
    Execute `apictl add-env --help` for detailed instructions
    > The flag `--environment` (-e) is mandatory
      You can either provide only the `--apim` flag, or all the other 5 flags (`--registration` `--publisher` `--devportal` `--admin` `--token`) without providing `--apim` flag.
      If you are omitting any of --registration --publisher --devportal --admin flags, you need to specify --apim flag with the API Manager endpoint.
      To add a micro integrator instance to an environment you can use the `--mi` flag.
    
- ### Command Autocomplete
    Copy the file `shell-completions/apictl_bash_completion.sh` to `/etc/bash_completion.d/` and source it with
    `source /etc/bash_completion.d/apictl_bash_completion.sh` to enable bash auto-completion.

***

## Command reference 

A reference for all commands can be found in [here](docs/apictl.md)
