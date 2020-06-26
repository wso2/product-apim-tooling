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
    "fmt"
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

// Read and return MainConfig. Silently catch the error  when config file is not found
func GetVCSConfigFromFileSilently(filePath string) *VCSConfig {
    var vcsConfig VCSConfig
    data, err := ioutil.ReadFile(filePath)
    if err == nil {
        if err := yaml.Unmarshal(data, &vcsConfig); err != nil {
            utils.HandleErrorAndExit("VCSConfig: Error parsing "+filePath, err)
        }
    }
    return &vcsConfig
}

func getVCSEnvironmentDetails(environment string) (VCSConfig, Environment, bool)  {
    vcsConfig := GetVCSConfigFromFileSilently(VCSConfigFilePath)
    if vcsConfig.Environments == nil {
        vcsConfig.Environments = make(map[string]Environment)
    }
    envVCSConfig, hasEnv := vcsConfig.Environments[environment]
    return *vcsConfig, envVCSConfig, hasEnv
}

func GetStatus(environment, fromRevType string) (int, map[string][]*params.ProjectParams){
    var envRevision string
    _, envVCSConfig, hasEnv := getVCSEnvironmentDetails(environment)
    if hasEnv {
        if fromRevType == FromRevTypeLastAttempted {
            envRevision = envVCSConfig.LastAttemptedRev
        } else if fromRevType == FromRevTypeLastSuccessful{
            envRevision = envVCSConfig.LastSuccessfulRev
        }
    }

    basePath, err := getRepoBaseDir()
    if err != nil {
        utils.HandleErrorAndExit("Error while getting repository base folder location", err)
    }

    var changedFiles string
    if envRevision == "" {
        changedFiles, _ = executeGitCommand("ls-tree", "-r", "HEAD", "--name-only")
    } else {
        changedFiles, _ = executeGitCommand("diff", "--name-only", envRevision)
    }
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
                failedProjectInEachType.FailedDuringPreviousPush = true
                totalProjectsToUpdate++
            }
        }
    }

    return totalProjectsToUpdate, updatedProjectsPerType
}

func failedDuringEarlierPush(vcsEnvConfig Environment, projectParams *params.ProjectParams) bool {
    failedProjectsForType := vcsEnvConfig.FailedProjects[projectParams.Type]
    for _, failedProject := range failedProjectsForType {
        if failedProject.RelativePath == projectParams.RelativePath {
            return true
        }
    }
    return false
}

func Rollback(accessToken, environment string) {
    totalProjectsToUpdate, updatedProjectsPerType := GetStatus(environment, FromRevTypeLastSuccessful)
    _, envVCSConfig, hasEnv := getVCSEnvironmentDetails(environment)

    if !hasEnv || envVCSConfig.LastSuccessfulRev == "" || len(envVCSConfig.FailedProjects) == 0{
        fmt.Println("Nothing to rollback")
        return
    }
    currentBranch := getCurrentBranch()
    tmpBranchName := "tmp-" + envVCSConfig.LastSuccessfulRev[0:8]
    checkoutNewBranchFromRevision(tmpBranchName, envVCSConfig.LastSuccessfulRev)
    pushUpdatedProjects(accessToken, environment, totalProjectsToUpdate, updatedProjectsPerType)
    checkoutBranch(currentBranch)
    deleteTmpBranch(tmpBranchName)
}

func checkoutNewBranchFromRevision(tmpBranchName, revision string) {
    _, err := executeGitCommand("checkout", "-b", tmpBranchName, revision)
    if err != nil {
        utils.HandleErrorAndExit("Error while checking out last successful commit ("+
            revision+") for rolling back", err)
    }
}

func checkoutBranch(branchName string) {
    _, err := executeGitCommand("checkout", branchName)
    if err != nil {
        utils.HandleErrorAndExit("Error while checking out branch " + branchName, err)
    }
}

func getCurrentBranch() string {
    branch, err := executeGitCommand("rev-parse", "--abbrev-ref", "HEAD")
    if err != nil {
        utils.HandleErrorAndExit("Error while getting current branch", err)
    }
    return strings.TrimSpace(branch)
}

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

