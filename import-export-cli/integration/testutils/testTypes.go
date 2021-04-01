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
	ApiProvider      Credentials
	CtlUser          Credentials
	Api              *apim.API
	SrcAPIM          *apim.Client
	DestAPIM         *apim.Client
	OverrideProvider bool
	ParamsFile       string
}

type ApiProductImportExportTestArgs struct {
	ApiProviders         map[string]Credentials
	ApiProductProvider   Credentials
	CtlUser              Credentials
	ApiProduct           *apim.APIProduct
	SrcAPIM              *apim.Client
	DestAPIM             *apim.Client
	ImportApisFlag       bool
	UpdateApisFlag       bool
	UpdateApiProductFlag bool
	ParamsFile           string
}

type AppImportExportTestArgs struct {
	AppOwner    Credentials
	CtlUser     Credentials
	Application *apim.Application
	SrcAPIM     *apim.Client
	DestAPIM    *apim.Client
}

type ApiGetKeyTestArgs struct {
	CtlUser    Credentials
	Api        *apim.API
	ApiProduct *apim.APIProduct
	Apim       *apim.Client
}

type loginTestArgs struct {
	ctlUser Credentials
	srcAPIM *apim.Client
}

type SetTestArgs struct {
	SrcAPIM             *apim.Client
	ExportDirectoryFlag string
	modeFlag            string
	TokenTypeFlag       string
	httpRequestTimeout  int
}

type InitTestArgs struct {
	CtlUser        Credentials
	SrcAPIM        *apim.Client
	InitFlag       string
	definitionFlag string
	ForceFlag      bool
	OasFlag        string
	APIName        string
	srcAPIM        *apim.Client
}

type ApiChangeLifeCycleStatusTestArgs struct {
	ApiProvider   Credentials
	CtlUser       Credentials
	Api           *apim.API
	APIM          *apim.Client
	Action        string
	Provider      string
	ExpectedState string
}

type AWSInitTestArgs struct {
	CtlUser        Credentials
	SrcAPIM        *apim.Client
	ApiNameFlag		 string 
	ApiStageNameFlag string
	InitFlag       string
}

type GenDeploymentDirTestArgs struct {
	Source      string
	Destination string
}
