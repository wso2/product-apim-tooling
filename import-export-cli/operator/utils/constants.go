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
const ApiOpCrdApi = "apis.wso2.com"
const ApiOpCrdRateLimiting = "ratelimitings.wso2.com"
const ApiOpCrdSecurity = "securities.wso2.com"
const ApiOpCrdTargetEndpoint = "targetendpoints.wso2.com"

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
