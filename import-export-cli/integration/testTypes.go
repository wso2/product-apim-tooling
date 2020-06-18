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

package integration

import "github.com/wso2/product-apim-tooling/import-export-cli/integration/apim"

type credentials struct {
	username string
	password string
}

type apiImportExportTestArgs struct {
	apiProvider credentials
	ctlUser     credentials
	api         *apim.API
	srcAPIM     *apim.Client
	destAPIM    *apim.Client
}

type apiProductImportExportTestArgs struct {
	apiProductProvider   credentials
	ctlUser              credentials
	apiProduct           *apim.APIProduct
	srcAPIM              *apim.Client
	destAPIM             *apim.Client
	importApisFlag       bool
	updateApisFlag       bool
	updateApiProductFlag bool
}

type appImportExportTestArgs struct {
	appOwner    credentials
	ctlUser     credentials
	application *apim.Application
	srcAPIM     *apim.Client
	destAPIM    *apim.Client
}

type apiGetKeyTestArgs struct {
	ctlUser    credentials
	api        *apim.API
	apiProduct *apim.APIProduct
	apim       *apim.Client
}
