/*
*  Copyright (c) WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
*
*  WSO2 Inc. licenses this file to you under the Apache License,
*  Version 2.0 (the "License"); you may not use this file except
*  in compliance with the License.
*  You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing,
* software distributed under the License is distributed on an
* "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
* KIND, either express or implied.  See the License for the
* specific language governing permissions and limitations
* under the License.
 */

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
const ApiOpCrdSecurity = "securities.wso2.com"

const ApiOperatorConfigsUrlTemplate = "https://github.com/wso2/K8s-api-operator/releases/download/%s/api-operator-configs.yaml"
const ApiOperatorVersionValidationUrlTemplate = "https://github.com/wso2/K8s-api-operator/tree/%s"
const ApiOperatorFindVersionUrl = "https://github.com/wso2/K8s-api-operator/releases"
const DefaultApiOperatorVersion = "v1.1.0"
const ApiOperatorVersionEnvVariable = "WSO2_API_OPERATOR_VERSION"

// constants of K8s ConfigMap: controller-config
const CtrlConfigRegType = "registryType"
const CtrlConfigReg = "repositoryName"

// Registry specific config maps and secrets names
const ConfigJsonVolume = "config-json"
const AwsCredentialsVolume = "aws-cred"
const GcrSvcAccKeyVolume = "gcr-key"
const GcrSvcAccKeyFile = "gcr_key.json"
