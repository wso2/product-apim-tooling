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
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/wso2/product-apim-tooling/import-export-cli/impl"
	"github.com/wso2/product-apim-tooling/import-export-cli/specs/params"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Jeffail/gabs"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var (
	importAPIFile                string
	importEnvironment            string
	importAPICmdPreserveProvider bool
	importAPIUpdate              bool
	importAPIParamsFile          string
	importAPISkipCleanup         bool
)

const (
	// ImportAPI command related usage info
	DefaultAPIMParamsFileName = "api_params.yaml"
	importAPICmdLiteral       = "import-api"
	importAPICmdShortDesc     = "Import API"
	importAPICmdLongDesc      = "Import an API to an environment"
)

const importAPICmdExamples = utils.ProjectName + ` ` + importAPICmdLiteral + ` -f qa/TwitterAPI.zip -e dev
` + utils.ProjectName + ` ` + importAPICmdLiteral + ` -f staging/FacebookAPI.zip -e production
` + utils.ProjectName + ` ` + importAPICmdLiteral + ` -f ~/myapi -e production --update
` + utils.ProjectName + ` ` + importAPICmdLiteral + ` -f ~/myapi -e production --update
NOTE: Both the flags (--file (-f) and --environment (-e)) are mandatory`

// ImportAPICmd represents the importAPI command
var ImportAPICmd = &cobra.Command{
	Use: importAPICmdLiteral + " --file <PATH_TO_API> --environment " +
		"<ENVIRONMENT>",
	Short:   importAPICmdShortDesc,
	Long:    importAPICmdLongDesc,
	Example: importAPICmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + importAPICmdLiteral + " called")
		credential, err := getCredentials(importEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error getting credentials", err)
		}
		accessOAuthToken, err := credentials.GetOAuthAccessToken(credential, importEnvironment)
		if err != nil {
			utils.HandleErrorAndExit("Error while getting an access token for importing API", err)
		}
		err = impl.ImportAPI(accessOAuthToken, importEnvironment, importAPIFile, importAPIParamsFile, importAPIUpdate,
			importAPICmdPreserveProvider, importAPISkipCleanup)
		if err != nil {
			utils.HandleErrorAndExit("Error importing API", err)
			return
		}
	},
}

// mergeAPI merges environmentParams to the API given in apiDirectory
// for now only Endpoints are merged
func mergeAPI(apiDirectory string, environmentParams *params.Environment) error {
	// read api from Meta-information
	apiPath := filepath.Join(apiDirectory, "Meta-information", "api")
	utils.Logln(utils.LogPrefixInfo + "Reading API definition: ")
	fp, jsonContent, err := resolveYamlOrJson(apiPath)
	if err != nil {
		return err
	}
	utils.Logln(utils.LogPrefixInfo+"Loaded definition from:", fp)
	api, err := gabs.ParseJSON(jsonContent)
	if err != nil {
		return err
	}
	// extract environmentParams from file
	apiEndpointData, err := params.ExtractAPIEndpointConfig(api.Bytes())
	if err != nil {
		return err
	}

	configData, err := json.Marshal(environmentParams.Endpoints)
	if err != nil {
		return err
	}

	mergedAPIEndpoints, err := utils.MergeJSON([]byte(apiEndpointData), configData)
	if err != nil {
		return err
	}

	utils.Logln(utils.LogPrefixInfo + "Merging API")
	// replace original endpointConfig with merged version
	if _, err := api.SetP(string(mergedAPIEndpoints), "endpointConfig"); err != nil {
		return err
	}

	// replace original GatewayEnvironments only if they are present in api-params file
	if environmentParams.GatewayEnvironments != nil {
		if _, err := api.SetP(environmentParams.GatewayEnvironments, "environments"); err != nil {
			return err
		}
	}

	// Handle security parameters in api_params.yaml
	err = handleSecurityEndpointsParams(environmentParams.Security, api)
	if err != nil {
		return err
	}

	apiPath = filepath.Join(apiDirectory, "Meta-information", "api.yaml")
	utils.Logln(utils.LogPrefixInfo+"Writing merged API to:", apiPath)
	// write this to disk
	content, err := utils.JsonToYaml(api.Bytes())
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(apiPath, content, 0644)
	if err != nil {
		return err
	}
	return nil
}

