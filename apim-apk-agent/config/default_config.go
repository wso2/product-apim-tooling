/*
 *  Copyright (c) 2024, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

package config

// Configuration object which is populated with default values.
var defaultConfig = &Config{
	ControlPlane: controlPlane{
		Enabled:              false,
		ServiceURL:           "https://apim:9443/",
		ServiceURLDeprecated: UnassignedAsDeprecated,
		Username:             "admin",
		Password:             "$env{cp_admin_pwd}",
		EnvironmentLabels:    []string{"Default"},
		RetryInterval:        5,
		SkipSSLVerification:  false,
		BrokerConnectionParameters: brokerConnectionParameters{
			EventListeningEndpoints: []string{"amqp://admin:$env{cp_admin_pwd}@apim:5672?retries='10'&connectdelay='30'"},
			ReconnectInterval:       5000, //in milli seconds
			ReconnectRetryCount:     60,
		},
		SendRevisionUpdate: false,
		HTTPClient: httpClient{
			RequestTimeOut: 30,
		},
		RequestWorkerPool: requestWorkerPool{
			PoolSize:              4,
			QueueSizePerPool:      1000,
			PauseTimeAfterFailure: 5,
		},
		InternalKeyIssuer: "http://am.wso2.com:443/token",
	},
	Agent: agent{
		Enabled: true,
		Keystore: keystore{
			KeyPath:  "/home/wso2/security/keystore/apk-agent.key",
			CertPath: "/home/wso2/security/keystore/apk-agent.crt",
		},
		TrustStore: truststore{
			Location: "/home/wso2/security/truststore",
		},
	},
	Metrics: metrics{
		Enabled: false,
		Port:    18006,
		Type:    "prometheus",
	},
}
