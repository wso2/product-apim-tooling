package k8s

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	wso2v1alpha2 "github.com/wso2/k8s-api-operator/api-operator/pkg/apis/wso2/v1alpha2"
	"github.com/wso2/product-apim-tooling/import-export-cli/box"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	"os"
	"path/filepath"
	"strings"
)

var genDeploymentDirDestination string
var genDeploymentDirSource string

type metaData struct {
	Name     string `yaml:"name"`
	Revision string `yaml:"revision"`
	Version  string `yaml:"version"`
}

var apiMetaData *metaData

// GetEnvsCmd related info
const GenDeploymentDirCmdLiteral = "deployment-dir"
const GenDeploymentDirCmdShortDesc = "Generate a sample deployment directory"

const GenDeploymentDirCmdLongDesc = `Generate a sample deployment directory based on the provided source artifact`

const GenDeploymentDirCmdExamples = utils.ProjectName + ` ` + K8sCmdLiteral + ` ` + K8sGenCmdLiteral + ` ` + GenDeploymentDirCmdLiteral + ` ` +
	`-s  ~/PizzaShackAPI_1.0.0.zip
` + utils.ProjectName + ` ` + K8sCmdLiteral + ` ` + K8sGenCmdLiteral + ` ` + GenDeploymentDirCmdLiteral + ` ` +
	`-s  ~/PizzaShackAPI_1.0.0.zip` + ` ` + ` -d /home/Deployment_repo/Dev`

// directories to be created
var directories = []string{
	"certificates",
}

