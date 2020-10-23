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
	apiProductIdHeader       = "ID"
	apiProductNameHeader     = "NAME"
	apiProductContextHeader  = "CONTEXT"
	apiProductProviderHeader = "PROVIDER"
	apiProductStatusHeader   = "STATUS"

	defaultApiProductTableFormat = "table {{.Id}}\t{{.Name}}\t{{.Context}}\t{{.LifeCycleStatus}}\t{{.Provider}}"
)

var listApiProductsCmdEnvironment string
var listApiProductsCmdFormat string
var listApiProductsCmdQuery string
var listApiProductsCmdLimit string

// apisCmd related info
const apiProductsCmdLiteral = "api-products"
const apiProductsCmdShortDesc = "Display a list of API Products in an environment"

const apiProductsCmdLongDesc = `Display a list of API Products in the environment specified by the flag --environment, -e`

var apiProductsCmdExamples = utils.ProjectName + ` ` + listCmdLiteral + ` ` + apiProductsCmdLiteral + ` -e dev
` + utils.ProjectName + ` ` + listCmdLiteral + ` ` + apiProductsCmdLiteral + ` -e dev -q provider:devops
` + utils.ProjectName + ` ` + listCmdLiteral + ` ` + apiProductsCmdLiteral + ` -e prod -q provider:admin context:/myproduct
` + utils.ProjectName + ` ` + listCmdLiteral + ` ` + apisCmdLiteral + ` -e prod -l 25
` + utils.ProjectName + ` ` + listCmdLiteral + ` ` + apiProductsCmdLiteral + ` -e staging
NOTE: The flag (--environment (-e)) is mandatory`

// apiProductsCmd represents the api-products command
var apiProductsCmd = &cobra.Command{
	Use:     apiProductsCmdLiteral,
	Short:   apiProductsCmdShortDesc,
	Long:    apiProductsCmdLongDesc,
	Example: apiProductsCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + apiProductsCmdLiteral + " called")
		cred, err := GetCredentials(listApiProductsCmdEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		//Since other flags does not use args[], query flag will own this
		if len(args) != 0 && listApiProductsCmdQuery != "" {
			for _, argument := range args {
				listApiProductsCmdQuery += " " + argument
			}
		}
		executeApiProductsCmd(cred)
	},
}

// apiProduct holds information about an API Product for outputting
type apiProduct struct {
	id              string
	name            string
	context         string
	provider        string
	lifeCycleStatus string
}

// creates a new API Product from utils.API
func newApiProductDefinitionFromAPI(a utils.APIProduct) *apiProduct {
	return &apiProduct{a.ID, a.Name, a.Context, a.Provider, a.LifeCycleStatus}
}

// Id of API Product
func (a apiProduct) Id() string {
	return a.id
}

// Name of API Product
func (a apiProduct) Name() string {
	return a.name
}

// Context of API Product
func (a apiProduct) Context() string {
	return a.context
}

// Lifecycle Status of API Product
func (a apiProduct) LifeCycleStatus() string {
	return a.lifeCycleStatus
}

// Provider of API Product
func (a apiProduct) Provider() string {
	return a.provider
}

// MarshalJSON marshals apiProduct using custom marshaller which uses methods instead of fields
func (a *apiProduct) MarshalJSON() ([]byte, error) {
	return formatter.MarshalJSON(a)
}

func executeApiProductsCmd(credential credentials.Credential) {
	accessToken, err := credentials.GetOAuthAccessToken(credential, listApiProductsCmdEnvironment)
	if err != nil {
		utils.Logln(utils.LogPrefixError + "calling 'list' " + err.Error())
		utils.HandleErrorAndExit("Error calling '"+apiProductsCmdLiteral+"'", err)
	}

	// Unified Search endpoint from the config file to search API Products
	_, apiProducts, err := impl.GetAPIProductListFromEnv(accessToken, listApiProductsCmdEnvironment, listApiProductsCmdQuery,
		listApiProductsCmdLimit)
	if err == nil {
		printAPIProducts(apiProducts, listApiProductsCmdFormat)
	} else {
		utils.Logln(utils.LogPrefixError+"Getting List of API Products", err)
	}
}

// printAPIProducts
func printAPIProducts(apiProducts []utils.APIProduct, format string) {
	if format == "" {
		format = defaultApiProductTableFormat
	}
	// create API Product context with standard output
	apiProductContext := formatter.NewContext(os.Stdout, format)

	// create a new renderer function which iterate collection
	renderer := func(w io.Writer, t *template.Template) error {
		for _, a := range apiProducts {
			if err := t.Execute(w, newApiProductDefinitionFromAPI(a)); err != nil {
				return err
			}
			_, _ = w.Write([]byte{'\n'})
		}
		return nil
	}

	// headers for table
	apiProductTableHeaders := map[string]string{
		"Id":              apiProductIdHeader,
		"Name":            apiProductNameHeader,
		"Context":         apiProductContextHeader,
		"LifeCycleStatus": apiProductStatusHeader,
		"Provider":        apiProductProviderHeader,
	}

	// execute context
	if err := apiProductContext.Write(renderer, apiProductTableHeaders); err != nil {
		fmt.Println("Error executing template:", err.Error())
	}
}

func init() {
	ListCmd.AddCommand(apiProductsCmd)

	apiProductsCmd.Flags().StringVarP(&listApiProductsCmdEnvironment, "environment", "e",
		"", "Environment to be searched")
	apiProductsCmd.Flags().StringVarP(&listApiProductsCmdQuery, "query", "q",
		"", "Query pattern")
	apiProductsCmd.Flags().StringVarP(&listApiProductsCmdLimit, "limit", "l",
		strconv.Itoa(utils.DefaultApiProductsDisplayLimit), "Maximum number of API Products to return")
	apiProductsCmd.Flags().StringVarP(&listApiProductsCmdFormat, "format", "", "", "Pretty-print API Products "+
		"using Go Templates. Use \"{{ jsonPretty . }}\" to list all fields")
	_ = apiProductsCmd.MarkFlagRequired("environment")
}
