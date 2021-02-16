/*
*  Copyright (c) WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
*
*  WSO2 Inc. licenses this file to you under the Apache License,
*  Version 2.0 (the "License"); you may not use this file except
*  in compliance with the License.
*  You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing,
* software distributed under the License is distributed on an
* "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
* KIND, either express or implied.  See the License for the
* specific language governing permissions and limitations
* under the License.
 */

package git

import (
    "bytes"
    "errors"
    "fmt"
    "github.com/google/uuid"
    "github.com/wso2/product-apim-tooling/import-export-cli/impl"
    "github.com/wso2/product-apim-tooling/import-export-cli/specs/params"
    "github.com/wso2/product-apim-tooling/import-export-cli/utils"
    "gopkg.in/yaml.v2"
    "io"
    "io/ioutil"
    "os"
    "os/exec"
    "path/filepath"
    "strconv"
    "strings"
)

// Reads the vcs configuration file and returns. Silently catch the error when config file is not found
// filePath is the path to look for the VCS configuration file
// returns *VCSConfig VCS configuration
func getVCSConfigFromFileSilently(filePath string) *VCSConfig {
    var vcsConfig VCSConfig
    data, err := ioutil.ReadFile(filePath)
    if err == nil {
        if err := yaml.Unmarshal(data, &vcsConfig); err != nil {
            utils.HandleErrorAndExit("VCSConfig: Error parsing "+filePath, err)
        }
    }
    return &vcsConfig
}

// Reads and returns the environment specific information from the VCS config along with the full VCS config
// repoId is the id of the git repository (located in vcs.yaml)
// environment is the name of the environment
// Returns VCSConfig, the full VCS configuration
// Returns Environment, the environment specific VCS configuration
// Returns bool, whether the environment is available in the VCS configuration or not
func getVCSEnvironmentDetails(repoId, environment string) (VCSConfig, Environment, bool)  {
    mainConfig := utils.GetMainConfigFromFile(utils.MainConfigFilePath)
    if mainConfig.Config.VCSConfigFilePath != "" {
        VCSConfigFilePath = mainConfig.Config.VCSConfigFilePath
    }
    vcsConfig := getVCSConfigFromFileSilently(VCSConfigFilePath)
    if vcsConfig.Repos == nil {
        vcsConfig.Repos = make(map[string]Repo)
    }
    envVCSConfig, hasEnv := vcsConfig.Repos[repoId].Environments[environment]
    return *vcsConfig, envVCSConfig, hasEnv
}

