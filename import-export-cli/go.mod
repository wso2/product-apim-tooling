go 1.14

require (
	github.com/Jeffail/gabs v1.4.0
	github.com/getkin/kin-openapi v0.2.0
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/go-openapi/loads v0.19.5
	github.com/go-resty/resty/v2 v2.4.0
	github.com/google/go-cmp v0.4.0
	github.com/google/uuid v1.1.1
	github.com/hashicorp/go-multierror v1.0.0
	github.com/json-iterator/go v1.1.10
	github.com/magiconair/properties v1.8.1
	github.com/mitchellh/mapstructure v1.3.2
	github.com/pavel-v-chernykh/keystore-go/v4 v4.1.0
	github.com/renstrom/dedent v1.0.0
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.1.1
	github.com/stretchr/testify v1.6.1
	github.com/wso2/k8s-api-operator/api-operator v0.0.0-20210223103109-66ee766c8413
	golang.org/x/crypto v0.0.0-20200414173820-0848c9571904
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/api v0.18.2 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)

replace k8s.io/client-go => k8s.io/client-go v0.18.2

module github.com/wso2/product-apim-tooling/import-export-cli
