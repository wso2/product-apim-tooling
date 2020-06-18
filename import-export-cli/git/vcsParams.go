package git

type VCSConfig struct {
    Environments map[string]Environment `yaml:"environments"`
}

type Environment struct {
    LastAttemptedRev string `yaml:"lastAttemptedRev"`
    LastSuccessfulRev string `yaml:"lastSuccessfulRev"`
    FailedProjects map[string][]string `yaml:"failedProjects"`
}