// Returns the status of the projects indicating the projects to deploy (need to save, delete or failed previously).
// Environment is the environment name
// fromRevType is the type of the revision the status should be taken by comparing with the current revision. The allowed values are "last_attempted", "last_successful"
// Returns string, id of the git repository (located in vcs.yaml)
// Returns int, the total number of projects to deploy
// Returns map[string][]*params.ProjectParams, the details of the projects that needs to deploy
func GetStatus(environment, fromRevType string) (string, int, map[string][]*params.ProjectParams) {
    var envRevision string
    mainConfig := utils.GetMainConfigFromFile(utils.MainConfigFilePath)
    repoId, err := getRepoId()
    if err != nil {
        utils.HandleErrorAndExit("Error while retrieving repository id", err)
    }
    if repoId == "" {
        utils.HandleErrorAndExit("The repository info: vcs.yaml is not found in the repository root. "+
            "If this is the first time you are using this repo, please initialize it with 'vcs init'.", nil)
    }
    _, envVCSConfig, hasEnv := getVCSEnvironmentDetails(repoId, environment)
    if hasEnv {
        if fromRevType == FromRevTypeLastAttempted {
            envRevision = envVCSConfig.LastAttemptedRev
        } else if fromRevType == FromRevTypeLastSuccessful {
            if len(envVCSConfig.LastSuccessfulRev) > 0 {
                envRevision = envVCSConfig.LastSuccessfulRev[0]
            } else {
                envRevision = ""
            }
        }
    }

    basePath, err := getRepoBaseDir()
    if err != nil {
        utils.HandleErrorAndExit("Error while getting repository base folder location", err)
    }

    var changedFiles string
    if envRevision == "" {
        changedFiles, _ = executeGitCommand("ls-tree", "-r", "HEAD", "--name-only", "--full-tree")
    } else if mainConfig.Config.VCSDeletionEnabled {
        changedFiles, _ = executeGitCommand("diff", "--name-only", envRevision)
    } else {
        changedFiles, _ = executeGitCommand("diff", "--diff-filter=d", "--name-only", envRevision)
    }

    //remove slashes (/) unix-format path separators with OS specific path separator
    changedFiles = strings.ReplaceAll(changedFiles,"/", string(filepath.Separator))

    changedFileList := strings.Split(changedFiles,"\n")
    // remove the last empty element
    if len(changedFileList) > 0 {
        changedFileList = changedFileList[:len(changedFileList)-1]
    }

    if utils.VerboseModeEnabled() {
        logChangedFiles(changedFileList)
    }

    changedPathInfoMap := make(map[string]*params.ProjectParams)
    updatedProjectsPerType := make(map[string][]*params.ProjectParams)
    updatedProjectsPerProjectPath := make(map[string]*params.ProjectParams)

    var totalProjectsToUpdate = 0
    for _, changedFile := range changedFileList {
        projectParam := getProjectInfoFromProjectFile(envVCSConfig, basePath, changedFile, changedPathInfoMap)
        if projectParam.Type != utils.ProjectTypeNone {
            if updatedProjectsPerType[projectParam.Type] == nil {
                updatedProjectsPerType[projectParam.Type] = []*params.ProjectParams{}
            }
            if updatedProjectsPerProjectPath[projectParam.AbsolutePath] == nil {
                updatedProjectsPerProjectPath[projectParam.AbsolutePath] = projectParam
                updatedProjectsPerType[projectParam.Type] = append(updatedProjectsPerType[projectParam.Type], projectParam)
                totalProjectsToUpdate++
            }
        }
    }

    //append failed projects to the updated project list if exists
    for _, failedProjectsInEachType := range envVCSConfig.FailedProjects {
        for _, failedProjectInEachType := range failedProjectsInEachType {
            if updatedProjectsPerProjectPath[failedProjectInEachType.AbsolutePath] == nil {
                updatedProjectsPerProjectPath[failedProjectInEachType.AbsolutePath] = failedProjectInEachType
                updatedProjectsPerType[failedProjectInEachType.Type] =
                    append(updatedProjectsPerType[failedProjectInEachType.Type], failedProjectInEachType)
                failedProjectInEachType.FailedDuringPreviousDeploy = true
                totalProjectsToUpdate++
            }
        }
    }

    return repoId, totalProjectsToUpdate, updatedProjectsPerType
}

// Returns whether the given project was failed to deploy previously
// environment is the environment name
// Returns bool indicating whether the given project was failed previously
func failedDuringEarlierDeploy(vcsEnvConfig Environment, projectParams *params.ProjectParams) bool {
    failedProjectsForType := vcsEnvConfig.FailedProjects[projectParams.Type]
    for _, failedProject := range failedProjectsForType {
        if failedProject.RelativePath == projectParams.RelativePath {
            return true
        }
    }
    return false
}

// Rollbacks the projects to the initial state when any of the projects were failed during deployment
// accesstoken is the access token to access the APIM product REST APIs
// environment is the environment name
func Rollback(accessToken, environment string) error {
    repoId, totalProjectsToUpdate, updatedProjectsPerType := GetStatus(environment, FromRevTypeLastSuccessful)
    _, envVCSConfig, hasEnv := getVCSEnvironmentDetails(repoId, environment)

    if !hasEnv || len(envVCSConfig.FailedProjects) == 0 {
        return errors.New("Nothing to rollback")
    }

    if len(envVCSConfig.LastSuccessfulRev) == 0 {
        return errors.New("Failed to rollback as there are no previous successful revisions")
    }
    lastSuccessfulRevision := envVCSConfig.LastSuccessfulRev[0]

    currentBranch := getCurrentBranch()
    tmpBranchName := "tmp-" + lastSuccessfulRevision[0:8]
    checkoutNewBranchFromRevision(tmpBranchName, lastSuccessfulRevision)
    deployUpdatedProjects(accessToken, repoId, environment, totalProjectsToUpdate, updatedProjectsPerType)
    checkoutBranch(currentBranch)
    deleteTmpBranch(tmpBranchName)
    return nil
}

