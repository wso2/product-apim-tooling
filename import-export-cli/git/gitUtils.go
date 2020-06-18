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

    //append failed projects to the list of changed files if exists
    for _, failedProjectsInEachType := range envVCSConfig.FailedProjects {
        changedFileList = append(changedFileList, failedProjectsInEachType...)
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

    return totalProjectsToUpdate, updatedProjectsPerType
}

func failedDuringEarlierPush(vcsEnvConfig Environment, projectParams *params.ProjectParams) bool {
    failedProjectsForType := vcsEnvConfig.FailedProjects[projectParams.Type]
    for _, failedProject := range failedProjectsForType {
        if failedProject == projectParams.RelativePath {
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
    pushChangedFiles(accessToken, environment, totalProjectsToUpdate, updatedProjectsPerType)
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
    if !strings.HasPrefix(tmpBranch, "tmp-") {
        utils.HandleErrorAndExit("Cannot remove branches not starting with 'tmp-'", nil)
    }
    _, err := executeGitCommand("branch", "-D", tmpBranch)
    if err != nil {
        utils.HandleErrorAndExit("Error while deleting the temp branch "+tmpBranch, err)
    }
}

func pushChangedFiles(accessToken, environment string, totalProjectsToUpdate int,
        updatedProjectsPerType map[string][]*params.ProjectParams) {
    if totalProjectsToUpdate == 0 {
        fmt.Println("Everything is up-to-date")
        return
    }

    fmt.Println("Updating Projects (" + strconv.Itoa(totalProjectsToUpdate) + ")..." )

    var failedProjects = make(map[string][]string)
    // pushing API projects
    apiProjects := updatedProjectsPerType[utils.ProjectTypeApi]
    if len(apiProjects) != 0 {
        fmt.Println("\nAPIs (" + strconv.Itoa(len(apiProjects)) + ") ...")
        for i, projectParam := range apiProjects {
            importParams := projectParam.ApiParams.Import
            fmt.Println(strconv.Itoa(i + 1) + ": " + projectParam.Name + ": (" + projectParam.RelativePath + ")")
            err := impl.ImportAPIToEnv(accessToken, environment, projectParam.AbsolutePath, "",
                importParams.Update, importParams.PreserveProvider, false)
            if err != nil {
                fmt.Println("Error... ", err)
                failedProjects[projectParam.Type] = append(failedProjects[projectParam.Type], projectParam.RelativePath)
            }
        }
    }

    // pushing API product projects
    apiProductProjects := updatedProjectsPerType[utils.ProjectTypeApiProduct]
    if len(apiProductProjects) != 0 {
        fmt.Println("\nAPI Products (" + strconv.Itoa(len(apiProductProjects)) + ") ...")
        for i, projectParam := range apiProductProjects {
            importParams := projectParam.ApiProductParams.Import
            fmt.Println(strconv.Itoa(i + 1) + ": " + projectParam.Name + ": (" + projectParam.RelativePath + ")")
            err := impl.ImportAPIProductToEnv(accessToken, environment, projectParam.AbsolutePath,
                importParams.ImportAPIs, importParams.UpdateAPIs, importParams.UpdateAPIProduct,
                importParams.PreserveProvider, false)
            if err != nil {
                fmt.Println("\terror... ", err)
            }
        }
    }

    // pushing Application projects
    applicationProjects := updatedProjectsPerType[utils.ProjectTypeApplication]
    if len(applicationProjects) != 0 {
        fmt.Println("\nApplications (" + strconv.Itoa(len(applicationProjects)) + ") ...")
        for i, projectParam := range applicationProjects {
            importParams := projectParam.ApplicationParams.Import
            fmt.Println(strconv.Itoa(i + 1) + ": " + projectParam.Name + ": (" + projectParam.RelativePath + ")")
            err := impl.ImportApplicationToEnv(accessToken, environment, projectParam.AbsolutePath,
                importParams.TargetOwner, importParams.Update, importParams.PreserveOwner,
                importParams.SkipSubscriptions, importParams.SkipKeys, false)
            if err != nil {
                fmt.Println("\terror... ", err)
            }
        }
    }
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

func PushChangedFiles(accessToken, environment string) {
    totalProjectsToUpdate, updatedProjectsPerType := GetStatus(environment, FromRevTypeLastAttempted)
    pushChangedFiles(accessToken, environment, totalProjectsToUpdate, updatedProjectsPerType)
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

    files, err := ioutil.ReadDir(fullPath)
    if err != nil {
        utils.HandleErrorAndExit("cannot open path " + fullPath + " for checking project type", err)
    }
    var projectParams = &params.ProjectParams{
        Type:         utils.ProjectTypeNone,
        AbsolutePath: fullPath,
        RelativePath: strings.Replace(fullPath, repoBasePath + string(os.PathSeparator), "", 1),
        Name:         filepath.Base(fullPath),
    }

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
    folderPath := path
    pathInfo, err := os.Stat(filepath.Join(parent, path))
    if err != nil {
        utils.HandleErrorAndExit("Error while checking details of path " + path, err)
    }
    if !pathInfo.IsDir() {
        folderPath, _ = filepath.Split(path)
    }
    folders := strings.Split(folderPath, string(os.PathSeparator))
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
