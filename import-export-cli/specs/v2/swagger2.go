package v2

import (
	"fmt"
	"path"

	"github.com/Jeffail/gabs"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/mitchellh/mapstructure"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

func swagger2XWO2BasePath(document *loads.Document) (string, bool) {
	if v, ok := document.Spec().Extensions["x-wso2-basePath"]; ok {
		str, ok := v.(string)
		return str, ok
	}
	return "", false
}

func swagger2XWSO2Cors(document *loads.Document) (*CorsConfiguration, bool, error) {
	if v, ok := document.Spec().Extensions["x-wso2-cors"]; ok {
		var cors CorsConfiguration
		err := mapstructure.Decode(v, &cors)
		if err != nil {
			return nil, true, err
		}
		cors.CorsConfigurationEnabled = true
		return &cors, true, nil
	}
	return nil, false, nil
}

func swagger2Tags(document *loads.Document) []string {
	tags := make([]string, len(document.Spec().Tags))
	for i, v := range document.Spec().Tags {
		tags[i] = v.Name
	}
	return tags
}

func swagger2XWSO2ProductionEndpoints(document *loads.Document) (*Endpoints, bool, error) {
	if v, ok := document.Spec().Extensions["x-wso2-production-endpoints"]; ok {
		var prodEp Endpoints
		err := mapstructure.Decode(v, &prodEp)
		if err != nil {
			return nil, true, err
		}
		return &prodEp, true, nil
	}
	return &Endpoints{}, false, nil
}

func swagger2XWSO2SandboxEndpoints(document *loads.Document) (*Endpoints, bool, error) {
	if v, ok := document.Spec().Extensions["x-wso2-sandbox-endpoints"]; ok {
		var sandboxEp Endpoints
		err := mapstructure.Decode(v, &sandboxEp)
		if err != nil {
			return nil, true, err
		}
		return &sandboxEp, true, nil
	}
	return &Endpoints{}, false, nil
}

// oai3GetHttpVerbs generates verbs for api definition
func swagger2GetHttpVerbs(item spec.PathItem) (verbs []string) {
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

// BuildAPIMEndpoints builds endpointConfig for given config
func BuildAPIMEndpoints(production, sandbox *Endpoints) (string, error) {
	epType := EpHttp
	if len(production.Urls) > 1 {
		epType = EpLoadbalance
		if production.Type == EpFailover {
			epType = EpFailover
		}
	}

	if len(production.Urls) == 0 {
		if len(sandbox.Urls) > 1 {
			epType = EpLoadbalance
		}
		if sandbox.Type == EpFailover {
			epType = EpFailover
		}
	}

	switch epType {
	case EpHttp:
		endpoint := buildHttpEndpoint(production, sandbox)
		return endpoint, nil
	case EpLoadbalance:
		endpoint := buildLoadBalancedEndpoints(production, sandbox)
		return endpoint, nil
	case EpFailover:
		endpoint := buildFailOver(production, sandbox)
		return endpoint, nil
	default:
		return "", fmt.Errorf("unknown endpoint type")
	}
}

func buildFailOver(production *Endpoints, sandbox *Endpoints) string {
	jsonObj, _ := gabs.ParseJSON([]byte(`
					{
						"endpoint_type": "failover",
		    			"algoCombo": "org.apache.synapse.endpoints.algorithms.RoundRobin",
		    			"algoClassName": "",
						"sessionManagement": "",
		    			"sessionTimeOut": "",
		    			"failOver": "True"
					}
				`))
	if len(production.Urls) > 0 {
		buildFailOverUrls(jsonObj, production, "production")
	}
	if len(sandbox.Urls) > 0 {
		buildFailOverUrls(jsonObj, sandbox, "sandbox")
	}
	return jsonObj.String()
}

func buildFailOverUrls(jsonObj *gabs.Container, endpoints *Endpoints, eptype string) {
	_, _ = jsonObj.Set(utils.Endpoint{Url: &endpoints.Urls[0]}, fmt.Sprintf("%s_endpoints", eptype))
	rest := endpoints.Urls[1:]
	if len(rest) > 0 {
		fo := make([]utils.Endpoint, len(rest))
		for i := 0; i < len(fo); i++ {
			fo[i] = utils.Endpoint{Url: &rest[i]}
		}
		if len(fo) > 0 {
			_, _ = jsonObj.Set(fo, fmt.Sprintf("%s_failovers", eptype))
		}
	}
}

func buildLoadBalancedEndpoints(production *Endpoints, sandbox *Endpoints) string {
	jsonObj, _ := gabs.ParseJSON([]byte(`
		{
			"endpoint_type": "load_balance",
		    "algoCombo": "org.apache.synapse.endpoints.algorithms.RoundRobin",
		    "algoClassName": "org.apache.synapse.endpoints.algorithms.RoundRobin",
		    "sessionManagement": "",
		    "sessionTimeOut": ""
		}
	`))
	prodEps := make([]utils.Endpoint, len(production.Urls))
	for i := 0; i < len(prodEps); i++ {
		prodEps[i] = utils.Endpoint{Url: &production.Urls[i]}
	}
	if len(prodEps) > 0 {
		_, _ = jsonObj.Set(prodEps, "production_endpoints")
	}

	sandboxEps := make([]utils.Endpoint, len(sandbox.Urls))
	for i := 0; i < len(sandboxEps); i++ {
		sandboxEps[i] = utils.Endpoint{Url: &sandbox.Urls[i]}
	}
	if len(sandboxEps) > 0 {
		_, _ = jsonObj.Set(sandboxEps, "sandbox_endpoints")
	}

	return jsonObj.String()
}

func buildHttpEndpoint(production *Endpoints, sandbox *Endpoints) string {
	jsonObj := gabs.New()
	_, _ = jsonObj.Set(EpHttp, "endpoint_type")
	if len(production.Urls) > 0 {
		var ep utils.Endpoint
		ep.Url = &production.Urls[0]
		_, _ = jsonObj.SetP(ep, "production_endpoints")
	}
	if len(sandbox.Urls) > 0 {
		var ep utils.Endpoint
		ep.Url = &sandbox.Urls[0]
		_, _ = jsonObj.SetP(ep, "sandbox_endpoints")
	}
	return jsonObj.String()
}

// generateFieldsFromSwagger3 using swagger
func Swagger2Populate(def *APIDefinition, document *loads.Document) error {
	def.ID.APIName = utils.ToPascalCase(document.Spec().Info.Title)
	def.ID.Version = document.Spec().Info.Version
	def.ID.ProviderName = "admin"
	def.Description = document.Spec().Info.Description
	def.Context = fmt.Sprintf("/%s/%s", def.ID.APIName, def.ID.Version)
	def.ContextTemplate = fmt.Sprintf("/%s/{version}", def.ID.APIName)
	def.Tags = swagger2Tags(document)

	// fill basepath from swagger
	if document.BasePath() != "" {
		def.Context = path.Clean(fmt.Sprintf("/%s/%s", document.BasePath(), def.ID.Version))
		def.ContextTemplate = path.Clean(fmt.Sprintf("/%s/{version}", document.BasePath()))
	}

	// override basepath if wso2 extension provided
	if basepath, ok := swagger2XWO2BasePath(document); ok {
		def.Context = path.Clean(fmt.Sprintf("/%s/%s", basepath, def.ID.Version))
		def.ContextTemplate = path.Clean(fmt.Sprintf("/%s/{version}", document.BasePath()))
	}

	cors, ok, err := swagger2XWSO2Cors(document)
	if err != nil && ok {
		return err
	}
	if ok {
		def.CorsConfiguration = cors
	}

	prodEp, foundProdEp, err := swagger2XWSO2ProductionEndpoints(document)
	if err != nil && foundProdEp {
		return err
	}
	sandboxEp, foundSandboxEp, err := swagger2XWSO2SandboxEndpoints(document)
	if err != nil && foundSandboxEp {
		return err
	}

	if foundProdEp || foundSandboxEp {
		ep, err := BuildAPIMEndpoints(prodEp, sandboxEp)
		if err != nil {
			return err
		}
		def.EndpointConfig = &ep
	}

	var uriTemplates []URITemplates
	for uri, info := range document.Spec().Paths.Paths {
		uriTemplate := URITemplates{}
		uriTemplate.URITemplate = uri
		verbs := swagger2GetHttpVerbs(info)
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
	return nil
}
