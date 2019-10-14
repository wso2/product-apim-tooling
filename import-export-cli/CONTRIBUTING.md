# Contributing

Follow the following steps to make a contribution to this tool.

 * [Install Go (1.12.x) locally](https://golang.org/doc/install)
 * Fork the respository: [wso2/product-apim-tooling](https://github.com/wso2/product-apim-tooling)
 * Clone the fork to any directory you want
 * Go to that directory and run `go mod vendor` or `go mod download` to download necessary dependenciese
 * Run ./build.sh to build the application
 * Make changes to the code
 * Commit to your fork
 * Send a pull request to [wso2/product-apim-tooling](https://github.com/wso2/product-apim-tooling)


## Tips
 * This command line tool is created on top of [Cobra](https://github.com/spf13/cobra)
 * Cobra's [README](https://github.com/spf13/cobra/blob/master/README.md) provides detailed information on developing command line tools.
 * For example: To add a new command named `apps` (so that the command would be `apictl apps`) the following cobra command needs to be executed

    `cobra add apps`