// Handle security parameters in api_params.yaml
// @param envSecurityEndpointParams : Environment security endpoint parameters from api_params.yaml
// @param api : Parameters from api.yaml
// @return error
func handleSecurityEndpointsParams(envSecurityEndpointParams *params.SecurityData, api *gabs.Container) error {
	// If the user has set (either true or false) the enabled field under security in api_params.yaml, the
	// following code should be executed. (if not set, the security endpoint settings will be made
	// according to the api.yaml file as usually)
	// In Go, irrespective of whether a boolean value is "" or false, it will contain false by default.
	// That is why, here a string comparison was  made since strings can have both "" and "false"
	if envSecurityEndpointParams != nil && envSecurityEndpointParams.Enabled != "" {
		// Convert the string enabled to boolean
		boolEnabled, err := strconv.ParseBool(envSecurityEndpointParams.Enabled)
		if err != nil {
			return err
		}
		if _, err := api.SetP(boolEnabled, "endpointSecured"); err != nil {
			return err
		}
		// If endpoint security is enabled
		if boolEnabled {
			// Set the security endpoint parameters when the enabled field is set to true
			err := setSecurityEndpointsParams(envSecurityEndpointParams, api)
			if err != nil {
				return err
			}
		} else {
			// If endpoint security is not enabled, the username and password should be empty.
			// (otherwise the security will be enabled after importing the API, considering there are values
			// for username and passwords)
			if _, err := api.SetP("", "endpointUTUsername"); err != nil {
				return err
			}
			if _, err := api.SetP("", "endpointUTPassword"); err != nil {
				return err
			}
		}
	}
	return nil
}

// Set the security endpoint parameters when the enabled field is set to true
// @param envSecurityEndpointParams : Environment security endpoint parameters from api_params.yaml
// @param api : Parameters from api.yaml
// @return error
func setSecurityEndpointsParams(envSecurityEndpointParams *params.SecurityData, api *gabs.Container) error {
	// Check whether the username, password and type fields have set in api_params.yaml
	if envSecurityEndpointParams.Username == "" {
		return errors.New("You have enabled endpoint security but the username is not found in the api_params.yaml")
	} else if envSecurityEndpointParams.Password == "" {
		return errors.New("You have enabled endpoint security but the password is not found in the api_params.yaml")
	} else if envSecurityEndpointParams.Type == "" {
		return errors.New("You have enabled endpoint security but the type is not found in the api_params.yaml")
	} else {
		// Override the username in api.yaml with the value in api_params.yaml
		if _, err := api.SetP(envSecurityEndpointParams.Username, "endpointUTUsername"); err != nil {
			return err
		}
		// Override the password in api.yaml with the value in api_params.yaml
		if _, err := api.SetP(envSecurityEndpointParams.Password, "endpointUTPassword"); err != nil {
			return err
		}
		// Set the fields in api.yaml according to the type field in api_params.yaml
		err := setEndpointSecurityType(envSecurityEndpointParams, api)
		if err != nil {
			return err
		}
	}
	return nil
}

// Set the fields in api.yaml according to the type field in api_params.yaml
// @param envSecurityEndpointParams : Environment security endpoint parameters from api_params.yaml
// @param api : Parameters from api.yaml
// @return error
func setEndpointSecurityType(envSecurityEndpointParams *params.SecurityData, api *gabs.Container) error {
	// Check whether the type is either basic or digest
	if envSecurityEndpointParams.Type == "digest" {
		if _, err := api.SetP(true, "endpointAuthDigest"); err != nil {
			return err
		}
	} else if envSecurityEndpointParams.Type == "basic" {
		if _, err := api.SetP(false, "endpointAuthDigest"); err != nil {
			return err
		}
	} else {
		// If the type is not either basic or digest, return an error
		return errors.New("Invalid endpoint security type found in the api_params.yaml. Should be either basic or digest")
	}
	return nil
}

