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

type TestCaseUsers struct {
	Description   string
	ApiCreator    Credentials
	ApiPublisher  Credentials
	ApiSubscriber Credentials
	Admin         Credentials
	CtlUser       Credentials
}

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
	ImportFilePath   string
	Revision         string
	IsDeployed       bool
	IsLatest         bool
	Update           bool
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
	Revision             string
	IsLatest             bool
}

type ThrottlePolicyImportExportTestArgs struct {
	Admin          Credentials
	CtlUser        Credentials
	Policy         map[string]interface{}
	SrcAPIM        *apim.Client
	DestAPIM       *apim.Client
	ParamsFile     string
	ImportFilePath string
	Update         bool
}

type AppImportExportTestArgs struct {
	AppOwner          Credentials
	CtlUser           Credentials
	Application       *apim.Application
	SrcAPIM           *apim.Client
	DestAPIM          *apim.Client
	PreserveOwner     bool
	UpdateFlag        bool
	WithKeys          bool
	SkipKeys          bool
	SkipSubscriptions bool
	ImportFilePath    string
}

type ApiGetKeyTestArgs struct {
	CtlUser    Credentials
	Api        *apim.API
	ApiProduct *apim.APIProduct
	Apim       *apim.Client
}

type UndeployTestArgs struct {
	CtlUser     Credentials
	Api         *apim.API
	ApiProduct  *apim.APIProduct
	SrcAPIM     *apim.Client
	RevisionNo  string
	GatewayEnvs []string
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
	DefinitionFlag string
	ForceFlag      bool
	OasFlag        string
	APIName        string
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
	CtlUser          Credentials
	SrcAPIM          *apim.Client
	ApiNameFlag      string
	ApiStageNameFlag string
	InitFlag         string
}

type GenDeploymentDirTestArgs struct {
	Source      string
	Destination string
}

type ApiProductChangeLifeCycleStatusTestArgs struct {
	ApiProvider   Credentials
	CtlUser       Credentials
	ApiProduct    *apim.APIProduct
	APIM          *apim.Client
	Action        string
	Provider      string
	ExpectedState string
}

type ApiLoggingTestArgs struct {
	Apis         []*apim.API
	APIM         *apim.Client
	CtlUser      Credentials
	TenantDomain string
	ApiId        string
	LogLevel     string
}