func pushDeletedProjects(accessToken, environment string, deletedProjectsPerType map[string][]*params.ProjectParams,
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
            resp, err := impl.DeleteApplication(accessToken, environment, appInfo.Name)
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

func handleIfError(err error, failedProjects map[string][]*params.ProjectParams, projectParam *params.ProjectParams) bool {
    if err != nil {
        fmt.Println("Error... ", err)
        failedProjects[projectParam.Type] = append(failedProjects[projectParam.Type], projectParam)
    }
    return err != nil
}

func pushUpdatedProjects(accessToken, environment string, totalProjectsToUpdate int,
        updatedProjectsPerType map[string][]*params.ProjectParams) (bool, map[string][]*params.ProjectParams,
        map[string][]*params.ProjectParams) {
    if totalProjectsToUpdate == 0 {
        fmt.Println("Everything is up-to-date")
        return false, nil, nil
    }

    fmt.Println("Updating Projects (" + strconv.Itoa(totalProjectsToUpdate) + ")..." )

    var failedProjects = make(map[string][]*params.ProjectParams)
    var hasDeletedProjects bool
    var deletedProjectsPerType =make(map[string][]*params.ProjectParams)

    // pushing API projects
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
            importParams := projectParam.ApiParams.Import
            fmt.Println(strconv.Itoa(i+1) + ": " + projectParam.NickName + ": (" + projectParam.RelativePath + ")")
            err := impl.ImportAPIToEnv(accessToken, environment, projectParam.AbsolutePath, "",
                importParams.Update, importParams.PreserveProvider, false)
            if err != nil {
                fmt.Println("Error... ", err)
                failedProjects[projectParam.Type] = append(failedProjects[projectParam.Type], projectParam)
            }
        }
    }

    // pushing API product projects
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
            importParams := projectParam.ApiProductParams.Import
            fmt.Println(strconv.Itoa(i+1) + ": " + projectParam.NickName + ": (" + projectParam.RelativePath + ")")
            err := impl.ImportAPIProductToEnv(accessToken, environment, projectParam.AbsolutePath,
                importParams.ImportAPIs, importParams.UpdateAPIs, importParams.UpdateAPIProduct,
                importParams.PreserveProvider, false)
            if err != nil {
                fmt.Println("\terror... ", err)
                failedProjects[projectParam.Type] = append(failedProjects[projectParam.Type], projectParam)
            }
        }
    }

    // pushing Application projects
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
            importParams := projectParam.ApplicationParams.Import
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
        updateVCSConfig(environment, failedProjects)
    }

    return hasDeletedProjects, deletedProjectsPerType, failedProjects
}

func updateVCSConfig(environment string, failedProjects map[string][]*params.ProjectParams) {
    vcsConfig, envVCSConfig, _ := getVCSEnvironmentDetails(environment)

    var err error
    envVCSConfig.LastAttemptedRev, err = getLatestCommitId()
    if err != nil {
        utils.HandleErrorAndExit("Error while getting latest commit-id", err)
    }
    envVCSConfig.FailedProjects = failedProjects

    if len(failedProjects) == 0 {
        envVCSConfig.LastSuccessfulRev = envVCSConfig.LastAttemptedRev
    }

    vcsConfig.Environments[environment] = envVCSConfig
    utils.WriteConfigFile(vcsConfig, VCSConfigFilePath)
}

func handleProjectDeletion(i int, projectParam *params.ProjectParams, deletedProjectsPerType map[string][]*params.ProjectParams) {
    fmt.Println(strconv.Itoa(i+1) + ": " + projectParam.NickName + ": (" + projectParam.RelativePath + ") awaiting deletion..")
    if deletedProjectsPerType[projectParam.Type] == nil {
        deletedProjectsPerType[projectParam.Type] = []*params.ProjectParams{}
    }
    deletedProjectsPerType[projectParam.Type] = append(deletedProjectsPerType[projectParam.Type], projectParam)
}

