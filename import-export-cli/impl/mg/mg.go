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

package mg

import (
	"errors"
	"strings"

	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const DefaultMgwAdapterEndpointSuffix = "/api/mgw/adapter/0.1"

const defaultTokenEndpointPath = "oauth2/token"
const apisResourcePath = "/apis"

func deriveTokenEndpointForMGAdapter(mgAdapterEndpoint string) string {
	if strings.HasSuffix(mgAdapterEndpoint, "/") {
		return mgAdapterEndpoint + defaultTokenEndpointPath
	} else {
		return mgAdapterEndpoint + "/" + defaultTokenEndpointPath
	}
}

func GetMgwAdapterInfo(env string) (mgwAdapterInfo MgwAdapterInfo, err error) {
	store, err := credentials.GetDefaultCredentialStore()
	if err != nil {
		return mgwAdapterInfo, err
	}
	mgToken, err := store.GetMGToken(env)
	if err != nil || mgToken.AccessToken == "" {
		err = errors.New("Error loading access token. " + err.Error())
		return mgwAdapterInfo, err
	}
	mgwAdapterEndpoints, err := utils.GetEndpointsOfMgwAdapterEnv(env, utils.MainConfigFilePath)
	if err != nil || mgwAdapterEndpoints.AdapterEndpoint == "" {
		err = errors.New("Error loading Adapter endpoint. " + err.Error())
		return mgwAdapterInfo, err
	}

	mgwAdapterInfo.Endpoint = mgwAdapterEndpoints.AdapterEndpoint
	mgwAdapterInfo.AccessToken = mgToken.AccessToken
	return mgwAdapterInfo, nil
}