// resolveImportFilePath resolves the archive/directory for import
// First will resolve in given path, if not found will try to load from exported directory
func resolveImportFilePath(file, defaultExportDirectory string) (string, error) {
	// check current path
	utils.Logln(utils.LogPrefixInfo + "Resolving for API path...")
	if _, err := os.Stat(file); os.IsNotExist(err) {
		// if the file not in given path it might be inside exported directory
		utils.Logln(utils.LogPrefixInfo+"Looking for API in", defaultExportDirectory)
		file = filepath.Join(defaultExportDirectory, file)
		if _, err := os.Stat(file); os.IsNotExist(err) {
			return "", err
		}
	}
	absPath, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}

	return absPath, nil
}

// extractArchive extracts the API and give the path.
// In API Manager archive there is a directory in the root which contains the API
// this function returns it appended to the destination path
func extractArchive(src, dest string) (string, error) {
	files, err := utils.Unzip(src, dest)
	if err != nil {
		return "", err
	}
	if len(files) == 0 {
		return "", fmt.Errorf("invalid API archive")
	}
	r := strings.TrimPrefix(files[0], src)
	return filepath.Join(dest, strings.Split(filepath.Clean(r), string(os.PathSeparator))[0]), nil
}


func getTempApiDirectory(file string) (string, error) {
	fileIsDir := false
	// create a temp directory
	tmpDir, err := ioutil.TempDir("", "apim")
	if err != nil {
		_ = os.RemoveAll(tmpDir)
		return "", err
	}

	if info, err := os.Stat(file); err == nil {
		fileIsDir = info.IsDir()
	} else {
		return "", err
	}
	if fileIsDir {
		// copy dir to a temp location
		utils.Logln(utils.LogPrefixInfo+"Copying from", file, "to", tmpDir)
		dest := filepath.Join(tmpDir, filepath.Base(file))
		err = utils.CopyDir(file, dest)
		if err != nil {
			return "", err
		}
		return dest, nil
	} else {
		// try to extract archive
		utils.Logln(utils.LogPrefixInfo+"Extracting", file, "to", tmpDir)
		finalPath, err := extractArchive(file, tmpDir)
		if err != nil {
			return "", err
		}
		return finalPath, nil
	}
}

// resolveYamlOrJson for a given filepath.
// first it will look for the yaml file, if not will fallback for json
// give filename without extension so resolver will resolve for file
// fn is resolved filename, jsonContent is file as a json object, error if anything wrong happen(or both files does not exists)
func resolveYamlOrJson(filename string) (string, []byte, error) {
	// lookup for yaml
	yamlFp := filename + ".yaml"
	if info, err := os.Stat(yamlFp); err == nil && !info.IsDir() {
		utils.Logln(utils.LogPrefixInfo+"Loading", yamlFp)
		// read it
		fn := yamlFp
		yamlContent, err := ioutil.ReadFile(fn)
		if err != nil {
			return "", nil, err
		}
		// load it as yaml
		jsonContent, err := utils.YamlToJson(yamlContent)
		if err != nil {
			return "", nil, err
		}
		return fn, jsonContent, nil
	}

	jsonFp := filename + ".json"
	if info, err := os.Stat(jsonFp); err == nil && !info.IsDir() {
		utils.Logln(utils.LogPrefixInfo+"Loading", jsonFp)
		// read it
		fn := jsonFp
		jsonContent, err := ioutil.ReadFile(fn)
		if err != nil {
			return "", nil, err
		}
		return fn, jsonContent, nil
	}

	return "", nil, fmt.Errorf("%s was not found as a YAML or JSON", filename)
}

func resolveCertPath(importPath, p string) (string, error) {
	// look in importPath
	utils.Logln(utils.LogPrefixInfo+"Resolving for", p)
	certfile := filepath.Join(importPath, p)
	utils.Logln(utils.LogPrefixInfo + "Looking in project directory")
	if info, err := os.Stat(certfile); err == nil && !info.IsDir() {
		// found file return it
		return certfile, nil
	}

	utils.Logln(utils.LogPrefixInfo + "Looking for absolute path")
	// try for an absolute path
	p, err := homedir.Expand(filepath.Clean(p))
	if err != nil {
		return "", err
	}

	if p != "" {
		return p, nil
	}

	return "", fmt.Errorf("%s not found", p)
}