// Creates a new branch from the given revision
// tmpBranchName is the new branch that is checkout from the revision
// revision is the git commit id
func checkoutNewBranchFromRevision(tmpBranchName, revision string) {
    _, err := executeGitCommand("checkout", "-b", tmpBranchName, revision)
    if err != nil {
        utils.HandleErrorAndExit("Error while checking out last successful commit ("+
            revision+") for rolling back", err)
    }
}

// Switches to the given branch name
// branchName is the name of the branch that should switch into
func checkoutBranch(branchName string) {
    _, err := executeGitCommand("checkout", branchName)
    if err != nil {
        utils.HandleErrorAndExit("Error while checking out branch " + branchName, err)
    }
}

// Returns the name of the current branch of the repository where the user is executing apictl commands
func getCurrentBranch() string {
    branch, err := executeGitCommand("rev-parse", "--abbrev-ref", "HEAD")
    if err != nil {
        utils.HandleErrorAndExit("Error while getting current branch", err)
    }
    return strings.TrimSpace(branch)
}

// Deletes the given branch. The branch name must start with "tmp-"
func deleteTmpBranch(tmpBranch string) {
    //done as a security check
    if !strings.HasPrefix(tmpBranch, "tmp-") {
        utils.HandleErrorAndExit("Cannot remove branches not starting with 'tmp-'", nil)
    }
    _, err := executeGitCommand("branch", "-D", tmpBranch)
    if err != nil {
        utils.HandleErrorAndExit("Error while deleting the temp branch "+tmpBranch, err)
    }
}

// Deletes the projects from the environment that are identified as deleted.
// accesstoken is the access token to access the APIM product REST APIs
// environment is the environment name
// deletedProjectsPerType A map that has keys as Apps/APIs or API Products and values as deleted projects of each type
// This will return the failed projects with the same structure at the end if such projects exist during deletion.
func deployProjectDeletions(accessToken, environment string, deletedProjectsPerType map[string][]*params.ProjectParams,
    failedProjects map[string][]*params.ProjectParams) map[string][]*params.ProjectParams {
    // Deleting Application projects
    applicationProjectsToDelete := deletedProjectsPerType[utils.ProjectTypeApplication]
    if len(applicationProjectsToDelete) != 0 {
        fmt.Println("\nApplications (" + strconv.Itoa(len(applicationProjectsToDelete)) + ") ...")
        for i, projectParam := range applicationProjectsToDelete {
            fmt.Println(strconv.Itoa(i+1) + ": " + projectParam.NickName + ": (" + projectParam.RelativePath + ")")
            appInfo, _, err := impl.GetApplicationDefinition(projectParam.AbsolutePath)
            if handleIfError(err, failedProjects, projectParam) {
                continue
            }
            projectParam.ProjectInfo.Name = appInfo.Name
            projectParam.ProjectInfo.Owner = appInfo.Subscriber.Name
            resp, err := impl.DeleteApplication(accessToken, environment, appInfo.Name, appInfo.Subscriber.Name)
            if handleIfError(err, failedProjects, projectParam) {
                continue
            }
            impl.PrintDeleteAppResponse(resp, err)
        }
    }

    // Deleting API Product projects
    apiProductProjectsToDelete := deletedProjectsPerType[utils.ProjectTypeApiProduct]
    if len(apiProductProjectsToDelete) != 0 {
        fmt.Println("\nAPI Products (" + strconv.Itoa(len(apiProductProjectsToDelete)) + ") ...")
        for i, projectParam := range apiProductProjectsToDelete {
            fmt.Println(strconv.Itoa(i+1) + ": " + projectParam.NickName + ": (" + projectParam.RelativePath + ")")
            apiProductInfo, _, err := impl.GetAPIProductDefinition(projectParam.AbsolutePath)
            if handleIfError(err, failedProjects, projectParam) {
                continue
            }
            projectParam.ProjectInfo.Name = apiProductInfo.ID.APIProductName
            projectParam.ProjectInfo.Owner = apiProductInfo.ID.ProviderName
            projectParam.ProjectInfo.Version = apiProductInfo.ID.Version
            resp, err := impl.DeleteAPIProduct(accessToken, environment, apiProductInfo.ID.APIProductName, apiProductInfo.ID.ProviderName)
            if handleIfError(err, failedProjects, projectParam) {
                continue
            }
            impl.PrintDeleteAPIProductResponse(resp, err)
        }
    }

    // Deleting API projects
    apiProjectsToDelete := deletedProjectsPerType[utils.ProjectTypeApi]
    if len(apiProjectsToDelete) != 0 {
        fmt.Println("\nAPIs (" + strconv.Itoa(len(apiProjectsToDelete)) + ") ...")
        for i, projectParam := range apiProjectsToDelete {
            fmt.Println(strconv.Itoa(i+1) + ": " + projectParam.NickName + ": (" + projectParam.RelativePath + ")")
            apiInfo, _, err := impl.GetAPIDefinition(projectParam.AbsolutePath)
            if handleIfError(err, failedProjects, projectParam) {
                continue
            }
            projectParam.ProjectInfo.Name = apiInfo.ID.APIName
            projectParam.ProjectInfo.Owner = apiInfo.ID.ProviderName
            projectParam.ProjectInfo.Version = apiInfo.ID.Version
            resp, err := impl.DeleteAPI(accessToken, environment, apiInfo.ID.APIName, apiInfo.ID.Version, apiInfo.ID.ProviderName)
            if handleIfError(err, failedProjects, projectParam) {
                continue
            }
            impl.PrintDeleteAPIResponse(resp, err)
        }
    }

    return failedProjects
}

