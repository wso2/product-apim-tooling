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
const UnableToReadFileMsg string = "unable to read file '%v'"

const DoYouWantToContinueMsg_DefaultYes string = "Do you want to continue? [Y/n] "
const DoYouWantToContinueMsg_DefaultNo string = "Do you want to continue? [y/N] "

const RootCmdShortDesc string = "CLI for Importing and Exporting APIs"

var RootCmdLongDesc string = dedent.Dedent(`
		` + CLIName + ` is a Command Line Tool for Importing and Exporting APIs between different environments
		(Dev, Production, Staging, QA etc.)
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
const VersionCmdLiteral string = "version"
const VersionCmdShortDesc string = "Display Version on current `+ProjectName+`"

var VersionCmdLongDesc string = dedent.Dedent(`
		Display the current version of this command line tool
		`)

var VersionCmdExamples = dedent.Dedent(`
		` + ProjectName + ` ` + VersionCmdLiteral + `
		`)

// ImportAPI command related usage info
const ImportAPICmdLiteral string = "import-api"
const ImportAPICmdShortDesc string = "Import API"

var ImportAPICmdLongDesc string = "Import an API to an environment"

var ImportAPICmdExamples = dedent.Dedent(`
		Examples:
		` + ProjectName + ` ` + ImportAPICmdLiteral + ` -n TwitterAPI.zip -e dev
		` + ProjectName + ` ` + ImportAPICmdLiteral + ` -n FacebookAPI.zip -e production
	`)

// ExportAPI command related usage info
const ExportAPICmdLiteral string = "export-api"
const ExportAPICmdShortDesc string = "Export API"

var ExportAPICmdLongDesc string = "Export an API from an environment"

var ExportAPICmdExamples = dedent.Dedent(`
		Examples:
		` + ProjectName + ` ` + ExportAPICmdLiteral + ` -n TwitterAPI -v 1.0.0 -e dev
		` + ProjectName + ` ` + ExportAPICmdLiteral + ` -n FacebookAPI -v 2.1.0 -e production
	`)

// List command related usage Info

const ListCmdLiteral string = "list"
const ListCmdShortDesc string = "List APIs in an environment"

var ListCmdLongDesc string = dedent.Dedent(`
			Display a list containing all the APIs available in the environment specified by flag (--environment, -e)
	`)

var ListCmdExamples = dedent.Dedent(`
		Examples:
		` + ProjectName + ` ` + ListCmdLiteral + ` -e dev
		` + ProjectName + ` ` + ListCmdLiteral + ` -e staging
	`)

// ResetUser command related usage Info

const ResetUserCmdLiteral string = "reset-user"
const ResetUserCmdShortDesc string = "Reset user of an environment"

var ResetUserCmdLongDesc = dedent.Dedent(`
		Reset user data of a particular environment (Clear the entry in env_keys_all.yaml file)
	`)

var ResetUserCmdExamples = dedent.Dedent(`
		Examples:
		` + ProjectName + ` ` + ResetUserCmdLiteral + ` -e dev
		` + ProjectName + ` ` + ResetUserCmdLiteral + `reset-user -e staging
	`)

// AddEnv command related Info
const AddEnvCmdLiteral string = "add-env"
const AddEnvCmdShortDesc string = "Add Environment to Config file"

var AddEnvCmdLongDesc = dedent.Dedent(`
		Add new environment and its related endpoints to the config file
	`)

var AddEnvCmdExamples = dedent.Dedent(`
		Examples:
		` + ProjectName + ` ` + AddEnvCmdLiteral + ` -n production  --registration http://localhost/reg \
						--apim http://localhost/apim \
						--token http://localhost/token
	`)

// RemoveEnv command related Info
const RemoveEnvCmdLiteral string = "remove-env"
const RemoveEnvCmdShortDesc string = "Remove Environment from Config file"

var RemoveEnvCmdLongDesc = dedent.Dedent(`
		Remove Environment and its related endpoints from the config file
    `)

var RemoveEnvCmdExamples string = dedent.Dedent(`
		Examples:
		` + ProjectName + ` ` + RemoveEnvCmdLiteral + ` -n production
	`)

// Set command related Info
const SetCmdLiteral string = "set"
const SetCmdShortDesc string = "Set configuration"

var SetCmdLongDesc = dedent.Dedent(`
			Set configuration parameters. Use at least one of the following flags
				* --http-request-timeout <time-in-milli-seconds>
				* --export-directory <path-to-directory-where-apis-should-be-saved>
	`)

var SetCmdExamples = dedent.Dedent(`
			Examples:
			` + ProjectName + ` ` + SetCmdLiteral + ` --http-request-timeout 3600 \
								  --export-directory /home/user/exported-apis

			` + ProjectName + ` ` + SetCmdLiteral + ` --http-request-timeout 5000 \
								  --export-directory /media/user/apis

			` + ProjectName + ` ` + SetCmdLiteral + ` --http-request-timeout 5000
	`)
