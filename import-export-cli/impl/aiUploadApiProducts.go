package impl

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/go-resty/resty/v2"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var apiProducts []utils.APIProduct

func AddAPIProductsToQueue(accessToken string, apiListQueue chan<- []map[string]interface{}) {
	if count == 0 {
		fmt.Println("No API Products available to be exported..!")
		return
	}
	for count > 0 {
		apiList := []map[string]interface{}{}
		for i := startingApiIndexFromList; i < len(apiProducts); i++ {
			apiPayload := GetAPIPayload(apiProducts[i], accessToken, CmdUploadEnvironment, true)
			if apiPayload != nil {
				apiList = append(apiList, apiPayload)
			}
		}
		atomic.AddInt32(&totalAPIs, int32(len(apiList)))
		if len(apiList) > 0 {
			apiListQueue <- apiList
		}
		apiListOffset += utils.MaxAPIsToExportOnce
		count, apiProducts, _ = GetAPIProductListFromEnv(accessToken, CmdUploadEnvironment, "", strconv.Itoa(utils.MaxAPIsToExportOnce)+"&offset="+strconv.Itoa(apiListOffset))
		startingApiIndexFromList = 0
	}
}

func GetAPIPayload(apiOrProduct interface{}, accessToken, cmdUploadEnvironment string, uploadProducts bool) map[string]interface{} {
	var name string
	var resp *resty.Response
	var err error

	if uploadProducts {
		api := apiOrProduct.(utils.APIProduct)
		if api.LifeCycleStatus != "PUBLISHED" && api.LifeCycleStatus != "PROTOTYPED" {
			return nil
		}
		resp, err = ExportAPIProductFromEnv(accessToken, api.Name, api.Version, "", api.Provider, "json", cmdUploadEnvironment, false, true)
		name = api.Name + "-" + api.Version
	} else {
		api := apiOrProduct.(utils.API)
		if api.LifeCycleStatus != "PUBLISHED" && api.LifeCycleStatus != "PROTOTYPED" {
			return nil
		}
		resp, err = ExportAPIFromEnv(accessToken, api.Name, api.Version, "", api.Provider, "json", cmdUploadEnvironment, true, false)
		name = api.Name + "-" + api.Version
	}

	if err != nil {
		utils.HandleErrorAndContinue("Error exporting API ", err)
		return nil
	}

	if resp.StatusCode() == http.StatusOK {
		zipReader, err := zip.NewReader(bytes.NewReader(resp.Body()), int64(len(resp.Body())))
		if err != nil {
			utils.HandleErrorAndContinue("Error reading zip file", err)
			return nil
		}

		apiPayload := map[string]interface{}{}

		for _, file := range zipReader.File {
			apiPayload = ReadZipFile(file, apiPayload, name)
			if apiPayload == nil {
				return nil
			}
		}
		return apiPayload
	}

	fmt.Println("Error exporting API: " + name + " Status: " + resp.Status())
	return nil
}

func ReadZipFile(file *zip.File, apiPayload map[string]interface{}, name string) map[string]interface{} {
	fileReader, err := file.Open()
	if err != nil {
		utils.HandleErrorAndContinue("Error while opening file", err)
		return nil
	}
	defer fileReader.Close()

	fileContents, err := ioutil.ReadAll(fileReader)
	if err != nil {
		utils.HandleErrorAndContinue("Error while reading file", err)
		return nil
	}

	if strings.HasSuffix(file.Name, name+"/api.json") || strings.HasSuffix(file.Name, name+"/api_product.json") {
		var jsonResp map[string]interface{}
		if err := json.Unmarshal(fileContents, &jsonResp); err != nil {
			utils.HandleErrorAndContinue("Error unmarshalling YAML content: %v\n", err)
			return nil
		}

		data, _ := jsonResp["data"].(map[string]interface{})

		if data["visibility"] != "PUBLIC" {
			return nil
		}

		apiPayload["uuid"] = data["id"]
		apiPayload["api_name"] = data["name"]
		apiPayload["version"] = data["version"]
		apiPayload["tenant_domain"] = Tenant
		if jsonResp["type"] == "api" {
			apiPayload["api_type"] = data["type"]
		} else {
			apiPayload["api_type"] = "APIPRODUCT"
		}

	} else if strings.HasSuffix(file.Name, name+"/Definitions/swagger.json") {
		var jsonResp map[string]interface{}
		if err := json.Unmarshal(fileContents, &jsonResp); err != nil {
			utils.HandleErrorAndContinue("Error unmarshalling YAML content: %v\n", err)
			return nil
		}
		info, _ := jsonResp["info"].(map[string]interface{})
		description := info["description"]
		apiPayload["description"] = description
		apiPayload["api_spec"] = string(fileContents)

	} else if strings.HasSuffix(file.Name, name+"/Definitions/schema.graphql") {
		apiPayload["description"] = ""
		apiPayload["sdl_schema"] = string(fileContents)

	} else if strings.HasSuffix(file.Name, name+"/Definitions/asyncapi.json") {
		var jsonResp map[string]interface{}
		if err := json.Unmarshal(fileContents, &jsonResp); err != nil {
			utils.HandleErrorAndContinue("Error unmarshalling YAML content: %v\n", err)
			return nil
		}
		info, _ := jsonResp["info"].(map[string]interface{})
		description := info["description"]
		apiPayload["description"] = description
		apiPayload["async_spec"] = string(fileContents)
	}

	return apiPayload
}
