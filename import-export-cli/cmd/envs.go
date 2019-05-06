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
	"text/template"

	"github.com/renstrom/dedent"
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/formatter"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const (
	envsNameHeader                 = "NAME"
	envsPublisherEndpointHeader    = "PUBLISHER ENDPOINT"
	envsRegistrationEndpointHeader = "REGISTRATION ENDPOINT"
	envsTokenEndpointHeader        = "TOKEN ENDPOINT"
	envsAdminEndpointHeader        = "ADMIN ENDPOINT"
	envsImportExportEndpoint       = "IMPORT/EXPORT ENDPOINT"
	envsApiManagerEndpoint         = "API MANAGER ENDPOINT"
	envsApplicationEndpoint        = "APPLICATION ENDPOINT"

	defaulEnvsTableFormat = "table {{.Name}}\t{{.PublisherEndpoint}}\t{{.RegistrationEndpoint}}\t{{.TokenEndpoint}}"
)

var envsCmdFormat string

// envsCmd related info
const EnvsCmdLiteral = "envs"
const EnvsCmdShortDesc = "Display the list of environments"

var EnvsCmdLongDesc = dedent.Dedent(`
		Display a list of environments defined in '` + utils.MainConfigFileName + `' file
	`)

var EnvsCmdExamples = dedent.Dedent(`
		` + utils.ProjectName + ` list envs
	`)

// endpoint contains information about endpoint of API Manager
type endpoints struct {
	name                 string
	publisherEndpoint    string
	registrationEndpoint string
	tokenEndpoint        string
	adminEndpoint        string
	importExportEndpoint string
	apiManagerEndpoint   string
	applicationEndpoint  string
}

func newEndpointFromEnvEndpoints(name string, e utils.EnvEndpoints) *endpoints {
	return &endpoints{
		name:                 name,
		adminEndpoint:        e.AdminEndpoint,
		apiManagerEndpoint:   e.ApiManagerEndpoint,
		applicationEndpoint:  e.AppListEndpoint,
		importExportEndpoint: e.ApiImportExportEndpoint,
		publisherEndpoint:    e.ApiListEndpoint,
		registrationEndpoint: e.RegistrationEndpoint,
		tokenEndpoint:        e.TokenEndpoint,
	}
}

// Name of endpoint
func (e endpoints) Name() string {
	return e.name
}

// PublisherEndpoint
func (e endpoints) PublisherEndpoint() string {
	return e.publisherEndpoint
}

// RegistrationEndpoint
func (e endpoints) RegistrationEndpoint() string {
	return e.registrationEndpoint
}

// TokenEndpoint
func (e endpoints) TokenEndpoint() string {
	return e.tokenEndpoint
}

// ApplicationEndpoint
func (e endpoints) ApplicationEndpoint() string {
	return e.applicationEndpoint
}

// ApiManagerEndpoint
func (e endpoints) ApiManagerEndpoint() string {
	return e.apiManagerEndpoint
}

// ImportExportEndpoint
func (e endpoints) ImportExportEndpoint() string {
	return e.importExportEndpoint
}

// AdminEndpoint
func (e endpoints) AdminEndpoint() string {
	return e.adminEndpoint
}

// MarshalJSON returns marshaled methods
func (e *endpoints) MarshalJSON() ([]byte, error) {
	return formatter.MarshalJSON(e)
}

// envsCmd represents the envs command
var envsCmd = &cobra.Command{
	Use:   EnvsCmdLiteral,
	Short: EnvsCmdShortDesc,
	Long:  EnvsCmdLongDesc + EnvsCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + EnvsCmdLiteral + " called")
		envs := utils.GetMainConfigFromFile(utils.MainConfigFilePath).Environments
		printEnvs(envs, envsCmdFormat)
	},
}

func printEnvs(envData map[string]utils.EnvEndpoints, format string) {
	if format == "" {
		format = defaulEnvsTableFormat
	}

	// create api context with standard output
	envsContext := formatter.NewContext(os.Stdout, format)

	// create a new renderer function which iterate collection
	renderer := func(w io.Writer, t *template.Template) error {
		for name, endpointDef := range envData {
			if err := t.Execute(w, newEndpointFromEnvEndpoints(name, endpointDef)); err != nil {
				return err
			}
			_, _ = w.Write([]byte{'\n'})
		}
		return nil
	}

	// headers for table
	envsTableHeaders := map[string]string{
		"Name":                 envsNameHeader,
		"PublisherEndpoint":    envsPublisherEndpointHeader,
		"RegistrationEndpoint": envsRegistrationEndpointHeader,
		"TokenEndpoint":        envsTokenEndpointHeader,
		"AdminEndpoint":        envsAdminEndpointHeader,
		"ImportExportEndpoint": envsImportExportEndpoint,
		"ApiManagerEndpoint":   envsApiManagerEndpoint,
		"ApplicationEndpoint":  envsApplicationEndpoint,
	}

	// execute context
	if err := envsContext.Write(renderer, envsTableHeaders); err != nil {
		fmt.Println("Error executing template:", err.Error())
	}
}

func init() {
	ListCmd.AddCommand(envsCmd)
	envsCmd.Flags().StringVarP(&envsCmdFormat, "format", "", "", "Pretty-print "+
		"environments using go templates")
}
