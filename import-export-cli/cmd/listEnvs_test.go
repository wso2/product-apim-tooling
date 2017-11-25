package cmd

import (
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"testing"
)

func TestPrintEnvs(t *testing.T) {
	envEndpoints := make(map[string]utils.EnvEndpoints)
	envEndpoints["dev"] = utils.EnvEndpoints{
		"publisher-enpdoint",
		"reg-endpoint",
		"token-endpoint",
	}
	printEnvs(envEndpoints)
}
