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

package cmd

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"text/template"

	"github.com/wso2/product-apim-tooling/import-export-cli/impl"

	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/formatter"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const (
	appIdHeader      = "ID"
	appNameHeader    = "NAME"
	appOwnerHeader   = "OWNER"
	appStatusHeader  = "STATUS"
	appGroupIdHeader = "GROUP ID"

	defaultAppTableFormat = "table {{.Id}}\t{{.Name}}\t{{.Owner}}\t{{.Status}}\t{{.GroupId}}"
)

var listAppsCmdEnvironment string
var listAppsCmdAppOwner string
var listAppsCmdFormat string
var listAppsCmdLimit string
var defaultAppsOwner string

// appsCmd related info
const appsCmdLiteral = "apps"
const appsCmdShortDesc = "Display a list of Applications in an environment specific to an owner"

// app contains information about util.Application
type app struct {
	id      string
	name    string
	owner   string
	status  string
	groupId string
}

// creates a new app definition from utils.Application
func newAppDefinitionFromApplication(a utils.Application) *app {
	return &app{a.ID, a.Name, a.Owner, a.Status, a.GroupID}
}

// Id of application
func (a app) Id() string {
	return a.id
}

// Name of application
func (a app) Name() string {
	return a.name
}

// Owner of application
func (a app) Owner() string {
	return a.owner
}

// Status of application
func (a app) Status() string {
	return a.status
}

// GroupId of application
func (a app) GroupId() string {
	return a.groupId
}

// MarshalJSON marshals api using custom marshaller which uses methods instead of fields
func (a *app) MarshalJSON() ([]byte, error) {
	return formatter.MarshalJSON(a)
}

const appsCmdLongDesc = "Display a list of Applications of the user in the environment specified by the flag --environment, -e"

const appsCmdExamples = utils.ProjectName + ` ` + listCmdLiteral + ` ` + appsCmdLiteral + ` -e dev 
` + utils.ProjectName + ` ` + listCmdLiteral + ` ` + appsCmdLiteral + ` -e dev -o sampleUser
` + utils.ProjectName + ` ` + listCmdLiteral + ` ` + appsCmdLiteral + ` -e prod -o sampleUser
` + utils.ProjectName + ` ` + listCmdLiteral + ` ` + appsCmdLiteral + ` -e staging -o sampleUser
` + utils.ProjectName + ` ` + listCmdLiteral + ` ` + appsCmdLiteral + ` -e dev -l 40
NOTE: The flag (--environment (-e)) is mandatory`

// appsCmd represents the apps command
var appsCmd = &cobra.Command{
	Use:     appsCmdLiteral,
	Short:   appsCmdShortDesc,
	Long:    appsCmdLongDesc,
	Example: appsCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + appsCmdLiteral + " called")
		cred, err := GetCredentials(listAppsCmdEnvironment)
		defaultAppsOwner = cred.Username
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		executeAppsCmd(cred, listAppsCmdAppOwner)
	},
}

func executeAppsCmd(credential credentials.Credential, appOwner string) {
	accessToken, err := credentials.GetOAuthAccessToken(credential, listAppsCmdEnvironment)
	if err != nil {
		utils.Logln(utils.LogPrefixError + "calling 'list' " + err.Error())
		utils.HandleErrorAndExit("Error calling '"+appsCmdLiteral+"'", err)
	}

	_, apps, err := impl.GetApplicationListFromEnv(accessToken, listAppsCmdEnvironment, appOwner, listAppsCmdLimit)

	if err == nil {
		// Printing the list of available Applications
		printApps(apps, listAppsCmdFormat)
	} else {
		utils.Logln(utils.LogPrefixError+"Getting List of Applications", err)
	}
}

func printApps(apps []utils.Application, format string) {
	if format == "" {
		format = defaultAppTableFormat
	}
	// create new app context with standard output
	appContext := formatter.NewContext(os.Stdout, format)

	// create a new renderer function which iterate collection of apps
	renderer := func(w io.Writer, t *template.Template) error {
		for _, a := range apps {
			if err := t.Execute(w, newAppDefinitionFromApplication(a)); err != nil {
				return err
			}
			// write a new line after executing template
			_, _ = w.Write([]byte{'\n'})
		}
		return nil
	}

	// headers for table
	appTableHeaders := map[string]string{
		"Id":      appIdHeader,
		"Name":    appNameHeader,
		"Status":  appStatusHeader,
		"Owner":   appOwnerHeader,
		"GroupId": appGroupIdHeader,
	}

	// execute context
	if err := appContext.Write(renderer, appTableHeaders); err != nil {
		fmt.Println("Error executing template:", err.Error())
	}
}

func init() {
	ListCmd.AddCommand(appsCmd)

	appsCmd.Flags().StringVarP(&listAppsCmdAppOwner, "owner", "o", defaultAppsOwner,
		"Owner of the Application")
	appsCmd.Flags().StringVarP(&listAppsCmdEnvironment, "environment", "e",
		"", "Environment to be searched")
	appsCmd.Flags().StringVarP(&listAppsCmdLimit, "limit", "l",
		strconv.Itoa(utils.DefaultAppsDisplayLimit), "Maximum number of applications to return")
	appsCmd.Flags().StringVarP(&listAppsCmdFormat, "format", "", "", "Pretty-print output"+
		"using Go templates. Use \"{{jsonPretty .}}\" to list all fields")
	_ = appsCmd.MarkFlagRequired("environment")
}