// Logs the error and appends the failed project given from projectParam into the failedProjects map.
func handleIfError(err error, failedProjects map[string][]*params.ProjectParams, projectParam *params.ProjectParams) bool {
    if err != nil {
        fmt.Println("Error... ", err)
        failedProjects[projectParam.Type] = append(failedProjects[projectParam.Type], projectParam)
    }
    return err != nil
}

// Deploys the updated projects. It will only handle new or updated projects and deleted projects will be tracked and
// skipped. Those deleted projects will be returned from the 2nd return argument.
// accesstoken is the access token to access the APIM product REST APIs
// repoId is the id of the git repository (located in vcs.yaml)
// environment is the environment name
// totalProjectsToUpdate is the number of total projects that needs to be deployed.
// updatedProjectsPerType is a map of string -> ProjectParams which consists of updated projects per each type (API, App..)
// Returns bool, true if any deleted projects exists so the process should continue with project deletion path
// Returns map[string][]*params.ProjectParams, a map of project type (API, App.. ) to each project detail which are
//  deleted projects
// Returns map[string][]*params.ProjectParams, a map of project type (API, App.. ) to each project detail which are
//  failed during the deployment
func deployUpdatedProjects(accessToken, repoId, environment string, totalProjectsToUpdate int,
        updatedProjectsPerType map[string][]*params.ProjectParams) (bool, map[string][]*params.ProjectParams,
        map[string][]*params.ProjectParams) {
    if totalProjectsToUpdate == 0 {
        fmt.Println("Everything is up-to-date")
        return false, nil, nil
    }

    fmt.Println("Deploying Projects (" + strconv.Itoa(totalProjectsToUpdate) + ")..." )

    var failedProjects = make(map[string][]*params.ProjectParams)
    var hasDeletedProjects bool
    var deletedProjectsPerType =make(map[string][]*params.ProjectParams)

    // deploying API projects
    apiProjects := updatedProjectsPerType[utils.ProjectTypeApi]
    if len(apiProjects) != 0 {
        fmt.Println("\nAPIs (" + strconv.Itoa(len(apiProjects)) + ") ...")
        for i, projectParam := range apiProjects {
            // if the project is a deleted one, we do it later. So keep it for now.
            if projectParam.Deleted {
                handleProjectDeletion(i, projectParam, deletedProjectsPerType)
                hasDeletedProjects = true
                continue
            }
            importParams := projectParam.ApiParams.Deploy.Import
            fmt.Println(strconv.Itoa(i+1) + ": " + projectParam.NickName + ": (" + projectParam.RelativePath + ")")
            err := impl.ImportAPIToEnv(accessToken, environment, projectParam.AbsolutePath, utils.ParamFileAPI,
                importParams.Update, importParams.PreserveProvider, false, false)
            if err != nil {
                fmt.Println("Error... ", err)
                failedProjects[projectParam.Type] = append(failedProjects[projectParam.Type], projectParam)
            }
        }
    }

    // deploying API product projects
    apiProductProjects := updatedProjectsPerType[utils.ProjectTypeApiProduct]
    if len(apiProductProjects) != 0 {
        fmt.Println("\nAPI Products (" + strconv.Itoa(len(apiProductProjects)) + ") ...")
        for i, projectParam := range apiProductProjects {
            // if the project is a deleted one, we do it later. So keep it for now.
            if projectParam.Deleted {
                handleProjectDeletion(i, projectParam, deletedProjectsPerType)
                hasDeletedProjects = true
                continue
            }
            importParams := projectParam.ApiProductParams.Deploy.Import
            fmt.Println(strconv.Itoa(i+1) + ": " + projectParam.NickName + ": (" + projectParam.RelativePath + ")")
            err := impl.ImportAPIProductToEnv(accessToken, environment, projectParam.AbsolutePath,
                importParams.ImportAPIs, importParams.UpdateAPIs, importParams.UpdateAPIProduct,
                importParams.PreserveProvider, false, false)
            if err != nil {
                fmt.Println("\terror... ", err)
                failedProjects[projectParam.Type] = append(failedProjects[projectParam.Type], projectParam)
            }
        }
    }

    // deploying Application projects
    applicationProjects := updatedProjectsPerType[utils.ProjectTypeApplication]
    if len(applicationProjects) != 0 {
        fmt.Println("\nApplications (" + strconv.Itoa(len(applicationProjects)) + ") ...")
        for i, projectParam := range applicationProjects {
            // if the project is a deleted one, we do it later. So keep it for now.
            if projectParam.Deleted {
                handleProjectDeletion(i, projectParam, deletedProjectsPerType)
                hasDeletedProjects = true
                continue
            }
            importParams := projectParam.ApplicationParams.Deploy.Import
            fmt.Println(strconv.Itoa(i+1) + ": " + projectParam.NickName + ": (" + projectParam.RelativePath + ")")
            _, err := impl.ImportApplicationToEnv(accessToken, environment, projectParam.AbsolutePath,
                importParams.TargetOwner, importParams.Update, importParams.PreserveOwner,
                importParams.SkipSubscriptions, importParams.SkipKeys, false)
            if err != nil {
                fmt.Println("\terror... ", err)
                failedProjects[projectParam.Type] = append(failedProjects[projectParam.Type], projectParam)
            }
        }
    }

    // If there are no deleted projects, update the VCS config file as there is nothing remaining to do.
    //  If there are deleted projects, this needs to handle after deleting those.
    if !hasDeletedProjects {
        updateVCSConfig(repoId, environment, failedProjects)
    }

    return hasDeletedProjects, deletedProjectsPerType, failedProjects
}

