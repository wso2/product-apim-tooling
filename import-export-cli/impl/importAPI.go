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

// mergeAPI merges environmentParams to the API given in apiDirectory
// for now only Endpoints are merged
func mergeAPI(apiDirectory string, environmentParams *params.Environment) error {
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
	// extract environmentParams from file
	apiEndpointData, err := params.ExtractAPIEndpointConfig(api.Bytes())
	if err != nil {
		return err
	}

	// if endpointType field is not specified in the api_params.yaml, it will be considered as HTTP/REST
	if isEmpty(environmentParams.EndpointType) {
		environmentParams.EndpointType = utils.HttpRESTEndpointType
	}

	configData, err := setupMultipleEndpoints(environmentParams)
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

	// if the mutualSslCert field is defined in the api_params.yaml, the apiSecurity type should contain mutualssl
	if environmentParams.MutualSslCerts != nil {
		apiSecurity := api.Path("apiSecurity").Data()

		// if the apiSecurity field already exists in the api.yaml file
		if apiSecurity != nil {
			// if the apiSecurity field does not have mutualssl type, append it
			if !strings.Contains(apiSecurity.(string), utils.APISecurityMutualSsl) {
				apiSecurity = apiSecurity.(string) + "," + utils.APISecurityMutualSsl
			}
		} else {
			// if the apiSecurity field does not exist in the api.yaml file, assign the value as mutualssl
			apiSecurity = utils.APISecurityMutualSsl
		}
		// assign the apiSecurity field with the correct modified value to enable mutualssl
		if _, err := api.SetP(apiSecurity, "apiSecurity"); err != nil {
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

// setupMultipleEndpoints will set up the endpoints accordingly, for the applicable type
// @param environmentParams : Environment parameters from api_params.yaml
// @return configData as a byte array
// @return error
func setupMultipleEndpoints(environmentParams *params.Environment) ([]byte, error) {
	var configData []byte
	var err error

	// if the endpoint routing policy or the endpoints field is not specified
	if environmentParams.EndpointRoutingPolicy == "" && environmentParams.Endpoints == nil {
		// if endpoint type is Dynamic
		if environmentParams.EndpointType == utils.DynamicEndpointType {
			configData = []byte(utils.DynamicEndpointConfig)
		} else if environmentParams.EndpointType == utils.AwsLambdaEndpointType { // if endpoint type is AWS Lambda
			if environmentParams.AWSLambdaEndpoints == nil {
				return nil, errors.New("Please specify awsLambdaEndpoints field for " + environmentParams.Name + " and continue...")
			}
			if environmentParams.AWSLambdaEndpoints.AccessMethod == utils.AwsLambdaRoleSuppliedAccessMethod {
				environmentParams.AWSLambdaEndpoints.AccessMethod = utils.AwsLambdaRoleSuppliedAccessMethodForJSON
			}
			environmentParams.AWSLambdaEndpoints.EndpointType = utils.AwsLambdaEndpointTypeForJSON
			configData, err = json.Marshal(environmentParams.AWSLambdaEndpoints)
		} else {
			return nil, errors.New("Please specify the endpoint routing policy or the endpoints field for " + environmentParams.Name + " and continue...")
		}
	}

	// if endpoint type is HTTP/REST
	if environmentParams.EndpointType == utils.HttpRESTEndpointType || environmentParams.EndpointType == utils.HttpRESTEndpointTypeForJSON {
		environmentParams.EndpointType = utils.HttpRESTEndpointTypeForJSON

		// if the endpoint routing policy is not specified, but the endpoints field is specified, this is the usual scenario
		if environmentParams.EndpointRoutingPolicy == "" && environmentParams.Endpoints != nil {
			configData, err = json.Marshal(environmentParams.Endpoints)
		}

		// if the endpoint routing policy is specified as load balanced
		if environmentParams.EndpointRoutingPolicy == utils.LoadBalanceEndpointRoutingPolicy {
			if environmentParams.LoadBalanceEndpoints == nil {
				return nil, errors.New("Please specify loadBalanceEndpoints field for " + environmentParams.Name + " and continue...")
			}
			// The default class of the algorithm to be used should be set to RoundRobin
			environmentParams.LoadBalanceEndpoints.AlgorithmClassName = utils.LoadBalanceAlgorithmClass
			environmentParams.LoadBalanceEndpoints.EndpointType = utils.LoadBalanceEndpointTypeForJSON
			if environmentParams.LoadBalanceEndpoints.SessionManagement == utils.LoadBalanceSessionManagementTransport {
				// If the user has specified this as "transport", this should be converted to an empty string.
				// Otherwise APIM won't recognize this as "transport".
				environmentParams.LoadBalanceEndpoints.SessionManagement = ""
			}
			configData, err = json.Marshal(environmentParams.LoadBalanceEndpoints)
		}

		// if the endpoint routing policy is specified as failover
		if environmentParams.EndpointRoutingPolicy == utils.FailoverRoutingPolicy {
			if environmentParams.FailoverEndpoints == nil {
				return nil, errors.New("Please specify failoverEndpoints field for " + environmentParams.Name + " and continue...")
			}
			environmentParams.FailoverEndpoints.EndpointType = utils.FailoverRoutingPolicy
			environmentParams.FailoverEndpoints.Failover = true
			configData, err = json.Marshal(environmentParams.FailoverEndpoints)
		}
	}

	// if endpoint type is HTTP/SOAP
	if environmentParams.EndpointType == utils.HttpSOAPEndpointType {

		// if the endpoint routing policy is not specified, but the endpoints field is specified
		if environmentParams.EndpointRoutingPolicy == "" && environmentParams.Endpoints != nil {
			environmentParams.Endpoints.EndpointType = utils.HttpSOAPEndpointTypeForJSON
			configData, err = json.Marshal(environmentParams.Endpoints)
		}

		// if the endpoint routing policy is specified as load balanced
		if environmentParams.EndpointRoutingPolicy == utils.LoadBalanceEndpointRoutingPolicy {
			if environmentParams.LoadBalanceEndpoints == nil {
				return nil, errors.New("Please specify loadBalanceEndpoints field for " + environmentParams.Name + " and continue...")
			}
			// The default class of the algorithm to be used should be set to RoundRobin
			environmentParams.LoadBalanceEndpoints.AlgorithmClassName = utils.LoadBalanceAlgorithmClass
			environmentParams.LoadBalanceEndpoints.EndpointType = utils.LoadBalanceEndpointTypeForJSON
			for index := range environmentParams.LoadBalanceEndpoints.Production {
				environmentParams.LoadBalanceEndpoints.Production[index].EndpointType = utils.HttpSOAPEndpointTypeForJSON
			}
			for index := range environmentParams.LoadBalanceEndpoints.Sandbox {
				environmentParams.LoadBalanceEndpoints.Sandbox[index].EndpointType = utils.HttpSOAPEndpointTypeForJSON
			}
			if environmentParams.LoadBalanceEndpoints.SessionManagement == utils.LoadBalanceSessionManagementTransport {
				// If the user has specified this as "transport", this should be converted to an empty string.
				// Otherwise APIM won't recognize this as "transport".
				environmentParams.LoadBalanceEndpoints.SessionManagement = ""
			}
			configData, err = json.Marshal(environmentParams.LoadBalanceEndpoints)
		}

		// if the endpoint routing policy is specified as failover
		if environmentParams.EndpointRoutingPolicy == utils.FailoverRoutingPolicy {
			if environmentParams.FailoverEndpoints == nil {
				return nil, errors.New("Please specify failoverEndpoints field for " + environmentParams.Name + " and continue...")
			}
			environmentParams.FailoverEndpoints.Production.EndpointType = utils.HttpSOAPEndpointTypeForJSON
			environmentParams.FailoverEndpoints.Sandbox.EndpointType = utils.HttpSOAPEndpointTypeForJSON
			for index := range environmentParams.FailoverEndpoints.ProductionFailovers {
				environmentParams.FailoverEndpoints.ProductionFailovers[index].EndpointType = utils.HttpSOAPEndpointTypeForJSON
			}
			for index := range environmentParams.FailoverEndpoints.SandboxFailovers {
				environmentParams.FailoverEndpoints.SandboxFailovers[index].EndpointType = utils.HttpSOAPEndpointTypeForJSON
			}
			environmentParams.FailoverEndpoints.EndpointType = utils.FailoverRoutingPolicy
			environmentParams.FailoverEndpoints.Failover = true
			configData, err = json.Marshal(environmentParams.FailoverEndpoints)
		}
	}
	return configData, err
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

// generateEndpointCertificates for the API
func generateEndpointCertificates(importPath string, environment *params.Environment) error {
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

// generateMutualSslCertificates for the API
func generateMutualSslCertificates(importPath string, environment *params.Environment, importAPICmdPreserveProvider bool) error {
	// reading the definition to get API Name and the version
	apiInfo, _, err := GetAPIDefinition(importPath)
	if err != nil {
		return err
	}

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
		mutualSslCert.APIIdentifier.APIName = apiInfo.ID.APIName
		mutualSslCert.APIIdentifier.Version = apiInfo.ID.Version
		if !importAPICmdPreserveProvider {
			// if the preserve-provider flag is set to false
			// the currently logged in user's username should be assigned as the provider name
			mutualSslCert.APIIdentifier.ProviderName = utils.GetUsernameOfEnv(environment.Name, utils.EnvKeysAllFilePath)
		} else {
			// if the preserve-provider flag is set to true (the default behaviour)
			// the original provider should be taken from api.yaml and should be assigned as the provider name
			mutualSslCert.APIIdentifier.ProviderName = apiInfo.ID.ProviderName
		}
		mutualSslCerts = append(mutualSslCerts, mutualSslCert)
	}

	data, err := json.Marshal(mutualSslCerts)
	if err != nil {
		return err
	}

	yamlContent, err := utils.JsonToYaml(data)
	if err != nil {
		return err
	}

	// filepath to save mutualssl certs
	fp := filepath.Join(importPath, "Meta-information", "client_certificates.yaml")
	utils.Logln(utils.LogPrefixInfo+"Writing", fp)
	err = ioutil.WriteFile(fp, yamlContent, os.ModePerm)

	return err
}

// injectParamsToAPI injects ApiParams to API located in importPath using importEnvironment and returns the path to
// injected API location
func injectParamsToAPI(importPath, paramsPath, importEnvironment string, preserveProvider bool) error {
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
		err = mergeAPI(importPath, envParams)
		if err != nil {
			return err
		}

		err = generateEndpointCertificates(importPath, envParams)
		if err != nil {
			return err
		}

		// generate certificates for mutualssl, only if the field is specified
		if envParams.MutualSslCerts != nil {
			err = generateMutualSslCertificates(importPath, envParams, preserveProvider)
			if err != nil {
				return err
			}
		}
	}
	return nil
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
		err := injectParamsToAPI(apiFilePath, paramsPath, importEnvironment, preserveProvider)
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
	apiParams.SetP(string(data),"Certs")
	fmt.Println(apiParams)
	return nil
}

// injectMutualSslCertificates details for the API
func injectClientCertificates(importPath string, environment *params.Environment,apiParams *gabs.Container) error {

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
	apiParams.SetP(string(data),"MutualSslCerts")
	return nil
}
