package utils

// Kubernetes Constants
const Kubectl = "kubectl"
const K8sCreate = "create"
const K8sApply = "apply"
const K8sDelete = "delete"
const K8sRollOut = "rollout"
const K8sGet = "get"
const K8sConfigMap = "configmap"
const K8sSecret = "secret"
const K8sSecretDockerRegType = "docker-registry"

// WSO2 API Operator constats
const DefaultKubernetesMode = false
const ApiOpControllerConfigMap = "controller-config"
const ApiOpWso2Namespace = "wso2-system"
const ApiOpCrdApi = "apis.wso2.com"
const ApiOpCrdRateLimiting = "ratelimitings.wso2.com"
const ApiOpCrdSecurity = "securities.wso2.com"
const ApiOpCrdTargetEndpoint = "targetendpoints.wso2.com"

// Operator Hub Constants
const OlmCrdUrlTemplate = "https://github.com/operator-framework/operator-lifecycle-manager/releases/download/%s/crds.yaml"
const OlmOlmUrlTemplate = "https://github.com/operator-framework/operator-lifecycle-manager/releases/download/%s/olm.yaml"
const OperatorYamlUrl = "https://operatorhub.io/install/api-operator.yaml"
const OperatorCsv = "csv"
const OlmVersion = "0.13.0" //TODO: renuka export default to env variable and read from it

// constants of K8s ConfigMap: controller-config
const CtrlConfigRegType = "registryType"
const CtrlConfigReg = "repositoryName"

// TODO: renuka replace this url (configuration?)
const OperatorConfigFileUrl = `https://gist.githubusercontent.com/renuka-fernando/6d6c64c786e6d13742e802534de3da4e/raw/72fe433e463c14952fd5685ac04bfc5ed8ae357c/controller_conf.yaml`

// Registry specific config maps and secrets names
// Common
const ConfigJsonVolume = "config-json"

// Amazon ECR
const AwsCredentialsVolume = "aws-cred"

// GCR
const GcrSvcAccKeyVolume = "gcr-key"
const GcrSvcAccKeyFile = "gcr_key.json"