// This method is responsible for updating the vcs configuration file at the end of the deployment
// repoId is the id of the git repository (located in vcs.yaml)
// environment is the environment name
// failedProjects are a map of project type to failed projects during the previous deployment
func updateVCSConfig(repoId, environment string, failedProjects map[string][]*params.ProjectParams) {
    vcsConfig, envVCSConfig, _ := getVCSEnvironmentDetails(repoId, environment)
    var err error
    envVCSConfig.LastAttemptedRev, err = getLatestCommitId()
    if err != nil {
        utils.HandleErrorAndExit("Error while getting latest commit-id", err)
    }
    envVCSConfig.FailedProjects = failedProjects

    if len(failedProjects) == 0 {
        if len(envVCSConfig.LastSuccessfulRev) == 0 || len(envVCSConfig.LastSuccessfulRev) > 0 &&
            envVCSConfig.LastSuccessfulRev[0] != envVCSConfig.LastAttemptedRev {
            persistedLast := envVCSConfig.LastSuccessfulRev[0:utils.Min(lastSuccessfulCommitsToKeep-1, len(envVCSConfig.LastSuccessfulRev))]
            envVCSConfig.LastSuccessfulRev = append([]string{envVCSConfig.LastAttemptedRev}, persistedLast...)
        }
    }
    _, hasRepo := vcsConfig.Repos[repoId]
    if !hasRepo {
        vcsConfig.Repos[repoId] = Repo{
            Environments: map[string]Environment{},
        }
    }
    vcsConfig.Repos[repoId].Environments[environment] = envVCSConfig
    utils.WriteConfigFile(vcsConfig, VCSConfigFilePath)
}

