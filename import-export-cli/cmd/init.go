package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"text/template"
	"unicode"

	"github.com/wso2/product-apim-tooling/import-export-cli/defaults"

	"github.com/ghodss/yaml"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/cobra"

	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

var (
	initCmdOutputDir         string
	initCmdSwaggerPath       string
	initCmdApiDefinitionPath string
	initCmdEnvInject         bool
	initCmdForced            bool
)

type Swagger2SpecPartial struct {
	BasePath string `json:"basePath,omitempty" yaml:"basePath,omitempty"`
}

// APIDefinition represents an API artifact in APIM
type APIDefinition struct {
	ID                                 ID                 `json:"id,omitempty"`
	UUID                               string             `json:"uuid,omitempty"`
	Description                        string             `json:"description,omitempty"`
	Type                               string             `json:"type,omitempty"`
	Context                            string             `json:"context"`
	ContextTemplate                    string             `json:"contextTemplate"`
	Tags                               []string           `json:"tags"`
	Documents                          []interface{}      `json:"documents,omitempty"`
	LastUpdated                        string             `json:"lastUpdated,omitempty"`
	AvailableTiers                     []AvailableTiers   `json:"availableTiers,omitempty"`
	AvailableSubscriptionLevelPolicies []interface{}      `json:"availableSubscriptionLevelPolicies,omitempty"`
	URITemplates                       []URITemplates     `json:"uriTemplates"`
	APIHeaderChanged                   bool               `json:"apiHeaderChanged,omitempty"`
	APIResourcePatternsChanged         bool               `json:"apiResourcePatternsChanged,omitempty"`
	Status                             string             `json:"status,omitempty"`
	TechnicalOwner                     string             `json:"technicalOwner,omitempty"`
	TechnicalOwnerEmail                string             `json:"technicalOwnerEmail,omitempty"`
	BusinessOwner                      string             `json:"businessOwner,omitempty"`
	BusinessOwnerEmail                 string             `json:"businessOwnerEmail,omitempty"`
	Visibility                         string             `json:"visibility,omitempty"`
	EndpointSecured                    bool               `json:"endpointSecured,omitempty"`
	EndpointAuthDigest                 bool               `json:"endpointAuthDigest,omitempty"`
	EndpointUTUsername                 string             `json:"endpointUTUsername,omitempty"`
	Transports                         string             `json:"transports,omitempty"`
	InSequence                         string             `json:"inSequence,omitempty"`
	OutSequence                        string             `json:"outSequence,omitempty"`
	FaultSequence                      string             `json:"faultSequence,omitempty"`
	AdvertiseOnly                      bool               `json:"advertiseOnly,omitempty"`
	CorsConfiguration                  *CorsConfiguration `json:"corsConfiguration,omitempty"`
	EndpointConfig                     *string            `json:"endpointConfig,omitempty"`
	ResponseCache                      string             `json:"responseCache,omitempty"`
	CacheTimeout                       int                `json:"cacheTimeout,omitempty"`
	Implementation                     string             `json:"implementation,omitempty"`
	AuthorizationHeader                string             `json:"authorizationHeader,omitempty"`
	Scopes                             []interface{}      `json:"scopes,omitempty"`
	IsDefaultVersion                   bool               `json:"isDefaultVersion,omitempty"`
	IsPublishedDefaultVersion          bool               `json:"isPublishedDefaultVersion,omitempty"`
	Environments                       []string           `json:"environments,omitempty"`
	CreatedTime                        string             `json:"createdTime,omitempty"`
	AdditionalProperties               map[string]string  `json:"additionalProperties,omitempty"`
	EnvironmentList                    []string           `json:"environmentList,omitempty"`
	APISecurity                        string             `json:"apiSecurity,omitempty"`
	AccessControl                      string             `json:"accessControl,omitempty"`
	Rating                             float64            `json:"rating,omitempty"`
	IsLatest                           bool               `json:"isLatest,omitempty"`
}
type ID struct {
	ProviderName string `json:"providerName"`
	APIName      string `json:"apiName"`
	Version      string `json:"version"`
}
type AvailableTiers struct {
	Name               string `json:"name,omitempty"`
	DisplayName        string `json:"displayName,omitempty"`
	Description        string `json:"description,omitempty"`
	RequestsPerMin     int    `json:"requestsPerMin,omitempty"`
	RequestCount       int    `json:"requestCount,omitempty"`
	UnitTime           int    `json:"unitTime,omitempty"`
	TimeUnit           string `json:"timeUnit,omitempty"`
	TierPlan           string `json:"tierPlan,omitempty"`
	StopOnQuotaReached bool   `json:"stopOnQuotaReached,omitempty"`
}
type Scopes struct {
	Key         string `json:"key,omitempty"`
	Name        string `json:"name,omitempty"`
	Roles       string `json:"roles,omitempty"`
	Description string `json:"description,omitempty"`
	ID          int    `json:"id,omitempty"`
}
type MediationScripts struct {
}
type URITemplates struct {
	URITemplate          string            `json:"uriTemplate,omitempty"`
	HTTPVerb             string            `json:"httpVerb,omitempty"`
	AuthType             string            `json:"authType,omitempty"`
	HTTPVerbs            []string          `json:"httpVerbs,omitempty"`
	AuthTypes            []string          `json:"authTypes,omitempty"`
	ThrottlingConditions []interface{}     `json:"throttlingConditions,omitempty"`
	ThrottlingTier       string            `json:"throttlingTier,omitempty"`
	ThrottlingTiers      []string          `json:"throttlingTiers,omitempty"`
	MediationScript      string            `json:"mediationScript,omitempty"`
	Scopes               []*Scopes         `json:"scopes,omitempty"`
	MediationScripts     *MediationScripts `json:"mediationScripts,omitempty"`
}
type CorsConfiguration struct {
	CorsConfigurationEnabled      bool     `json:"corsConfigurationEnabled,omitempty"`
	AccessControlAllowOrigins     []string `json:"accessControlAllowOrigins,omitempty"`
	AccessControlAllowCredentials bool     `json:"accessControlAllowCredentials,omitempty"`
	AccessControlAllowHeaders     []string `json:"accessControlAllowHeaders,omitempty"`
	AccessControlAllowMethods     []string `json:"accessControlAllowMethods,omitempty"`
}