func PushChangedFiles(accessToken, environment string) {
    totalProjectsToUpdate, updatedProjectsPerType := GetStatus(environment, FromRevTypeLastAttempted)
    hasDeletedProjects, deletedProjectsPerType, failedProjects :=
        pushUpdatedProjects(accessToken, environment, totalProjectsToUpdate, updatedProjectsPerType)

    if hasDeletedProjects {
        // work on deleted files
        _, envVCSConfig, hasEnv := getVCSEnvironmentDetails(environment)
        if !hasEnv || envVCSConfig.LastSuccessfulRev == "" {
            utils.HandleErrorAndExit("Error: there are projects to delete by no last successful "+
                "revision available in vcs config (vcs_config.yaml)", nil)
            return
        }
        currentBranch := getCurrentBranch()
        tmpBranchName := "tmp-" + envVCSConfig.LastSuccessfulRev[0:8]

        fmt.Println("\nDeleting projects ..")
        checkoutNewBranchFromRevision(tmpBranchName, envVCSConfig.LastSuccessfulRev)
        failedProjects = pushDeletedProjects(accessToken, environment, deletedProjectsPerType, failedProjects)
        checkoutBranch(currentBranch)
        deleteTmpBranch(tmpBranchName)

        // Update the VCS config with failed projects, last attempted and last successful revisions
        updateVCSConfig(environment, failedProjects)
    }
}

func getRepoBaseDir() (string, error) {
    baseDir, err := executeGitCommand("rev-parse", "--show-toplevel")
    if err != nil {
        return "", err
    }
    return strings.TrimSpace(baseDir), nil
}

func getLatestCommitId() (string, error) {
    latestCommit, err := executeGitCommand("rev-parse", "HEAD")
    if err != nil {
        return "", err
    }
    return strings.TrimSpace(latestCommit), nil
}

func logChangedFiles(changedFileList []string) {
    utils.Logln("Total changed files: " + strconv.Itoa(len(changedFileList)))
    for i, changedFile := range changedFileList {
        utils.Logln(strconv.Itoa(i+1) + ": " + changedFile)
    }
    utils.Logln()
}

func getProjectInfoFromProjectFile(envVCSConfig Environment, repoBasePath string, subPath string, pathInfoMap map[string]*params.ProjectParams) *params.ProjectParams {
    subPaths := getSubPaths(repoBasePath, subPath)
    for _, s := range subPaths {
        projectParams := checkProjectTypeOfSpecificPath(repoBasePath, s, pathInfoMap)
        if projectParams.Type != utils.ProjectTypeNone {
            // once we identified the project type, check whether the project is failed previously. If so, mark it as
            //  failed. This is used to show failed projects by the "status" command.
            projectParams.FailedDuringPreviousPush = failedDuringEarlierPush(envVCSConfig, projectParams)
            return projectParams
        }
    }
    return &params.ProjectParams{
        Type: utils.ProjectTypeNone,
    }
}

func checkProjectTypeOfSpecificPath(repoBasePath, fullPath string, pathInfoMap map[string]*params.ProjectParams) *params.ProjectParams {
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

// ExecuteCommand executes the command with args and prints output, errors in standard output, error
func executeGitCommand(args ...string) (string, error) {
    cmd := exec.Command(Git, args...)
    var errBuf bytes.Buffer
    cmd.Stderr = io.MultiWriter(os.Stderr, &errBuf)

    output, err := cmd.Output()
    return string(output), err
}


// GetCommandOutput executes a command and returns the output
func GetCommandOutput(command string, args ...string) (string, error) {
    cmd := exec.Command(command, args...)
    var errBuf bytes.Buffer
    cmd.Stderr = io.MultiWriter(os.Stderr, &errBuf)

    output, err := cmd.Output()
    return string(output), err
}

// setCommandOutAndError sets the output and error of the command cmd to the standard output and error
func setCommandOutAndError(cmd *exec.Cmd) {
    var errBuf, outBuf bytes.Buffer
    cmd.Stderr = io.MultiWriter(os.Stderr, &errBuf)
    cmd.Stdout = io.MultiWriter(os.Stdout, &outBuf)
}