// createDeploymentContentDirectories will create directories in current working directory
func createDeploymentContentDirectories(name string) error {
	for _, directory := range directories {
		directoryPath := filepath.Join(name, filepath.FromSlash(directory))
		utils.Logln(utils.LogPrefixInfo + "Creating directory " + directoryPath)
		err := os.MkdirAll(directoryPath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

// executeGenDeploymentDirCmd will run gen deployment-dir command
func executeGenDeploymentDirCmd() error {
	var deploymentDirParent, deploymentDirName, sourceDirectoryPath, tempDirPath string

	// Check the validity of destination path when it is given if not given use the working directory
	if genDeploymentDirDestination != "" {
		err := os.MkdirAll(genDeploymentDirDestination, os.ModePerm)
		if err != nil {
			return err
		}
		p, err := filepath.Abs(genDeploymentDirDestination)
		if err != nil {
			return err
		}
		deploymentDirParent = p
	} else {
		pwd, err := os.Getwd()
		if err != nil {
			return err
		}
		deploymentDirParent = pwd
	}

	// Check whether the source is existed in the given location
	if _, err := os.Stat(genDeploymentDirSource); os.IsNotExist(err) {
		utils.HandleErrorAndContinue("Error retrieving the source file from the given path "+sourceDirectoryPath+" ", err)
	}
	// Get the source artifact name
	deploymentDirName = filepath.Base(genDeploymentDirSource)
	if info, err := os.Stat(genDeploymentDirSource); err == nil && !info.IsDir() {
		// if artifact is given as zip remove the ".zip" suffix to get the name for deployment directory
		deploymentDirName = strings.Trim(deploymentDirName, utils.ZipFileSuffix)
		//extract zip to a temp directory

		tempDirPath := os.TempDir()
		path, err := utils.Unzip(genDeploymentDirSource, tempDirPath)
		if err != nil {
			return err
		}
		sourceDirectoryPath = tempDirPath + "/" + path[0]
	} else {
		sourceDirectoryPath = genDeploymentDirSource
	}

	deploymentDirPath, err := filepath.Abs(filepath.Join(deploymentDirParent, deploymentDirName))
	if err != nil {
		return err
	}

	//Create the deployment directory
	err = utils.CreateDir(deploymentDirPath)
	if err != nil {
		return err
	}

	// Copy *_meta.yaml file from source to deployment directory based on the artifact type
	files, err := ioutil.ReadDir(sourceDirectoryPath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	var metaDataFileFound bool = false

	//add api_crd.yaml file to deployment directory
	//get API definition from file
	apiConfigMapData, _ := box.Get("/kubernetes_resources/api_cr.yaml")
	apiCrd := &wso2v1alpha2.API{}
	errUnmarshal := yaml.Unmarshal(apiConfigMapData, apiCrd)
	if errUnmarshal != nil {
		utils.HandleErrorAndExit("Error unmarshal api configmap into struct ", errUnmarshal)
	}

	for _, file := range files {
		fileName := file.Name()
		// if project artifact is a API project
		if strings.EqualFold(fileName, utils.MetaFileAPI) {
			metaDataFileFound = true
			err := utils.CopyFile(filepath.Join(sourceDirectoryPath, fileName), filepath.Join(deploymentDirPath, utils.MetaFileAPI))
			if err != nil {
				utils.HandleErrorAndExit("Cannot copy metadata file from the source directory ", err)
			}
			metaDataYamlFile, err := ioutil.ReadFile(filepath.Join(sourceDirectoryPath, fileName))
			if err != nil {
				utils.HandleErrorAndExit("Cannot read the meta file", err)
			}
			errUnmarshal := yaml.Unmarshal(metaDataYamlFile, &apiMetaData)
			if errUnmarshal != nil {
				utils.HandleErrorAndExit("Error unmarshal api configmap into struct ", errUnmarshal)
			}
			apiName := utils.GetValidK8sName(apiMetaData.Name)
			apiCrd.Name = apiName
			apiCrd.Spec.SwaggerConfigMapName = fmt.Sprintf("%v-cm", apiName)
			apiCrd.Spec.ParamsValues = fmt.Sprintf("%v-params", apiName)
			apiCrd.Spec.CertsValues = fmt.Sprintf("%v-certs", apiName)
			break
		} else if strings.EqualFold(fileName, utils.MetaFileAPIProduct) { // if project artifact is a APIProduct project
			metaDataFileFound = true
			err := utils.CopyFile(filepath.Join(sourceDirectoryPath, fileName), filepath.Join(deploymentDirPath, utils.MetaFileAPIProduct))
			if err != nil {
				utils.HandleErrorAndExit("Cannot copy metadata file from the source directory ", err)
			}
			fmt.Println(fileName)
			metaDataYamlFile, err := ioutil.ReadFile(filepath.Join(sourceDirectoryPath, fileName))
			if err != nil {
				utils.HandleErrorAndExit("Cannot read the meta file", err)
			}
			errUnmarshal := yaml.Unmarshal(metaDataYamlFile, &apiMetaData)
			if errUnmarshal != nil {
				utils.HandleErrorAndExit("Error unmarshal api configmap into struct ", errUnmarshal)
			}
			apiName := utils.GetValidK8sName(apiMetaData.Name)
			apiCrd.Name = apiName
			apiCrd.Spec.SwaggerConfigMapName = fmt.Sprintf("%v-cm", apiName)
			apiCrd.Spec.ParamsValues = fmt.Sprintf("%v-params", apiName)
			apiCrd.Spec.CertsValues = fmt.Sprintf("%v-certs", apiName)

			break
		} else if strings.EqualFold(fileName, utils.MetaFileApplication) { // if project artifact is a Application project
			metaDataFileFound = true
			err := utils.CopyFile(filepath.Join(sourceDirectoryPath, fileName), filepath.Join(deploymentDirPath, utils.MetaFileApplication))
			if err != nil {
				utils.HandleErrorAndExit("Cannot copy metadata file from the source directory ", err)
			}
			metaDataYamlFile, err := ioutil.ReadFile(deploymentDirPath + fileName)
			if err != nil {
				utils.HandleErrorAndExit("Cannot read the meta file", err)
			}
			errUnmarshal := yaml.Unmarshal(metaDataYamlFile, &apiMetaData)
			if errUnmarshal != nil {
				utils.HandleErrorAndExit("Error unmarshal api configmap into struct ", errUnmarshal)
			}
			apiName := utils.GetValidK8sName(apiMetaData.Name)
			apiCrd.Name = apiName
			apiCrd.Spec.SwaggerConfigMapName = fmt.Sprintf("%v-cm", apiName)
			apiCrd.Spec.ParamsValues = fmt.Sprintf("%v-params", apiName)
			apiCrd.Spec.CertsValues = fmt.Sprintf("%v-certs", apiName)
			break
		}
	}
	// if *_meta.yaml is not found inside the source directory
	if !metaDataFileFound {
		utils.HandleErrorAndExit("Cannot find metadata file inside the source directory ", err)
	}

	// write to api_crd.yaml file
	byteVal, errMarshal := yaml.Marshal(apiCrd)
	if errMarshal != nil {
		utils.HandleErrorAndExit("Error marshal API configmap ", errMarshal)
	}
	err = ioutil.WriteFile(filepath.Join(deploymentDirPath, "api_crd.yaml"), byteVal, os.ModePerm)
	if err != nil {
		utils.HandleErrorAndExit("Error creating api_crd.yaml file", err)
	}

	apiParamsData, _ := box.Get("/sample/api_params.yaml")
	apiParamsCmData, _ := box.Get("/kubernetes_resources/api_params_cm.yaml")
	apiParamsCm := &corev1.ConfigMap{}
	errUnmarshal = yaml.Unmarshal(apiParamsCmData, apiParamsCm)
	if errUnmarshal != nil {
		utils.HandleErrorAndExit("Error unmarshal api configmap into struct ", errUnmarshal)
	}
	apiParamsCm.Name = apiCrd.Spec.ParamsValues
	apiParamsCm.Data = map[string]string{"params.yaml": string(apiParamsData)}
	byteParamsVal, errParamsMarshal := yaml.Marshal(apiParamsCm)
	if errParamsMarshal != nil {
		utils.HandleErrorAndExit("Error marshal API configmap ", errMarshal)
	}
	err = ioutil.WriteFile(filepath.Join(deploymentDirPath, fmt.Sprintf("%v-params.yaml", apiMetaData.Name)),
		byteParamsVal, os.ModePerm)
	if err != nil {
		utils.HandleErrorAndExit("Error creating sample api_params.yaml file", err)
	}

	// Generate required directories inside the deployment directory
	err = createDeploymentContentDirectories(deploymentDirPath)
	if err != nil {
		return err
	}

	//remove temporary directories
	err = os.RemoveAll(tempDirPath)
	if err != nil {
		return err
	}

	fmt.Println("The deployment directory for " + genDeploymentDirSource + " file is generated at " + deploymentDirParent)

	return nil
}

// getEnvsCmd represents the envs command
var genDeploymentDirCmd = &cobra.Command{
	Use:     GenDeploymentDirCmdLiteral,
	Short:   GenDeploymentDirCmdShortDesc,
	Long:    GenDeploymentDirCmdLongDesc,
	Example: GenDeploymentDirCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + GenDeploymentDirCmdLiteral + " called")

		// check the destination directory is existed if it is provided
		if genDeploymentDirDestination != "" {
			if stat, err := os.Stat(genDeploymentDirDestination); !os.IsNotExist(err) {
				if !stat.IsDir() {
					fmt.Printf("%s is not a directory\n", genDeploymentDirDestination)
					os.Exit(1)
				}
			}
		}

		err := executeGenDeploymentDirCmd()
		if err != nil {
			utils.HandleErrorAndContinue("Error initializing the Deployment directory", err)
		}
	},
}

func init() {
	GenCmd.AddCommand(genDeploymentDirCmd)
	genDeploymentDirCmd.Flags().StringVarP(&genDeploymentDirDestination, "destination", "d", "", "Path of "+
		"the directory where the directory should be generated")
	genDeploymentDirCmd.Flags().StringVarP(&genDeploymentDirSource, "source", "s", "", "Path of "+
		"the source directory to be used when generating the directory")
	_ = genDeploymentDirCmd.MarkFlagRequired("source")
}