// Logs the deletion project info message and appends the project to delete (projectParam) into deletedProjectsPerType map.
// i is the index of the project
// projectParam is the project to be deleted
// deletedProjectsPerType is the map of project type -> projects which are keeping the projects to delete
func handleProjectDeletion(i int, projectParam *params.ProjectParams, deletedProjectsPerType map[string][]*params.ProjectParams) {
    fmt.Println(strconv.Itoa(i+1) + ": " + projectParam.NickName + ": (" + projectParam.RelativePath + ") awaiting deletion..")
    if deletedProjectsPerType[projectParam.Type] == nil {
        deletedProjectsPerType[projectParam.Type] = []*params.ProjectParams{}
    }
    deletedProjectsPerType[projectParam.Type] = append(deletedProjectsPerType[projectParam.Type], projectParam)
}

// Scan and detects all the changes in projects by comparing the current revision with the last successful revision.
// Deploy all the changes to the specified environment.
// accesstoken is the access token to access the APIM product REST APIs
// environment is the environment name
func DeployChangedFiles(accessToken, environment string) map[string][]*params.ProjectParams {
    repoId, totalProjectsToUpdate, updatedProjectsPerType := GetStatus(environment, FromRevTypeLastAttempted)
    hasDeletedProjects, deletedProjectsPerType, failedProjects :=
        deployUpdatedProjects(accessToken, repoId, environment, totalProjectsToUpdate, updatedProjectsPerType)

    if hasDeletedProjects {
        //check whether project deletion is disabled
        mainConfig := utils.GetMainConfigFromFile(utils.MainConfigFilePath)
        if !mainConfig.Config.VCSDeletionEnabled {
            utils.HandleErrorAndExit("Error: there are projects to delete while project " +
                "deletion is disabled via VCS", nil)
        }

        // work on deleted files
        _, envVCSConfig, hasEnv := getVCSEnvironmentDetails(repoId, environment)
        if !hasEnv || len(envVCSConfig.LastSuccessfulRev) == 0 {
            utils.HandleErrorAndExit("Error: there are projects to delete but no last successful "+
                "revision available in vcs config (vcs_config.yaml)", nil)
            return nil
        }
        currentBranch := getCurrentBranch()
        lastSuccessfulRev := envVCSConfig.LastSuccessfulRev[0]
        tmpBranchName := "tmp-" + lastSuccessfulRev[0:8]

        fmt.Println("\nDeleting projects ..")
        checkoutNewBranchFromRevision(tmpBranchName, lastSuccessfulRev)
        failedProjects = deployProjectDeletions(accessToken, environment, deletedProjectsPerType, failedProjects)
        checkoutBranch(currentBranch)
        deleteTmpBranch(tmpBranchName)

        // Update the VCS config with failed projects, last attempted and last successful revisions
        updateVCSConfig(repoId, environment, failedProjects)
    }
    return failedProjects
}

// Create 'vcs.yaml' in the repository root folder with a unique id (uuid) for the repository.
//  If the value of force is false, and the file is already created, gives an error.
//  If the value of force is true, It will reinitialize the file even if it already exists.
func InitializeRepo(force bool) error {
    vcsInfoPath, err := getVcsYamlPath()
    if err != nil {
        return err
    }
    if !force && utils.IsFileExist(vcsInfoPath) {
        return errors.New("the repository is already initialized")
    }
    repoInfo := RepoInfo{
        Id: uuid.New().String(),
    }
    utils.WriteConfigFile(repoInfo, vcsInfoPath)
    return nil
}

