package impl

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/go-resty/resty/v2"
	"github.com/wso2/product-apim-tooling/import-export-cli/credentials"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const (
	DefaultTenant = "carbon.super"
)

func AIUploadAPIs(credential credentials.Credential, cmdUploadEnvironment, onPremKey, endpointUrl string, uploadAll, uploadProducts bool) {

	CmdUploadEnvironment = cmdUploadEnvironment
	Credential = credential
	UploadAll = uploadAll
	UploadProducts = uploadProducts

	if !strings.Contains(credential.Username, "@") {
		Tenant = DefaultTenant
	} else {
		Tenant = strings.Split(credential.Username, "@")[1]
	}

	if endpointUrl != "" {
		Endpoint = endpointUrl
	}

	if onPremKey != "" {
		OnPremKey = onPremKey
	} else {
		OnPremKey = utils.OnPremKey
	}

	if OnPremKey == "" {
		fmt.Println("You have to provide your on prem key (that you generated for ai features) to do this operation.")
		os.Exit(1)
	}

	accessToken, err := credentials.GetOAuthAccessToken(credential, cmdUploadEnvironment)

	if err != nil {
		utils.HandleErrorAndExit("Error getting OAuth Tokens", err)
	}

	apiListQueue := make(chan []map[string]interface{}, 10)

	go ProduceAPIPayloads(accessToken, apiListQueue)

	numConsumers := utils.AIThreadCount
	var wg sync.WaitGroup
	for i := 0; i < numConsumers; i++ {
		wg.Add(1)
		go ConsumeAPIPayloads(apiListQueue, &wg)
	}

	wg.Wait()

	fmt.Printf("\nTotal number of public APIs present in the API Manager: %d\nTotal number of APIs successfully uploaded: %d\n\n", totalAPIs, uploadedAPIs)
}

func ProduceAPIPayloads(accessToken string, apiListQueue chan<- []map[string]interface{}) {
	ProcessAPIs(accessToken, apiListQueue)
	close(apiListQueue)
}

func ProcessAPIs(accessToken string, apiListQueue chan<- []map[string]interface{}) {
	apiListOffset = 0
	startingApiIndexFromList = 0
	if UploadAll {
		count, apis = getAPIList(Credential, CmdUploadEnvironment, "")
		AddAPIsToQueue(accessToken, apiListQueue)
		apiListOffset = 0
		count, apiProducts, _ = GetAPIProductListFromEnv(accessToken, CmdUploadEnvironment, "", strconv.Itoa(utils.MaxAPIsToExportOnce)+"&offset="+strconv.Itoa(apiListOffset))
		AddAPIProductsToQueue(accessToken, apiListQueue)
	} else if UploadProducts {
		count, apiProducts, _ = GetAPIProductListFromEnv(accessToken, CmdUploadEnvironment, "", strconv.Itoa(utils.MaxAPIsToExportOnce)+"&offset="+strconv.Itoa(apiListOffset))
		AddAPIProductsToQueue(accessToken, apiListQueue)
	} else {
		count, apis = getAPIList(Credential, CmdUploadEnvironment, "")
		AddAPIsToQueue(accessToken, apiListQueue)
	}
}

func AddAPIsToQueue(accessToken string, apiListQueue chan<- []map[string]interface{}) {
	if count == 0 {
		fmt.Println("No APIs available to be uploaded..!")
	} else {
		for count > 0 {
			apiList := []map[string]interface{}{}
			for i := startingApiIndexFromList; i < len(apis); i++ {
				apiPayload := GetAPIPayload(apis[i], accessToken, CmdUploadEnvironment, false)
				if apiPayload != nil {
					apiList = append(apiList, apiPayload)
				}
			}
			atomic.AddInt32(&totalAPIs, int32(len(apiList)))
			if len(apiList) > 0 {
				apiListQueue <- apiList
			}
			apiListOffset += utils.MaxAPIsToExportOnce
			count, apis = getAPIList(Credential, CmdUploadEnvironment, "")
			startingApiIndexFromList = 0
		}
	}
}

func ConsumeAPIPayloads(apiListQueue <-chan []map[string]interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

	for apiList := range apiListQueue {
		InvokePOSTRequest(apiList)
	}
}

func InvokePOSTRequest(apiList []map[string]interface{}) {
	fmt.Printf("Uploading %d APIs for tenant: %s\n", len(apiList), apiList[0]["tenant_domain"])
	payload, err := json.Marshal(map[string]interface{}{"apis": apiList})
	if err != nil {
		utils.HandleErrorAndContinue("Error in marshalling payload:", err)
		return
	}

	headers := make(map[string]string)
	headers["API-KEY"] = OnPremKey
	headers[utils.HeaderContentType] = utils.HeaderValueApplicationJSON

	var resp *resty.Response
	var uploadErr error

	for attempt := 1; attempt <= 2; attempt++ {
		resp, uploadErr = utils.InvokePOSTRequest(Endpoint+"/ai/spec-populator/bulk-upload", headers, payload)
		if uploadErr != nil {
			fmt.Printf("API upload failed (attempt %d). Reason: %v\n", attempt, uploadErr)
			continue
		}

		if resp.StatusCode() != 200 {
			fmt.Printf("Failed to upload %d APIs for tenant %s with status %d %s (attempt %d).\n", len(apiList), apiList[0]["tenant_domain"], resp.StatusCode(), resp.Body(), attempt)
			continue
		}

		jsonResp := map[string]map[string]int32{}

		err := json.Unmarshal(resp.Body(), &jsonResp)

		if err != nil {
			utils.HandleErrorAndContinue("Error in unmarshalling response:", err)
			continue
		}

		fmt.Printf("%d APIs uploaded successfully for tenant: %s (attempt %d)\n", len(apiList), apiList[0]["tenant_domain"], attempt)
		atomic.AddInt32(&uploadedAPIs, jsonResp["message"]["upsert_count"])
		break
	}

	if uploadErr != nil {
		utils.HandleErrorAndContinue("API upload failed after retry. Reason: ", uploadErr)
	}
}
