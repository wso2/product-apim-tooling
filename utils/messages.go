package utils

import "github.com/renstrom/dedent"

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
const LogoutCmdShortDesc string = "Logout from the current session"

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
