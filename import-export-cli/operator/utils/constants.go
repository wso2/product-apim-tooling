package utils

// constants of K8s ConfigMap: controller-config
const CtrlConfigRegType = "registryType"
const CtrlConfigReg = "repositoryName"

// TODO: renuka replace this url (configuration?)
const OperatorConfigFileUrl = `https://gist.githubusercontent.com/renuka-fernando/6d6c64c786e6d13742e802534de3da4e/raw/72fe433e463c14952fd5685ac04bfc5ed8ae357c/controller_conf.yaml`

// registry specific config maps and secrets names
// Amazon ECR
const AwsCredentialsVolume = "aws-cred"
const GcrSvcAccKeyVolume = "gcr-key"
const GcrSvcAccKeyFile = "gcr_key.json"