// generateCertificates for the API
func generateCertificates(importPath string, environment *params.Environment) error {
	var certs []params.Cert

	if len(environment.Certs) == 0 {
		return nil
	}

	for _, cert := range environment.Certs {
		// read cert
		p, err := resolveCertPath(importPath, cert.Path)
		if err != nil {
			return err
		}
		pubPEMData, err := ioutil.ReadFile(p)
		if err != nil {
			return err
		}
		// get cert
		block, _ := pem.Decode(pubPEMData)
		enc := credentials.Base64Encode(string(block.Bytes))
		cert.Certificate = enc
		certs = append(certs, cert)
	}

	data, err := json.Marshal(certs)
	if err != nil {
		return err
	}

	yamlContent, err := utils.JsonToYaml(data)
	if err != nil {
		return err
	}

	// filepath to save certs
	fp := filepath.Join(importPath, "Meta-information", "endpoint_certificates.yaml")
	utils.Logln(utils.LogPrefixInfo+"Writing", fp)
	err = ioutil.WriteFile(fp, yamlContent, os.ModePerm)

	return err
}

// injectParamsToAPI injects ApiParams to API located in importPath using importEnvironment and returns the path to
// injected API location
func injectParamsToAPI(importPath, paramsPath, importEnvironment string) error {
	utils.Logln(utils.LogPrefixInfo+"Loading parameters from", paramsPath)
	apiParams, err := params.LoadApiParamsFromFile(paramsPath)
	if err != nil {
		return err
	}
	// check whether import environment is included in api configuration
	envParams := apiParams.GetEnv(importEnvironment)
	if envParams == nil {
		fmt.Println("Using default values as the environment is not present in api_param.yaml file")
	} else {
		//If environment parameters are present in parameter file
		err = mergeAPI(importPath, envParams)
		if err != nil {
			return err
		}

		err = generateCertificates(importPath, envParams)
		if err != nil {
			return err
		}
	}

	return nil
}

// isEmpty returns true when a given string is empty
func isEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}


// Substitutes environment variables in the project files.
func replaceEnvVariables(apiFilePath string) error {
	for _, replacePath := range utils.EnvReplaceFilePaths {
		absFile := filepath.Join(apiFilePath, replacePath)
		// check if the path exists. If exists, proceed with processing. Otherwise, continue with the next items
		if fi, err := os.Stat(absFile); err != nil {
			if !os.IsNotExist(err) {
				return err
			}
		} else {
			switch mode := fi.Mode(); {
			case mode.IsDir():
				utils.Logln(utils.LogPrefixInfo+"Substituting env variables of files in folder path: ", absFile)
				err = utils.EnvSubstituteInFolder(absFile)
			case mode.IsRegular():
				utils.Logln(utils.LogPrefixInfo+"Substituting env of file: ", absFile)
				err = utils.EnvSubstituteInFile(absFile)
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// init using Cobra
func init() {
	RootCmd.AddCommand(ImportAPICmd)
	ImportAPICmd.Flags().StringVarP(&importAPIFile, "file", "f", "",
		"Name of the API to be imported")
	ImportAPICmd.Flags().StringVarP(&importEnvironment, "environment", "e",
		"", "Environment from the which the API should be imported")
	ImportAPICmd.Flags().BoolVar(&importAPICmdPreserveProvider, "preserve-provider", true,
		"Preserve existing provider of API after importing")
	ImportAPICmd.Flags().BoolVarP(&importAPIUpdate, "update", "", false, "Update an "+
		"existing API or create a new API")
	ImportAPICmd.Flags().StringVarP(&importAPIParamsFile, "params", "", DefaultAPIMParamsFileName,
		"Provide a API Manager params file")
	ImportAPICmd.Flags().BoolVarP(&importAPISkipCleanup, "skipCleanup", "", false, "Leave "+
		"all temporary files created during import process")
	// Mark required flags
	_ = ImportAPICmd.MarkFlagRequired("environment")
	_ = ImportAPICmd.MarkFlagRequired("file")
}
