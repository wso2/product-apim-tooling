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

// AmazonEcrRegistry represents Amazon ECR registry
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
		k8sUtils.K8sCreateSecretFromFile(k8sUtils.AwsCredentialsVolume, amazonEcrValues.credFile, "credentials")
	},
	Flags: Flags{
		RequiredFlags: &map[string]bool{k8sUtils.FlagBmRepository: true, k8sUtils.FlagBmKeyFile: true},
		OptionalFlags: &map[string]bool{},
	},
}

// readAmazonEcrInputs reads file path for amazon credential file
func readAmazonEcrInputs() (string, string) {
	isConfirm := false
	repository := ""
	credFile := ""
	var err error

	// validation regex for repository URI validation
	amazonRepositoryRegex := `\.amazonaws\.com\/.*$`

	for !isConfirm {
		repository, err = utils.ReadInputString("Enter Repository URI (<aws_account_id.dkr.ecr.region.amazonaws.com>/repository)", utils.Default{IsDefault: false}, amazonRepositoryRegex, true)
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

		fmt.Println("\nRepository     : " + repository)
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
		k8sUtils.Kubectl, k8sUtils.K8sCreate, k8sUtils.K8sConfigMap,
		k8sUtils.ConfigJsonVolume, "--from-file=config.json="+tempFile,
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

func init() {
	add(AmazonEcrRegistry)
}
