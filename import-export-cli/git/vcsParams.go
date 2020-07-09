package git

import "github.com/wso2/product-apim-tooling/import-export-cli/specs/params"

type VCSConfig struct {
    Environments map[string]Environment `yaml:"environments"`
}

type Environment struct {
    LastAttemptedRev string `yaml:"lastAttemptedRev"`
    LastSuccessfulRev string `yaml:"lastSuccessfulRev"`
    FailedProjects map[string][]*params.ProjectParams `yaml:"failedProjects"`
}
