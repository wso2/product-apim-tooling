/*
*  Copyright (c) 2005-2017, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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

import "github.com/renstrom/dedent"

const DoneMsg string = "Done!\n"
const UnableToConnectMsg = "There was a problem connecting, please try again"
const InvalidCredentialsMsg string = "Invalid Credentials"
const UsernamePasswordEmptyMsg string = "Username and Password cannot be empty"
const InvalidExpiredRefreshTokenMsg string = "Your session has timed out"
const ErrorReadingResponseMsg string = "There was an error reading the response from the server"
const ErrorProcessingResponseMsg string = "There was an error processing the response from the server"

const DoYouWantToContinueMsg_DefaultYes string = "Do you want to continue? [Y/n] "
const DoYouWantToContinueMsg_DefaultNo string = "Do you want to continue? [y/N] "

const RootCmdShortDesc string = "CLI for Importing and Exporting APIs"

var RootCmdLongDesc string = dedent.Dedent(`
		` + ProjectName + ` is a CLI for Importing and Exporting APIs between different environments
		(Production, Staging, QA etc.)
		`)

// Init command related usage info
const InitCmdShortDesc string = "Initialize " + ProjectName + " with your WSO2 credentials"

var InitCmdLongDesc = dedent.Dedent(`
		Initialize ` + ProjectName + ` with your WSO2 credentials

		You need a WSO2 account to start using ` + ProjectName + `.
		Don't have one yet? Sign up at https://wso2.com/user/register
		`)

var InitCmdExamples = dedent.Dedent(`
		<InitCmdExamples>
		`)

// Version command related usage info
const VersionCmdShortDesc string = "Display Version on current `+ProjectName+`"

var VersionCmdLongDesc string = dedent.Dedent(`
		Display the current version of this command line tool
		`)

var VersionCmdExamples = dedent.Dedent(`
		` + ProjectName + ` version
		`)

// ImportAPI command related usage info
const ImportAPICmdShortDesc string = "Import API"

var ImportAPICmdLongDesc string = "Import an API to an environment"

var ImportAPICmdExamples = dedent.Dedent(`
		Examples:
		` + ProjectName + ` importAPI -n TwitterAPI -v 1.0.0 -e dev
		` + ProjectName + ` importAPI -n FacebookAPI -v 2.1.0 -e production
	`)

// ExportAPI command related usage info
const ExportAPICmdShortDesc string = "Export API"

var ExportAPICmdLongDesc string = "Export an API from an environment"

var ExportAPICmdExamples = dedent.Dedent(`
		Examples:
		` + ProjectName + ` exportAPI -n TwitterAPI -v 1.0.0 -e dev
		` + ProjectName + ` exportAPI -n FacebookAPI -v 2.1.0 -e production
	`)

// List command related usage Info

const ListCmdShortDesc string = "List APIs in an environment"

var ListCmdLongDesc string = dedent.Dedent(`
			Display a list containing all the APIs available in the environment specified by flag (--environment, -e)
	`)

var ListCmdExamples = dedent.Dedent(`
		Examples:
		` + ProjectName + ` list -e dev
		` + ProjectName + ` list -e staging
	`)

// ResetUser command related usage Info

const ResetUserCmdShortDesc string = "Reset user of an environment"

var ResetUserCmdLongDesc = dedent.Dedent(`
		Reset user data of a particular environment (Clear the entry in env_keys_all.yaml file)
	`)

var ResetUserCmdExamples = dedent.Dedent(`
		Examples:
		` + ProjectName + ` reset-user -e dev
		` + ProjectName + ` reset-user -e staging
	`)

// AddEnv command related Info
const AddEnvCmdShortDesc string = "Add Environment to Config file"

var AddEnvCmdLongDesc = dedent.Dedent(`
		Add new environment and its related endpoints to the config file
	`)

var AddEnvCmdExamples = dedent.Dedent(`
		Examples:
		` + ProjectName + ` add-env -n production  --reg http://localhost/reg \
						--apim http://localhost/apim \
						--token http://localhost/token
	`)

// RemoveEnv command related Info
const RemoveEnvCmdShortDesc string = "Remove Environment from Config file"

var RemoveEnvCmdLongDesc = dedent.Dedent(`
		Remove Environment and its related endpoints from the config file
    `)

var RemoveEnvCmdExamples string = dedent.Dedent(`
		Examples:
		` + ProjectName + ` remove-env -n production
	`)

// Set command related Info
const SetCmdShortDesc string = "Set configuration"

var SetCmdLongDesc = dedent.Dedent(`
			Set configuration parameters. Use at least one of the following flags
				* --skip-tls-verfication <true>/<false>
				* --http-request-timeout <time-in-milli-seconds>
				* --export-directory <path-to-directory-where-apis-should-be-saved>
	`)

var SetCmdExamples = dedent.Dedent(`
			Examples:
			`+ ProjectName +` set --skip-tls-verification true \
				--http-request-timeout 3600 \
				--export-directory /home/user/exported-apis

			`+ ProjectName +` set --skip-tls-verification false \
				--http-request-timeout 5000 \
				--export-directory /media/user/apis

			`+ ProjectName +` set --skip-tls-verification false

			`+ ProjectName +` set --http-request-timeout 5000
	`)
