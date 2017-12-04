# Contributing

Follow the following steps to make a contribution to this tool.

 * [Install Go (1.8.x) locally](https://golang.org/doc/install)
 * [Setup `$GOPATH` and `$GOROOT`](https://www.goinggo.net/2016/05/installing-go-and-your-workspace.html)
 * [Install Glide - Go Dependency Manager](https://github.com/Masterminds/glide#install)
 * Fork the respository: [wso2/product-apim-tooling](https://github.com/wso2/product-apim-tooling)
 * Do `git clone https://github.com/<your-username>/product-apim-tooling`
 
 NOTE: Project should be cloned such that the location would be `$GOPATH/src/github.com/wso2/product-apim-tooling` (This is mandatory for Go's dependencies to identity each other)
 * `cd` into `product-apim-tooling/import-export-cli`
 * Execute `glide install` inside `product-apim-tooling/import-export-cli` to download all dependencies
 * Make changes to the code
 * Commit to your fork
 * Send a pull request to [wso2/product-apim-tooling](https://github.com/wso2/product-apim-tooling)
 

## Tips
 * This command line tool is created on top of [Cobra](https://github.com/spf13/cobra)
 * Cobra's [README](https://github.com/spf13/cobra/blob/master/README.md) provides detailed information on developing command line tools.
 * For example: To add a new command named `apps` (so that the command would be `apimcli apps`) the following cobra command needs to be executed
 
    `cobra add apps`