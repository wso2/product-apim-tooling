# Agent to integrate APK with APIM
## For WSO2 API Manager 4.3.0

## Getting Started

- ### Setting up the development environment
    1. Install [Go 1.19](https://golang.org/dl)
    2. Fork the [repository](https://github.com/wso2/product-apim-tooling)
    3. Clone your fork into any directory
    5. `cd` into cloned directory and then cd into `product-apim-tooling/apim-apk-agent`
    6. Execute `./gradlew build` to build the docker image
    7. `cd` into cloned directory and then cd into `product-apim-tooling/helm-charts`
    8. Run `helm install apim-apk-agent . -n apk` to deploy the agent in K8s
    9. Run `helm uninstall apim-apk-agent` to undeploy the agent in K8s