// Returns the absolute path of the vcs.yaml in the current working repository
func getVcsYamlPath() (string, error) {
    baseDir, err := getRepoBaseDir()
    if err != nil {
        return "", err
    }
    vcsInfoPath := filepath.Join(baseDir, VCSRepoInfoFileName)
    return vcsInfoPath, nil
}

// Returns the id of the current working repository by reading vcs.yaml
func getRepoId() (string, error) {
    vcsInfoPath, err := getVcsYamlPath()
    if err != nil {
        return "", err
    }
    var repoInfo RepoInfo
    data, err := ioutil.ReadFile(vcsInfoPath)
    if err == nil {
        if err := yaml.Unmarshal(data, &repoInfo); err != nil {
            utils.HandleErrorAndExit("Error parsing "+vcsInfoPath, err)
        }
    }
    return repoInfo.Id, nil
}

// Retrieves the base location of the git repository
func getRepoBaseDir() (string, error) {
    baseDir, err := executeGitCommand("rev-parse", "--show-toplevel")
    if err != nil {
        return "", err
    }
    return strings.TrimSpace(baseDir), nil
}

// Retrieves the latest commit of the git repository
func getLatestCommitId() (string, error) {
    latestCommit, err := executeGitCommand("rev-parse", "HEAD")
    if err != nil {
        return "", err
    }
    return strings.TrimSpace(latestCommit), nil
}

// Prints the changed files in the terminal
func logChangedFiles(changedFileList []string) {
    utils.Logln("Total changed files: " + strconv.Itoa(len(changedFileList)))
    for i, changedFile := range changedFileList {
        utils.Logln(strconv.Itoa(i+1) + ": " + changedFile)
    }
    utils.Logln()
}

// subPath denotes a single changed file retrieved from "git diff" command. This path is scanned and then identifies the
//  project type (API, App,.. ) from this method.
// envVCSConfig is the environment related VCS configuration
// repoBasePath is the basepath of the git repository
// pathInfoMap is a map of path (string) to project info. This is used for caching and avoid repetitive checking
// Returns the identified project details. If it is not related to a project, a NONE project info item will be returned
func getProjectInfoFromProjectFile(envVCSConfig Environment, repoBasePath string, subPath string,
    pathInfoMap map[string]*params.ProjectParams) *params.ProjectParams {
    subPaths := getSubPaths(repoBasePath, subPath)
    for _, s := range subPaths {
        projectParams := checkProjectTypeOfSpecificPath(repoBasePath, s, pathInfoMap)
        if projectParams.Type != utils.ProjectTypeNone {
            // once we identified the project type, check whether the project is failed previously. If so, mark it as
            //  failed. This is used to show failed projects by the "status" command.
            projectParams.FailedDuringPreviousDeploy = failedDuringEarlierDeploy(envVCSConfig, projectParams)
            return projectParams
        }
    }
    return &params.ProjectParams{
        Type: utils.ProjectTypeNone,
    }
}

