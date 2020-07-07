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

package testutils

import "github.com/wso2/product-apim-tooling/import-export-cli/integration/apim"

type Credentials struct {
	Username string
	Password string
}

type ApiImportExportTestArgs struct {
	ApiProvider Credentials
	CtlUser     Credentials
	Api         *apim.API
	SrcAPIM     *apim.Client
	DestAPIM    *apim.Client
}

type ApiProductImportExportTestArgs struct {
	ApiProductProvider   Credentials
	CtlUser              Credentials
	ApiProduct           *apim.APIProduct
	SrcAPIM              *apim.Client
	DestAPIM             *apim.Client
	ImportApisFlag       bool
	UpdateApisFlag       bool
	UpdateApiProductFlag bool
}

type appImportExportTestArgs struct {
	appOwner    Credentials
	ctlUser     Credentials
	application *apim.Application
	srcAPIM     *apim.Client
	destAPIM    *apim.Client
}

type apiGetKeyTestArgs struct {
	ctlUser    Credentials
	api        *apim.API
	apiProduct *apim.APIProduct
	apim       *apim.Client
}

type loginTestArgs struct {
	ctlUser Credentials
	srcAPIM *apim.Client
}

type setTestArgs struct {
	srcAPIM             *apim.Client
	exportDirectoryFlag string
	modeFlag            string
	tokenTypeFlag       string
	httpRequestTimeout  int
}

type initTestArgs struct {
	ctlUser        Credentials
	srcAPIM        *apim.Client
	initFlag       string
	definitionFlag string
	forceFlag      bool
	oasFlag        string
}