// directories to be created
var dirs = []string{
	"Meta-information",
	"Image",
	"Docs",
	"Docs/FileContents",
	"Sequences",
	"Sequences/fault-sequence",
	"Sequences/in-sequence",
	"Sequences/out-sequence",
}

// createDirectories will create dirs in current working directory
func createDirectories(name string) error {
	for _, dir := range dirs {
		dirPath := filepath.Join(name, filepath.FromSlash(dir))
		utils.Logln(utils.LogPrefixInfo + "Creating directory " + dirPath)
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

// loadDefaultSpecFromDisk loads api definition stored in HOME/.wso2apimcli/default_api.yaml
func loadDefaultSpecFromDisk() (*APIDefinition, error) {
	defaultData, err := ioutil.ReadFile(utils.DefaultAPISpecFilePath)
	if err != nil {
		return nil, err
	}
	def := &APIDefinition{}
	err = yaml.Unmarshal(defaultData, &def)
	if err != nil {
		return nil, err
	}
	return def, nil
}

// loads swagger from swaggerDoc
// swagger2.0/OpenAPI3.0 specs are supported
func loadSwagger(swaggerDoc string) (*openapi3.Swagger, []byte, error) {
	utils.Logln(utils.LogPrefixInfo + "Loading swagger from " + swaggerDoc)
	buffer, err := ioutil.ReadFile(swaggerDoc)
	if err != nil {
		return nil, nil, err
	}

	sw, err := openapi3.NewSwaggerLoader().LoadSwaggerFromData(buffer)
	if err != nil {
		return nil, nil, err
	}
	return sw, buffer, nil
}

// generateFieldsFromSwagger3 using swagger
func (def *APIDefinition) generateFieldsFromSwagger3(swagger *openapi3.Swagger) {
	def.ID.APIName = utils.ToPascalCase(swagger.Info.Title)
	def.ID.Version = swagger.Info.Version
	def.Description = swagger.Info.Description
	def.Context = fmt.Sprintf("/%s/%s", def.ID.APIName, def.ID.Version)
	def.ContextTemplate = fmt.Sprintf("/%s/{version}", def.ID.APIName)

	var uriTemplates []URITemplates
	for uri, info := range swagger.Paths {
		uriTemplate := URITemplates{}
		uriTemplate.URITemplate = uri
		verbs := getHttpVerbs(info)
		uriTemplate.HTTPVerbs = verbs
		if len(verbs) > 0 {
			uriTemplate.HTTPVerb = verbs[0]
		}
		authTypes := make([]string, len(verbs))
		throttlingTiers := make([]string, len(verbs))
		for i := 0; i < len(verbs); i++ {
			authTypes[i] = "Any"
			throttlingTiers[i] = "Unlimited"
		}
		uriTemplate.AuthType = "Any"
		uriTemplate.AuthTypes = authTypes
		uriTemplate.ThrottlingTier = "Unlimited"
		uriTemplate.ThrottlingTiers = throttlingTiers
		uriTemplate.Scopes = make([]*Scopes, len(verbs))
		uriTemplates = append(uriTemplates, uriTemplate)
	}
	def.URITemplates = uriTemplates
}

// getHttpVerbs generates verbs for api definition
func getHttpVerbs(item *openapi3.PathItem) (verbs []string) {
	if item.Get != nil {
		verbs = append(verbs, "GET")
	}
	if item.Post != nil {
		verbs = append(verbs, "POST")
	}
	if item.Put != nil {
		verbs = append(verbs, "PUT")
	}
	if item.Delete != nil {
		verbs = append(verbs, "DELETE")
	}
	if item.Patch != nil {
		verbs = append(verbs, "PATCH")
	}
	if item.Head != nil {
		verbs = append(verbs, "HEAD")
	}
	if item.Options != nil {
		verbs = append(verbs, "OPTIONS")
	}
	return
}

// hasJSONPrefix returns true if the provided buffer appears to start with
// a JSON open brace.
func hasJSONPrefix(buf []byte) bool {
	return hasPrefix(buf, []byte("{"))
}

// Return true if the first non-whitespace bytes in buf is prefix.
func hasPrefix(buf []byte, prefix []byte) bool {
	trim := bytes.TrimLeftFunc(buf, unicode.IsSpace)
	return bytes.HasPrefix(trim, prefix)
}

func generateConfig(file string) error {
	envs := utils.GetMainConfigFromFile(utils.MainConfigFilePath)
	t, err := template.New("").Parse(defaults.ApiVarsTmpl)
	if err != nil {
		return err
	}

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	err = t.Execute(f, envs.Environments)
	if err != nil {
		return err
	}
	return nil
}

// executeInitCmd will run init command
func executeInitCmd() error {
	var dir string
	swaggerSavePath := filepath.Join(initCmdOutputDir, filepath.FromSlash("Meta-information/swagger.json"))

	if initCmdOutputDir != "" {
		err := os.MkdirAll(initCmdOutputDir, os.ModePerm)
		if err != nil {
			return err
		}
		p, err := filepath.Abs(initCmdOutputDir)
		if err != nil {
			return err
		}
		dir = p
	} else {
		pwd, err := os.Getwd()
		if err != nil {
			return err
		}
		dir = pwd
	}
	fmt.Println("Initializing a new WSO2 API Manager project in", dir)

	def, err := loadDefaultSpecFromDisk()
	if err != nil {
		return err
	}

	err = createDirectories(initCmdOutputDir)
	if err != nil {
		return err
	}

	// use swagger to auto generate
	if initCmdSwaggerPath != "" {
		// load swagger from path
		sw, buff, err := loadSwagger(initCmdSwaggerPath)
		if err != nil {
			return err
		}
		def.generateFieldsFromSwagger3(sw)

		// put swagger file to corresponding directory
		// if swagger is either from json or yaml source it will be properly indented with two spaces
		// before saving into directory
		var holder map[string]interface{}
		if hasJSONPrefix(buff) {
			// try to unmarshal json
			err = json.Unmarshal(buff, &holder)
			if err != nil {
				return err
			}
		} else {
			// try to unmarshal yaml
			err = yaml.Unmarshal(buff, &holder)
			if err != nil {
				return err
			}
		}
		// set context based on basePath if presented(only for swagger 2.0)
		if basePath, ok := holder["basePath"]; ok {
			def.Context = path.Clean(fmt.Sprintf("/%s/%s", basePath, sw.Info.Version))
			def.ContextTemplate = path.Clean(fmt.Sprintf("/%s/{version}", basePath))
		}

		// add indention with two spaces
		utils.Logln(utils.LogPrefixInfo + "Writing " + swaggerSavePath)
		data, err := json.MarshalIndent(holder, "", "  ")
		if err != nil {
			return err
		}
		// write to file
		err = ioutil.WriteFile(swaggerSavePath, data, os.ModePerm)
		if err != nil {
			return err
		}
	} else {
		// create an empty swagger
		utils.Logln(utils.LogPrefixInfo + "Writing " + swaggerSavePath)
		err = ioutil.WriteFile(swaggerSavePath, defaults.Swagger, os.ModePerm)
		if err != nil {
			return err
		}
	}

	// use api definition if given
	if initCmdApiDefinitionPath != "" {
		// read definition file
		utils.Logln(utils.LogPrefixInfo + "Reading API Definition from " + initCmdApiDefinitionPath)
		content, err := ioutil.ReadFile(initCmdApiDefinitionPath)
		if err != nil {
			return err
		}

		apiDef := &APIDefinition{}
		// inject from env if requested
		if initCmdEnvInject {
			utils.Logln(utils.LogPrefixInfo + "Injecting variables to definition from environment")
			data, err := utils.InjectEnv(string(content))
			if err != nil {
				return err
			}
			content = []byte(data)
		}
		// read from yaml definition
		err = yaml.Unmarshal(content, &apiDef)
		if err != nil {
			return err
		}

		// marshal original def
		originalDefBytes, err := json.Marshal(def)
		if err != nil {
			return err
		}
		// marshal new def
		newDefBytes, err := json.Marshal(apiDef)
		if err != nil {
			return err
		}

		// merge two definitions
		finalDefBytes, err := utils.MergeJSON(originalDefBytes, newDefBytes)
		if err != nil {
			return err
		}
		tmpDef := &APIDefinition{}
		err = json.Unmarshal(finalDefBytes, &tmpDef)
		if err != nil {
			return err
		}
		def = tmpDef
	}
	// indent json with two spaces
	indentedDefBytes, err := json.MarshalIndent(def, "", "  ")
	if err != nil {
		return err
	}
	// write to the disk
	apiJSONPath := filepath.Join(initCmdOutputDir, filepath.FromSlash("Meta-information/api.json"))
	utils.Logln(utils.LogPrefixInfo + "Writing " + apiJSONPath)
	err = ioutil.WriteFile(apiJSONPath, indentedDefBytes, os.ModePerm)
	if err != nil {
		return err
	}

	apimProjConfigFilePath := filepath.Join(initCmdOutputDir, DefaultAPIMParamsFileName)
	utils.Logln(utils.LogPrefixInfo + "Writing " + apimProjConfigFilePath)
	err = generateConfig(apimProjConfigFilePath)
	if err != nil {
		return err
	}

	apimProjReadmeFilePath := filepath.Join(initCmdOutputDir, "README.txt")
	utils.Logln(utils.LogPrefixInfo + "Writing " + apimProjReadmeFilePath)
	err = ioutil.WriteFile(apimProjReadmeFilePath, defaults.ProjectReadme, os.ModePerm)
	if err != nil {
		return err
	}
	fmt.Println("Project initialized")
	return nil
}

var InitCommand = &cobra.Command{
	Use:     "init [project path]",
	Short:   "initialize a new project in given path",
	Long:    "initialize a new project in given path. If a openAPI definition provided API will be populated with details from it",
	Example: "apimcli init myapi --openapi petstore.yaml",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + "init called")
		initCmdOutputDir = args[0]

		// check for dir existence, if so stop it unless forced flag is present
		if stat, err := os.Stat(initCmdOutputDir); !os.IsNotExist(err) {
			fmt.Printf("%s already exists\n", initCmdOutputDir)
			if !stat.IsDir() {
				fmt.Printf("%s is not a directory\n", initCmdOutputDir)
				os.Exit(1)
			}
			if !initCmdForced {
				fmt.Println("Run with -f or --force to overwrite directory and create project")
				os.Exit(1)
			}
			fmt.Println("Running command in forced mode")
		}

		err := executeInitCmd()
		if err != nil {
			utils.HandleErrorAndExit("Error initializing project", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(InitCommand)
	InitCommand.Flags().StringVarP(&initCmdApiDefinitionPath, "definition", "d", "", "Provide a "+
		"YAML definition of API")
	InitCommand.Flags().StringVarP(&initCmdSwaggerPath, "openapi", "", "", "Provide an OpenAPI "+
		"definition for the API (json/yaml)")
	InitCommand.Flags().BoolVarP(&initCmdEnvInject, "env-inject", "", false, "Inject "+
		"environment variables to definition file")
	InitCommand.Flags().BoolVarP(&initCmdForced, "force", "f", false, "Force create project")
}
