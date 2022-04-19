package impl

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/json-iterator/go"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

// ExportThrottlingPolicyFromEnv function is used with export policy rate-limiting command
func ExportThrottlingPolicyFromEnv(accessToken, exportEnvironment string, exportThrottlePolicyName string, exportThrottlePolicyType string, exportFormat string) (*resty.Response, error) {
	adminEndpoint := utils.GetAdminEndpointOfEnv(exportEnvironment, utils.MainConfigFilePath)
	return exportThrottlePolicy(adminEndpoint, accessToken, exportThrottlePolicyName, exportThrottlePolicyType, exportFormat)
}

func exportThrottlePolicy(adminEndpoint, accessToken string, ThrottlePolicyName string, ThrottlePolicyType string, exportFormat string) (*resty.Response, error) {
	var PolicyType string
	adminEndpoint = utils.AppendSlashToString(adminEndpoint)
	ThrottlePolicyResource := "throttling/policies/export?"
	switch ThrottlePolicyType {
	case "sub":
		PolicyType = "subscription"
	case "app":
		PolicyType = "application"
	case "advanced":
		PolicyType = "advanced"
	case "custom":
		PolicyType = "custom"
	}

	query := `name=` + ThrottlePolicyName + `&type=` + PolicyType + `&format=` + exportFormat
	ThrottlePolicyResource += query
	url := adminEndpoint + ThrottlePolicyResource
	utils.Logln(utils.LogPrefixInfo+"ExportThrottlingPolicy: URL:", url)
	headers := make(map[string]string)
	headers[utils.HeaderAuthorization] = utils.HeaderValueAuthBearerPrefix + " " + accessToken
	resp, err := utils.InvokeGETRequest(url, headers)

	if err != nil {
		return nil, err
	}
	return resp, nil
}

func WriteThrottlePolicyToFile(ExportLocationPath string, resp *resty.Response, ExportFormat string, runningExportThrottlePolicyCommand bool) {
	err := utils.CreateDirIfNotExist(ExportLocationPath)
	if err != nil {
		utils.HandleErrorAndExit("Error creating dir to store zip archives: "+ExportLocationPath, err)
	}

	fileName, marshaledData := ResolveThrottlePolicy(ExportFormat, resp)

	_, _ = ThrottlingPolicyWrite(ExportLocationPath, fileName, marshaledData)

	if runningExportThrottlePolicyCommand {
		fmt.Println("Successfully exported Throttling Policy!")
		fmt.Println("Find the exported Throttling Policies at " + utils.AppendSlashToString(ExportLocationPath) + fileName)
	}
}

func ResolveThrottlePolicy(exportThrottlePolicyFormat string, resp *resty.Response) (string, []byte) {

	var subscriptionPolicy utils.ExportThrottlePolicySubscription
	var applicationPolicy utils.ExportThrottlePolicyApplication
	var advancedPolicy utils.ExportThrottlePolicyAdvanced
	var customPolicy utils.ExportThrottlePolicyCustom
	var ExportThrottlingPolicyResponse utils.ExportThrottlePolicyGeneral
	var marshaledData []byte
	err := json.Unmarshal(resp.Body(), &ExportThrottlingPolicyResponse)
	if err != nil {
		utils.HandleErrorAndExit("Error unmarshalling response data", err)
	}

	throttlingPolicyType := ExportThrottlingPolicyResponse.Subtype
	policyName := ExportThrottlingPolicyResponse.Data.PolicyName
	fileName := throttlingPolicyType + `-` + policyName

	switch throttlingPolicyType {
	case "subcription policy":
		err = json.Unmarshal(resp.Body(), &subscriptionPolicy)
		if err != nil {
			utils.HandleErrorAndExit("Error unmarshalling policy data", err)
		}

		marshaledData, err = jsoniter.MarshalIndent(subscriptionPolicy, "", " ")
		if err != nil {
			utils.HandleErrorAndExit("Error marshaling policy content", err)
		}
	case "application policy":
		err = json.Unmarshal(resp.Body(), &applicationPolicy)
		if err != nil {
			utils.HandleErrorAndExit("Error unmarshalling policy data", err)
		}

		marshaledData, err = jsoniter.MarshalIndent(applicationPolicy, "", " ")
		if err != nil {
			utils.HandleErrorAndExit("Error marshaling policy content", err)
		}
	case "advanced policy":
		err = json.Unmarshal(resp.Body(), &advancedPolicy)
		if err != nil {
			utils.HandleErrorAndExit("Error unmarshalling policy data", err)
		}

		marshaledData, err = jsoniter.MarshalIndent(advancedPolicy, "", " ")
		if err != nil {
			utils.HandleErrorAndExit("Error marshaling policy content", err)
		}
	case "custom rule":
		err = json.Unmarshal(resp.Body(), &customPolicy)
		if err != nil {
			utils.HandleErrorAndExit("Error unmarshalling policy data", err)
		}

		marshaledData, err = jsoniter.MarshalIndent(customPolicy, "", " ")
		if err != nil {
			utils.HandleErrorAndExit("Error marshaling policy content", err)
		}
	}
	fmt.Println(exportThrottlePolicyFormat)
	if exportThrottlePolicyFormat == utils.DefaultExportFormat {
		m := yaml.MapSlice{}
		err := yaml.Unmarshal(marshaledData, &m)
		fileName += ".yaml"
		if err != nil {
			utils.HandleErrorAndExit("Error marshaling policy content", err)
		}
		marshaledData, _ = yaml.Marshal(m)
	} else {
		fileName += ".json"

	}
	return fileName, marshaledData
}

func ThrottlingPolicyWrite(FilePath string, Filename string, marshaledData []byte) (string, error) {

	Filename = filepath.Join(FilePath, Filename)
	err := ioutil.WriteFile(Filename, marshaledData, 0644)
	if err != nil {
		utils.HandleErrorAndExit("Error writing file", err)
	}
	return FilePath, err
}
