package git

import "github.com/wso2/product-apim-tooling/import-export-cli/specs/params"

type VCSConfig struct {
    Repos map[string]Repo `yaml:"repos"`
}

type Environment struct {
    LastAttemptedRev  string                             `yaml:"lastAttemptedRev"`
    LastSuccessfulRev string                             `yaml:"lastSuccessfulRev"`
    FailedProjects    map[string][]*params.ProjectParams `yaml:"failedProjects"`
}

type Repo struct {
    Environments map[string]Environment `yaml:"environments"`
}

type RepoInfo struct {
    Id string `yaml:"id"`
}
