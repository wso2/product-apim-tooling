package utils

// constants of K8s ConfigMap: controller-config
const CtrlConfigRegType = "registryType"
const CtrlConfigReg = "dockerRegistry"

// TODO: renuka replace this url (configuration?)
const OperatorConfigFileUrl = `https://gist.githubusercontent.com/renuka-fernando/6d6c64c786e6d13742e802534de3da4e/raw/2ecc5a6f4ce8f1c463605783db31e0bd9489c602/controller_conf.yaml`

// registry specific config maps and secrets names
// Amazon ECR
const AmazonEcrDockerConfig = "docker-config"
