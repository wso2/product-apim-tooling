package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"
	"unicode"

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
)

type Swagger2SpecPartial struct {
	BasePath string `json:"basePath,omitempty" yaml:"basePath,omitempty"`
}

const (
	// Time format which used to output date for api
	timeFormat = "Jan _2, 2006 03:04:05 PM"
	// default tiers of API Manager
	defaultTiers = `
[
	{
      "name": "Bronze",
      "displayName": "Bronze",
      "description": "Allows 1000 requests per minute",
      "requestsPerMin": 1000,
      "requestCount": 1000,
      "unitTime": 1,
      "timeUnit": "min",
      "tierPlan": "FREE",
      "stopOnQuotaReached": true
    },
    {
      "name": "Gold",
      "displayName": "Gold",
      "description": "Allows 5000 requests per minute",
      "requestsPerMin": 5000,
      "requestCount": 5000,
      "unitTime": 1,
      "timeUnit": "min",
      "tierPlan": "FREE",
      "stopOnQuotaReached": true
    },
    {
      "name": "Silver",
      "displayName": "Silver",
      "description": "Allows 2000 requests per minute",
      "requestsPerMin": 2000,
      "requestCount": 2000,
      "unitTime": 1,
      "timeUnit": "min",
      "tierPlan": "FREE",
      "stopOnQuotaReached": true
    },
    {
      "name": "Unlimited",
      "displayName": "Unlimited",
      "description": "Allows unlimited requests",
      "requestsPerMin": 2147483647,
      "requestCount": 2147483647,
      "unitTime": 0,
      "timeUnit": "ms",
      "tierPlan": "FREE",
      "stopOnQuotaReached": true
    }
]
`
	// default cors configuration
	defaultCorsConfig = `
{
    "corsConfigurationEnabled": false,
    "accessControlAllowOrigins": [
      "*"
    ],
    "accessControlAllowCredentials": false,
    "accessControlAllowHeaders": [
      "authorization",
      "Access-Control-Allow-Origin",
      "Content-Type",
      "SOAPAction"
    ],
    "accessControlAllowMethods": [
      "GET",
      "PUT",
      "POST",
      "DELETE",
      "PATCH",
      "OPTIONS"
    ]
}
`
)

// APIDefinition represents an API artifact in APIM
type APIDefinition struct {
	ID                                 ID                `json:"id,omitempty"`
	UUID                               string            `json:"uuid,omitempty"`
	Description                        string            `json:"description,omitempty"`
	Type                               string            `json:"type,omitempty"`
	Context                            string            `json:"context"`
	ContextTemplate                    string            `json:"contextTemplate"`
	Tags                               []string          `json:"tags"`
	Documents                          []interface{}     `json:"documents"`
	LastUpdated                        string            `json:"lastUpdated,omitempty"`
	AvailableTiers                     []AvailableTiers  `json:"availableTiers,omitempty"`
	AvailableSubscriptionLevelPolicies []interface{}     `json:"availableSubscriptionLevelPolicies"`
	URITemplates                       []URITemplates    `json:"uriTemplates"`
	APIHeaderChanged                   bool              `json:"apiHeaderChanged"`
	APIResourcePatternsChanged         bool              `json:"apiResourcePatternsChanged"`
	Status                             string            `json:"status"`
	TechnicalOwner                     string            `json:"technicalOwner,omitempty"`
	TechnicalOwnerEmail                string            `json:"technicalOwnerEmail,omitempty"`
	BusinessOwner                      string            `json:"businessOwner,omitempty"`
	BusinessOwnerEmail                 string            `json:"businessOwnerEmail,omitempty"`
	Visibility                         string            `json:"visibility"`
	EndpointSecured                    bool              `json:"endpointSecured"`
	EndpointAuthDigest                 bool              `json:"endpointAuthDigest"`
	EndpointUTUsername                 string            `json:"endpointUTUsername,omitempty"`
	Transports                         string            `json:"transports"`
	InSequence                         string            `json:"inSequence,omitempty"`
	OutSequence                        string            `json:"outSequence,omitempty"`
	FaultSequence                      string            `json:"faultSequence,omitempty"`
	AdvertiseOnly                      bool              `json:"advertiseOnly"`
	CorsConfiguration                  CorsConfiguration `json:"corsConfiguration"`
	EndpointConfig                     *string           `json:"endpointConfig"`
	ResponseCache                      string            `json:"responseCache"`
	CacheTimeout                       int               `json:"cacheTimeout"`
	Implementation                     string            `json:"implementation"`
	AuthorizationHeader                string            `json:"authorizationHeader,omitempty"`
	Scopes                             []interface{}     `json:"scopes"`
	IsDefaultVersion                   bool              `json:"isDefaultVersion"`
	IsPublishedDefaultVersion          bool              `json:"isPublishedDefaultVersion"`
	Environments                       []string          `json:"environments"`
	CreatedTime                        string            `json:"createdTime,omitempty"`
	AdditionalProperties               map[string]string `json:"additionalProperties,omitempty"`
	EnvironmentList                    []string          `json:"environmentList"`
	APISecurity                        string            `json:"apiSecurity"`
	AccessControl                      string            `json:"accessControl"`
	Rating                             float64           `json:"rating"`
	IsLatest                           bool              `json:"isLatest"`
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
	Key         string `json:"key"`
	Name        string `json:"name"`
	Roles       string `json:"roles"`
	Description string `json:"description"`
	ID          int    `json:"id,omitempty"`
}
type MediationScripts struct {
}
type URITemplates struct {
	URITemplate          string           `json:"uriTemplate,omitempty"`
	HTTPVerb             string           `json:"httpVerb,omitempty"`
	AuthType             string           `json:"authType,omitempty"`
	HTTPVerbs            []string         `json:"httpVerbs,omitempty"`
	AuthTypes            []string         `json:"authTypes,omitempty"`
	ThrottlingConditions []interface{}    `json:"throttlingConditions,omitempty"`
	ThrottlingTier       string           `json:"throttlingTier,omitempty"`
	ThrottlingTiers      []string         `json:"throttlingTiers,omitempty"`
	MediationScript      string           `json:"mediationScript,omitempty"`
	Scopes               []*Scopes        `json:"scopes,omitempty"`
	MediationScripts     MediationScripts `json:"mediationScripts,omitempty"`
}
type CorsConfiguration struct {
	CorsConfigurationEnabled      bool     `json:"corsConfigurationEnabled"`
	AccessControlAllowOrigins     []string `json:"accessControlAllowOrigins"`
	AccessControlAllowCredentials bool     `json:"accessControlAllowCredentials"`
	AccessControlAllowHeaders     []string `json:"accessControlAllowHeaders"`
	AccessControlAllowMethods     []string `json:"accessControlAllowMethods"`
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

	def.LastUpdated = time.Now().Format(timeFormat)
	def.CreatedTime = strconv.FormatInt(time.Now().Unix(), 10)
	return def, nil
}

