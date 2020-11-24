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

package impl

import (
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/wso2/product-apim-tooling/import-export-cli/specs/params"

	"github.com/mitchellh/go-homedir"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	v2 "github.com/wso2/product-apim-tooling/import-export-cli/specs/v2"

	"github.com/Jeffail/gabs"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var (
	reAPIName = regexp.MustCompile(`[~!@#;:%^*()+={}|\\<>"',&/$]`)
)

// extractAPIDefinition extracts API information from jsonContent
func extractAPIDefinition(jsonContent []byte) (*v2.APIDefinition, error) {
	api := &v2.APIDefinition{}
	err := json.Unmarshal(jsonContent, &api)
	if err != nil {
		return nil, err
	}

	return api, nil
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

// resolveAPIParamsPath resolves api_params.yaml path
// First it will look at AbsolutePath of the import path (the last directory)
// If not found it will look at current working directory
// If a path is provided search ends looking up at that path
func resolveAPIParamsPath(importPath, paramPath string) (string, error) {
	utils.Logln(utils.LogPrefixInfo + "Scanning for parameters file")
	if paramPath == utils.ParamFileAPI {
		// look in importpath
		if stat, err := os.Stat(importPath); err == nil && stat.IsDir() {
			loc := filepath.Join(importPath, utils.ParamFileAPI)
			utils.Logln(utils.LogPrefixInfo+"Scanning for", loc)
			if info, err := os.Stat(loc); err == nil && !info.IsDir() {
				// found api_params.yml in the importpath
				return loc, nil
			}
		}

		// look in the basepath of importPath
		base := filepath.Dir(importPath)
		fp := filepath.Join(base, utils.ParamFileAPI)
		utils.Logln(utils.LogPrefixInfo+"Scanning for", fp)
		if info, err := os.Stat(fp); err == nil && !info.IsDir() {
			// found api_params.yml in the base path
			return fp, nil
		}

		// look in the current working directory
		wd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		utils.Logln(utils.LogPrefixInfo+"Scanning for", wd)
		fp = filepath.Join(wd, utils.ParamFileAPI)
		if info, err := os.Stat(fp); err == nil && !info.IsDir() {
			// found api_params.yml in the cwd
			return fp, nil
		}

		// no luck, it means paramPath is missing
		return "", fmt.Errorf("could not find %s. Please check %s exists in basepath of "+
			"import location or current working directory", utils.ParamFileAPI, utils.ParamFileAPI)
	} else {
		if info, err := os.Stat(paramPath); err == nil && !info.IsDir() {
			return paramPath, nil
		}
		return "", fmt.Errorf("could not find %s", paramPath)
	}
}

// resolveYamlOrJSON for a given filepath.
// first it will look for the yaml file, if not will fallback for json
// give filename without extension so resolver will resolve for file
// fn is resolved filename, jsonContent is file as a json object, error if anything wrong happen(or both files does not exists)
func resolveYamlOrJSON(filename string) (string, []byte, error) {
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

// isEmpty returns true when a given string is empty
func isEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func preProcessAPI(apiDirectory string) error {
	dirty := false
	apiPath, jsonData, err := resolveYamlOrJSON(filepath.Join(apiDirectory, "Meta-information", "api"))
	if err != nil {
		return err
	}
	utils.Logln(utils.LogPrefixInfo+"Loading API definition from: ", apiPath)

	api, err := gabs.ParseJSON(jsonData)
	if err != nil {
		return err
	}

	// preprocess endpoint config
	if !api.Exists("endpointConfig") {
		dirty = true
		conf, err := gabs.ParseJSON([]byte(`{"endpoint_type":"http"}`))
		if err != nil {
			return err
		}

		if api.Exists("productionUrl") {
			_, err = conf.SetP(api.Path("productionUrl").Data(), "production_endpoints.url")
			if err != nil {
				return err
			}
			_, err = conf.SetP("null", "production_endpoints.config")
			if err != nil {
				return err
			}
			_ = api.Delete("productionUrl")
		}
		if api.Exists("sandboxUrl") {
			_, err = conf.SetP(api.Path("sandboxUrl").Data(), "sandbox_endpoints.url")
			if err != nil {
				return err
			}
			_, err = conf.SetP("null", "sandbox_endpoints.config")
			if err != nil {
				return err
			}
			_ = api.Delete("sandboxUrl")
		}

		_, err = api.SetP(conf.String(), "endpointConfig")
		if err != nil {
			return err
		}
	}

	if dirty {
		yamlAPIPath := filepath.Join(apiDirectory, "Meta-information", "api.yaml")
		utils.Logln(utils.LogPrefixInfo+"Writing preprocessed API to:", yamlAPIPath)
		content, err := utils.JsonToYaml(api.Bytes())
		if err != nil {
			return err
		}
		// write this to disk
		err = ioutil.WriteFile(yamlAPIPath, content, 0644)
		if err != nil {
			return err
		}
	}

	return nil
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

func populateAPIWithDefaults(def *v2.APIDefinition) (dirty bool) {
	dirty = false
	def.Context = strings.ReplaceAll(def.Context, " ", "")
	if def.ContextTemplate == "" {
		if !strings.Contains(def.Context, "{version}") {
			def.ContextTemplate = path.Clean(def.Context + "/{version}")
			def.Context = strings.ReplaceAll(def.ContextTemplate, "{version}", def.ID.Version)
		} else {
			def.Context = path.Clean(def.Context)
			def.ContextTemplate = def.Context
			def.Context = strings.ReplaceAll(def.Context, "{version}", def.ID.Version)
		}
		dirty = true
	}
	if def.Tags == nil {
		def.Tags = []string{}
		dirty = true
	}
	if def.URITemplates == nil {
		def.URITemplates = []v2.URITemplates{}
		dirty = true
	}
	if def.Implementation == "" {
		def.Implementation = "ENDPOINT"
		dirty = true
	}
	if def.KeyManagers == nil || len(def.KeyManagers) == 0 {
		def.KeyManagers = []string{"all"}
		dirty = true
	}
	return
}

// validateAPIDefinition validates an API against basic rules
func validateAPIDefinition(def *v2.APIDefinition) error {
	utils.Logln(utils.LogPrefixInfo + "Validating API")
	if isEmpty(def.ID.APIName) {
		return errors.New("apiName is required")
	}
	if reAPIName.MatchString(def.ID.APIName) {
		return errors.New(`apiName contains one or more illegal characters (~!@#;:%^*()+={}|\\<>"',&\/$)`)
	}

	if isEmpty(def.ID.Version) {
		return errors.New("version is required")
	}
	if isEmpty(def.Context) {
		return errors.New("context is required")
	}
	if isEmpty(def.ContextTemplate) {
		return errors.New("contextTemplate is required")
	}
	if !strings.HasPrefix(def.Context, "/") {
		return errors.New("context should begin with a /")
	}
	if !strings.HasPrefix(def.ContextTemplate, "/") {
		return errors.New("contextTemplate should begin with a /")
	}
	return nil
}

// importAPI imports an API to the API manager
func importAPI(endpoint, filePath, accessToken string, extraParams map[string]string) error {
	resp, err := ExecuteNewFileUploadRequest(endpoint, extraParams, "file",
		filePath, accessToken)
	if err != nil {
		utils.Logln(utils.LogPrefixError, err)
		return err
	}
	if resp.StatusCode() == http.StatusCreated || resp.StatusCode() == http.StatusOK {
		// 201 Created or 200 OK
		fmt.Println("Successfully imported API.")
		return nil
	} else {
		// We have an HTTP error
		fmt.Println("Error importing API.")
		fmt.Println("Status: " + resp.Status())
		fmt.Println("Response:", resp)
		return errors.New(resp.Status())
	}
}

// ImportAPIToEnv function is used with import-api command
func ImportAPIToEnv(accessOAuthToken, importEnvironment, importPath, apiParamsPath string, importAPIUpdate, preserveProvider,
	importAPISkipCleanup bool) error {
	adminEndpoint := utils.GetAdminEndpointOfEnv(importEnvironment, utils.MainConfigFilePath)
	return ImportAPI(accessOAuthToken, adminEndpoint, importEnvironment, importPath, apiParamsPath, importAPIUpdate,
		preserveProvider, importAPISkipCleanup)
}

// ImportAPI function is used with import-api command
func ImportAPI(accessOAuthToken, adminEndpoint, importEnvironment, importPath, apiParamsPath string, importAPIUpdate, preserveProvider,
	importAPISkipCleanup bool) error {
	exportDirectory := filepath.Join(utils.ExportDirectory, utils.ExportedApisDirName)
	resolvedAPIFilePath, err := resolveImportFilePath(importPath, exportDirectory)
	if err != nil {
		return err
	}
	utils.Logln(utils.LogPrefixInfo+"API Location:", resolvedAPIFilePath)

	utils.Logln(utils.LogPrefixInfo + "Creating workspace")
	tmpPath, err := utils.GetTempCloneFromDirOrZip(resolvedAPIFilePath)
	if err != nil {
		return err
	}
	defer func() {
		if importAPISkipCleanup {
			utils.Logln(utils.LogPrefixInfo+"Leaving", tmpPath)
			return
		}
		utils.Logln(utils.LogPrefixInfo+"Deleting", tmpPath)
		err := os.RemoveAll(tmpPath)
		if err != nil {
			utils.Logln(utils.LogPrefixError + err.Error())
		}
	}()
	apiFilePath := tmpPath

	utils.Logln(utils.LogPrefixInfo + "Substituting environment variables in API files...")
	err = replaceEnvVariables(apiFilePath)
	if err != nil {
		return err
	}

	utils.Logln(utils.LogPrefixInfo + "Pre Processing API...")
	err = preProcessAPI(apiFilePath)
	if err != nil {
		return err
	}

	utils.Logln(utils.LogPrefixInfo + "Attempting to inject parameters to the API from api_params.yaml (if exists)")
	paramsPath, err := resolveAPIParamsPath(resolvedAPIFilePath, apiParamsPath)
	if err != nil && apiParamsPath != utils.ParamFileAPI && apiParamsPath != "" {
		return err
	}
	if paramsPath != "" {
		//Reading API params file and populate api.yaml
		err := handleCustomizedParameters(apiFilePath, paramsPath, importEnvironment, preserveProvider)
		if err != nil {
			return err
		}
	}

	// Get API info
	apiInfo, originalContent, err := GetAPIDefinition(apiFilePath)
	if err != nil {
		return err
	}
	// Fill with defaults
	if populateAPIWithDefaults(apiInfo) {
		utils.Logln(utils.LogPrefixInfo + "API is populated with defaults")
		// api is dirty, write it to disk
		buf, err := json.Marshal(apiInfo)
		if err != nil {
			return err
		}

		newContent, err := gabs.ParseJSON(buf)
		if err != nil {
			return err
		}
		originalContent, err := gabs.ParseJSON(originalContent)
		if err != nil {
			return err
		}
		result, err := utils.MergeJSON(originalContent.Bytes(), newContent.Bytes())
		if err != nil {
			return err
		}

		yamlContent, err := utils.JsonToYaml(result)
		if err != nil {
			return err
		}
		p := filepath.Join(apiFilePath, "Meta-information", "api.yaml")
		utils.Logln(utils.LogPrefixInfo+"Writing", p)

		err = ioutil.WriteFile(p, yamlContent, 0644)
		if err != nil {
			return err
		}
	}
	// validate definition
	if err = validateAPIDefinition(apiInfo); err != nil {
		return err
	}

	// if apiFilePath contains a directory, zip it. Otherwise, leave it as it is.
	apiFilePath, err, cleanupFunc := utils.CreateZipFileFromProject(apiFilePath, importAPISkipCleanup)
	if err != nil {
		return err
	}

	//cleanup the temporary artifacts once consuming the zip file
	if cleanupFunc != nil {
		defer cleanupFunc()
	}

	updateAPI := false
	if importAPIUpdate {
		// check for API existence
		id, err := GetAPIId(accessOAuthToken, importEnvironment, apiInfo.ID.APIName, apiInfo.ID.Version,
			apiInfo.ID.ProviderName)
		if err != nil {
			return err
		}

		if id == "" {
			utils.Logln("The specified API was not found.")
			utils.Logln("Creating: %s %s\n", apiInfo.ID.APIName, apiInfo.ID.Version)
		} else {
			utils.Logln("Existing API found, attempting to update it...")
			utils.Logln("API ID:", id)
			updateAPI = true
		}
	}
	extraParams := map[string]string{}
	adminEndpoint += "/import/api"
	if updateAPI {
		adminEndpoint += "?overwrite=" + strconv.FormatBool(true) + "&preserveProvider=" +
			strconv.FormatBool(preserveProvider)
	} else {
		adminEndpoint += "?preserveProvider=" + strconv.FormatBool(preserveProvider)
	}
	utils.Logln(utils.LogPrefixInfo + "Import URL: " + adminEndpoint)

	err = importAPI(adminEndpoint, apiFilePath, accessOAuthToken, extraParams)
	return err
}

// injectParamsToAPI injects ApiParams to API located in importPath using importEnvironment and returns the path to
// injected API location
func handleCustomizedParameters(importPath, paramsPath, importEnvironment string, preserveProvider bool) error {
	utils.Logln(utils.LogPrefixInfo+"Loading parameters from", paramsPath)
	apiParams, err := params.LoadApiParamsFromFile(paramsPath)
	if err != nil {
		return err
	}
	// check whether import environment is included in api configuration
	envParams := apiParams.GetEnv(importEnvironment)
	if envParams == nil {
		utils.Logln(utils.LogPrefixInfo + "Using default values as the environment is not present in api_param.yaml file")
	} else {
		//If environment parameters are present in parameter file
		err = handleEnvParams(importPath, envParams)
		if err != nil {
			return err
		}
	}
	return nil
}

// injectEndpointCertificates details for the API
func injectEndpointCerts(importPath string, environment *params.Environment, apiParams *gabs.Container) error {

	var certs []params.Cert

	if len(environment.Certs) == 0 {
		return nil
	}

	//Read certificate list provided in user api-params file
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

	//Inject cert details to api_params file
	apiParams.SetP(string(data), "Certs")
	return nil
}

// injectMutualSslCertificates details for the API
func injectClientCertificates(importPath string, environment *params.Environment, apiParams *gabs.Container) error {

	var mutualSslCerts []params.MutualSslCert

	if len(environment.MutualSslCerts) == 0 {
		return nil
	}

	for _, mutualSslCert := range environment.MutualSslCerts {
		// read cert
		p, err := resolveCertPath(importPath, mutualSslCert.Path)
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
		mutualSslCert.Certificate = enc
		mutualSslCert.APIIdentifier.APIName = ""
		mutualSslCert.APIIdentifier.Version = ""
		mutualSslCert.APIIdentifier.ProviderName = ""

		mutualSslCerts = append(mutualSslCerts, mutualSslCert)
	}

	data, err := json.Marshal(mutualSslCerts)
	if err != nil {
		return err
	}

	//Inject mutualSSL cert details to  api_params file
	apiParams.SetP(string(data), "MutualSslCerts")
	return nil
}

//Process env params and create a temp env_parmas.yaml in temp artifact
func handleEnvParams(apiDirectory string, environmentParams *params.Environment) error {
	// read api from Meta-information
	apiPath := filepath.Join(apiDirectory, "Meta-information", "api")
	utils.Logln(utils.LogPrefixInfo + "Reading API definition: ")
	fp, jsonContent, err := resolveYamlOrJSON(apiPath)
	if err != nil {
		return err
	}
	utils.Logln(utils.LogPrefixInfo+"Loaded definition from:", fp)
	api, err := gabs.ParseJSON(jsonContent)
	if err != nil {
		return err
	}

	envParamsJson, err := json.Marshal(environmentParams)

	if err == nil {
		s := string(envParamsJson)
		fmt.Println(s)
	}

	apiPath = filepath.Join(apiDirectory, "Meta-information", "api.yaml")
	var apiParamsPath string
	apiParamsPath = filepath.Join(apiDirectory, "Meta-information", "env_params.yaml")
	utils.Logln(utils.LogPrefixInfo+"Writing merged API to:", apiPath)

	// write this to disk
	content, err := utils.JsonToYaml(api.Bytes())
	if err != nil {
		return err
	}

	apiParams, err := gabs.ParseJSON(envParamsJson)

	err = injectEndpointCerts(apiDirectory, environmentParams, apiParams)
	if err != nil {
		return err
	}

	err = injectClientCertificates(apiDirectory, environmentParams, apiParams)
	if err != nil {
		return err
	}

	paramsContent, err := utils.JsonToYaml(apiParams.Bytes())
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(apiPath, content, 0644)
	err = ioutil.WriteFile(apiParamsPath, paramsContent, 0644)
	if err != nil {
		return err
	}
	return nil
}