// The files in the path "fullPath" is scanned and identifies the project type (API, App,.. ) from this method.
// repoBasePath is the basepath of the git repository
// pathInfoMap is a map of path (string) to project info. This is used for caching and avoid repetitive checking
// Returns the identified project details. If it is not related to a project, a NONE project info item will be returned
func checkProjectTypeOfSpecificPath(repoBasePath, fullPath string,
    pathInfoMap map[string]*params.ProjectParams) *params.ProjectParams {
    if pathInfoMap[fullPath] != nil {
        return pathInfoMap[fullPath]
    }

    var projectParams = &params.ProjectParams{
        Type:         utils.ProjectTypeNone,
        AbsolutePath: fullPath,
        RelativePath: strings.Replace(fullPath, repoBasePath + string(os.PathSeparator), "", 1),
        NickName:     filepath.Base(fullPath),
    }

    // in case fullPath contains a deleted file/folder, the path does not exist.
    files, err := ioutil.ReadDir(fullPath)
    if err != nil && os.IsNotExist(err) {
        // if the path doesn't exist, mark it as deleted
        projectParams.Deleted = true

        // checks if fullPath represents a *_params.yaml file, then set the project type accordingly
        if strings.HasSuffix(fullPath, utils.ParamFileAPI) {
            projectParams.Type = utils.ProjectTypeApi
        }
        if strings.HasSuffix(fullPath, utils.ParamFileAPIProduct) {
            projectParams.Type = utils.ProjectTypeApiProduct
        }
        if strings.HasSuffix(fullPath, utils.ParamFileApplication) {
            projectParams.Type = utils.ProjectTypeApplication
        }
        //This means project type is set from any of the above condition.
        //  Then set the correct basePath of the project.
        if projectParams.Type != utils.ProjectTypeNone {
            // remove the *_params.yaml part from the paths
            projectParams.RelativePath, _ = filepath.Split(projectParams.RelativePath)
            projectParams.AbsolutePath, _ = filepath.Split(projectParams.AbsolutePath)
            //remove "/" suffix if exists
            if strings.HasSuffix(projectParams.RelativePath, string(os.PathSeparator)) {
                projectParams.RelativePath = strings.TrimSuffix(projectParams.RelativePath, string(os.PathSeparator))
                projectParams.AbsolutePath = strings.TrimSuffix(projectParams.AbsolutePath, string(os.PathSeparator))
            }
            projectParams.NickName = filepath.Base(projectParams.RelativePath)
        }
        // return the projectParams as a deleted project
        return projectParams
    }

    //If the path exists (checked previously), read through the file names of the specific path and check for
    //  *_params.yaml to determine the project type
    for _, f := range files {
        fullPathWithFileName := filepath.Join(fullPath, f.Name())
        switch f.Name() {
        case utils.ParamFileAPI:
            projectParams.Type = utils.ProjectTypeApi;
            apiParams, err := params.LoadApiParamsFromFile(fullPathWithFileName)
            if err != nil {
                utils.HandleErrorAndExit("Error while parsing " + utils.ParamFileAPI + " file:" +fullPathWithFileName, err)
            }
            projectParams.ApiParams = apiParams
            break
        case utils.ParamFileAPIProduct:
            projectParams.Type = utils.ProjectTypeApiProduct;
            apiProductParams, err := params.LoadApiProductParamsFromFile(fullPathWithFileName)
            if err != nil {
                utils.HandleErrorAndExit("Error while parsing " + utils.ParamFileAPIProduct + " file:" +fullPathWithFileName, err)
            }
            projectParams.ApiProductParams = apiProductParams
            break
        case utils.ParamFileApplication:
            projectParams.Type = utils.ProjectTypeApplication;
            applicationParams, err := params.LoadApplicationParamsFromFile(fullPathWithFileName)
            if err != nil {
                utils.HandleErrorAndExit("Error while parsing " + utils.ParamFileApplication + " file:" +fullPathWithFileName, err)
            }
            projectParams.ApplicationParams = applicationParams
            break
        }
        if projectParams.Type != utils.ProjectTypeNone {
            //breaks from for loop
            break
        }
    }
    pathInfoMap[fullPath] = projectParams
    return projectParams
}

// Returns all the subpaths as an array from a given path
func getSubPaths(parent string, path string) (paths []string) {
    var subPaths []string
    folders := strings.Split(path, string(os.PathSeparator))
    nextPath := parent
    for _, folder := range folders {
        nextPath = filepath.Join(nextPath, folder)
        subPaths = append(subPaths, nextPath)
    }
    return subPaths
}

// Executes the give git command as args list and returns the output
func executeGitCommand(args ...string) (string, error) {
    cmd := exec.Command(Git, args...)

    if utils.VerboseModeEnabled() {
        utils.Logln("Executing command: " + Git + " " + strings.Join(args, " "))
    }

    var errBuf bytes.Buffer
    cmd.Stderr = io.MultiWriter(os.Stderr, &errBuf)

    output, err := cmd.Output()

    if utils.VerboseModeEnabled() {
        if err != nil {
            utils.HandleErrorAndContinue("Error occurred while executing command: ", err)
        } else {
            utils.Logln("Output : " + string(output))
        }
    }
    return string(output), err
}