// newApiDefinitionWithDefaults creates a definition with defaults
func newApiDefinitionWithDefaults() *APIDefinition {
	def := &APIDefinition{}
	def.ID.ProviderName = "admin"
	def.CorsConfiguration = getDefaultCORS()
	def.AvailableTiers = getDefaultTiers()
	def.LastUpdated = time.Now().Format(timeFormat)
	def.CreatedTime = strconv.FormatInt(time.Now().Unix(), 10)
	def.Status = "CREATED"
	def.Environments = []string{"Production and Sandbox"}
	def.EnvironmentList = []string{"SANDBOX", "PRODUCTION"}
	def.CacheTimeout = 300
	def.IsPublishedDefaultVersion = false
	def.ResponseCache = "Disabled"
	def.EndpointConfig = nil
	def.APISecurity = "oauth2"
	def.Rating = 0.0
	def.AccessControl = "all"
	def.Visibility = "public"
	def.Type = "HTTP"
	def.Implementation = "ENDPOINT"
	def.EndpointSecured = false
	def.EndpointAuthDigest = false
	def.AdvertiseOnly = false
	def.APIHeaderChanged = false
	def.APIResourcePatternsChanged = false
	def.IsLatest = false
	def.IsDefaultVersion = false
	def.IsPublishedDefaultVersion = false
	def.Transports = "http,https"
	def.Tags = []string{}
	def.Documents = []interface{}{}
	def.AvailableSubscriptionLevelPolicies = []interface{}{}

	return def
}

// getDefaultTiers populates default tiers
func getDefaultTiers() []AvailableTiers {
	var tiers []AvailableTiers
	err := json.Unmarshal([]byte(defaultTiers), &tiers)
	if err != nil {
		panic(err)
	}
	return tiers
}

// getDefaultCORS populates cors config
func getDefaultCORS() CorsConfiguration {
	var cors CorsConfiguration
	err := json.Unmarshal([]byte(defaultCorsConfig), &cors)
	if err != nil {
		panic(err)
	}
	return cors
}

// loads swagger from swaggerDoc
// swagger2.0/openapi3.0 specs are supported
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

// executeInitCmd will run init command
func executeInitCmd() error {
	var dir string
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
	fmt.Println("Initializing a new APIM project in", dir)

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
		// set context based on basepath if presented(only for swagger 2.0)
		if basePath, ok := holder["basePath"]; ok {
			def.Context = path.Clean(fmt.Sprintf("/%s/%s", basePath, sw.Info.Version))
			def.ContextTemplate = path.Clean(fmt.Sprintf("/%s/{version}", basePath))
		}

		// add indention with two spaces
		swaggerSavePath := filepath.Join(initCmdOutputDir, filepath.FromSlash("Meta-information/swagger.json"))
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
	fmt.Println("Project initialized")
	return err
}

var InitCommand = &cobra.Command{
	Use:   "init",
	Short: "initialize a new project in current directory",
	Long:  "initialize a new project in current directory. If a swagger file provided API will be populated with details from swagger",
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + "init called")
		err := executeInitCmd()
		if err != nil {
			utils.HandleErrorAndExit("Error initializing project", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(InitCommand)
	InitCommand.Flags().StringVarP(&initCmdOutputDir, "output", "o", "", "Output directory for API. When not "+
		"provided project will be initialized in current working directory")
	InitCommand.Flags().StringVarP(&initCmdApiDefinitionPath, "definition", "d", "", "Provide a "+
		"YAML definition of API")
	InitCommand.Flags().StringVarP(&initCmdSwaggerPath, "swagger", "s", "", "Provide a swagger"+
		"file for the API (json/yaml)")
	InitCommand.Flags().BoolVarP(&initCmdEnvInject, "env-inject", "", false, "Inject "+
		"environment variables to definition file")
}
