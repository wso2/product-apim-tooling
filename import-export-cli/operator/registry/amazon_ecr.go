package registry

import (
	"fmt"
	k8sUtils "github.com/wso2/product-apim-tooling/import-export-cli/operator/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"os"
	"path/filepath"
	"strings"
)

var amazoneEcrRepo = new(string)

var amazonEcrValues = struct {
	repository string
	credFile   string
}{}

var AmazonEcrRegistry = &Registry{
	Name:       "AMAZON_ECR",
	Caption:    "Amazon ECR",
	Repository: amazoneEcrRepo,
	Option:     2,
	Read: func() {
		repository, credFile := readAmazonEcrInputs()
		amazonEcrValues.repository = repository
		amazonEcrValues.credFile = credFile
		*amazoneEcrRepo = repository
	},
	Run: func() {
		createAmazonEcrConfig()
		createAmazonEcrCred(amazonEcrValues.credFile)
	},
}

// readAmazonEcrInputs reads file path for amazon credential file
func readAmazonEcrInputs() (string, string) {
	isConfirm := false
	repository := ""
	credFile := ""
	var err error

	amazonRepositoryRegex := `\.amazonaws\.com\/.*$` //TODO: renuka make this regex more specif with finding repo syntax

	for !isConfirm {
		repository, err = utils.ReadInputString("Enter repository name (<aws_account_id.dkr.ecr.region.amazonaws.com>/repository)", utils.Default{IsDefault: false}, amazonRepositoryRegex, true)
		if err != nil {
			utils.HandleErrorAndExit("Error reading DockerHub repository name from user", err)
		}

		defaultLocation, err := os.UserHomeDir()
		if err == nil {
			defaultLocation = filepath.Join(defaultLocation, ".aws", "credentials")
		} // else ignore and make defaultLocation = ""

		credFile, err = utils.ReadInput("Amazon credential file", utils.Default{Value: defaultLocation, IsDefault: true}, utils.IsFileExist, "Invalid file", true)
		if err != nil {
			utils.HandleErrorAndExit("Error reading amazon credential file from user", err)
		}

		fmt.Println("")
		fmt.Println("Repository     : " + repository)
		fmt.Println("Credential File: " + credFile)

		isConfirmStr, err := utils.ReadInputString("Confirm configurations", utils.Default{Value: "Y", IsDefault: true}, "", false)
		if err != nil {
			utils.HandleErrorAndExit("Error reading user input Confirmation", err)
		}

		isConfirmStr = strings.ToUpper(isConfirmStr)
		isConfirm = isConfirmStr == "Y" || isConfirmStr == "YES"
	}

	return repository, credFile
}

// createAmazonEcrConfig creates K8S secret with credentials for Amazon ECR
func createAmazonEcrConfig() {
	configJson := `{ "credsStore": "ecr-login" }`

	tempFile, err := utils.CreateTempFile("config-*.json", []byte(configJson))
	if err != nil {
		utils.HandleErrorAndExit("Error writing configs to temporary file", err)
	}
	defer os.Remove(tempFile)

	// render config map
	configMap, err := k8sUtils.GetCommandOutput(
		utils.Kubectl, utils.Create, utils.K8sConfigMap,
		utils.ConfigJsonVolume, "--from-file=config.json="+tempFile,
		"--dry-run", "-o", "yaml",
	)
	if err != nil {
		utils.HandleErrorAndExit("Error creating docker config for Amazon ECR", err)
	}

	// apply config map
	if err = k8sUtils.K8sApplyFromStdin(configMap); err != nil {
		utils.HandleErrorAndExit("Error creating docker config for Amazon ECR", err)
	}
}

func createAmazonEcrCred(credFile string) {
	// render secret
	secret, err := k8sUtils.GetCommandOutput(
		utils.Kubectl, utils.Create, utils.K8sSecret, "generic",
		k8sUtils.AwsCredentialsVolume, "--from-file=credentials="+credFile,
		"--dry-run", "-o", "yaml",
	)
	if err != nil {
		utils.HandleErrorAndExit("Error creating secret for Amazon ECR", err)
	}

	// apply secret
	if err = k8sUtils.K8sApplyFromStdin(secret); err != nil {
		utils.HandleErrorAndExit("Error creating secret for Amazon ECR", err)
	}
}

func init() {
	add(AmazonEcrRegistry)
}
