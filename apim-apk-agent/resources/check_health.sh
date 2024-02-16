#!/bin/bash
# --------------------------------------------------------------------
# Copyright (c) 2023, WSO2 LLC. (http://wso2.com) All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
# -----------------------------------------------------------------------

# Debugging: Print out the values of variables
echo "APIM_APK_AGENT_GRPC_PORT: ${APIM_APK_AGENT_GRPC_PORT}"
echo "APIM_APK_AGENT_SERVER_NAME: ${APIM_APK_AGENT_SERVER_NAME}"
echo "APIM_APK_AGENT_PUBLIC_CERT_PATH: ${APIM_APK_AGENT_PUBLIC_CERT_PATH}"
echo "APIM_APK_AGENT_PRIVATE_KEY_PATH: ${APIM_APK_AGENT_PRIVATE_KEY_PATH}"

# Run grpc_health_probe with debugging information
grpc_health_probe -addr "127.0.0.1:${APIM_APK_AGENT_GRPC_PORT}" \
    -tls \
    -tls-ca-cert "${APIM_APK_AGENT_PUBLIC_CERT_PATH}" \
    -tls-client-cert "${APIM_APK_AGENT_PUBLIC_CERT_PATH}" \
    -tls-client-key "${APIM_APK_AGENT_PRIVATE_KEY_PATH}" \
    -tls-server-name ${APIM_APK_AGENT_SERVER_NAME} \
    -connect-timeout=3s