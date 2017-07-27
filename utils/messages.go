package utils

import "github.com/renstrom/dedent"

const FindMoreInfoMsg string = "<FindMoreInfoMsg>"
const DoneMsg string = "Done!\n"
const NoAPIsFoundMsg = "No APIs found for the given query!"
const UnableToConnectMsg = "There was a problem connecting, please try again"
const InvalidCredentialsMsg string = "Invalid Credentials"
const UsernamePasswordEmptyMsg string = "Username and Password cannot be empty"
const InvalidExpiredRefreshTokenMsg string = "Your session has timed out"
const ErrorReadingResponseMsg string = "There was an error reading the response from the server"
const ErrorProcessingResponseMsg string = "There was an error processing the response from the server"

const DoYouWantToContinueMsg_DefaultYes string = "Do you want to continue? [Y/n] "
const DoYouWantToContinueMsg_DefaultNo string = "Do you want to continue? [y/N] "
const RunWSO2APIMInitToContinueMsg = "Run 'wso2apim init' to continue"

const RootCmdShortDesc string = "CLI for Importing and Exporting APIs"

var RootCmdLongDesc string = dedent.Dedent(`
		wso2api-cli is a CLI for Importing and Exporting APIs between different environments
		(Production, Staging, QA etc.)
		`)

// Init command related usage info
const InitCmdShortDesc string = "Initialize wso2apim-cli with your WSO2 credentials"

var InitCmdLongDesc = dedent.Dedent(`
		Initialize wso2apim-cli with your WSO2 credentials

		You need a WSO2 account to start using wso2apim-cli.
		Don't have one yet? Sign up at https://wso2.com/user/register
		`)

var InitCmdExamples = dedent.Dedent(`
		<InitCmdExamples>
		`)

// Logout command related usage info
const LogoutCmdShortDesc string = "Logout from current session"

var LogoutCmdLongDesc = dedent.Dedent(`
		<LogoutCmdLongDesc>
		`)

// Version command related usage info
const VersionCmdShortDesc string = "Display Version on current wso2apim-cli"

var VersionCmdLongDesc string = dedent.Dedent(`
		<VersionCmdLongDesc>
		`)

var VersionCmdExamples = dedent.Dedent(`
		<VersionCmdExamples>
		`)

// Config command related usage info
const ConfigCmdShortDesc string = "Configure wso2apim-cli"

var ConfigCmdLongDesc string = dedent.Dedent(`
		<ConfigCmdLongDesc>
		`)

var ConfigCmdExamples = dedent.Dedent(`
		<ConfigCmdExamples>
		`)

// ImportAPI command related usage info
const ImportAPICmdShortDesc string = "Import API"

var ImportAPICmdLongDesc string = "Import an API in a different environment to the current environment"

var ImportAPICmdExamples = dedent.Dedent(`
	wso2apim-cli importAPI -n TwitterAPI -v 1.0.0 -e dev
	wso2apim-cli importAPI -n FacebookAPI -v 2.1.0 -e production
	`)

// ExportAPI command related usage info
const ExportAPICmdShortDesc string = "Export API"

var ExportAPICmdLongDesc string = "Export an API in the current environment to a different environment"

var ExportAPICmdExamples = dedent.Dedent(`
	wso2apim-cli exportAPI -n TwitterAPI -v 1.0.0 -e dev
	wso2apim-cli exportAPI -n FacebookAPI -v 2.1.0 -e production
	`)
