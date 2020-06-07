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
    "io"
    "io/ioutil"
    "os"
    "os/exec"
    "path/filepath"
    "strconv"
    "strings"
)

func GetChangedFiles(accessToken, environment string) {
    changedFiles, _ := executeGitCommand("diff", "--name-only")
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
        projectParam := getProjectInfoFromProjectFile("/home/malintha/wso2apim/cur/apictl/gitint/repo2",
            changedFile, changedPathInfoMap)
        if projectParam.Type != utils.ProjectTypeNone {
            if updatedProjectsPerType[projectParam.Type] == nil {
                updatedProjectsPerType[projectParam.Type] = []*params.ProjectParams{}
            }
            if updatedProjectsPerProjectPath[projectParam.BasePath] == nil {
                updatedProjectsPerProjectPath[projectParam.BasePath] = projectParam
                updatedProjectsPerType[projectParam.Type] = append(updatedProjectsPerType[projectParam.Type], projectParam)
                totalProjectsToUpdate++
            }
        }
    }

    if totalProjectsToUpdate == 0 {
        fmt.Println("Everything is up-to-date")
        return
    }

    fmt.Println("Updating Projects (" + strconv.Itoa(totalProjectsToUpdate) + ")..." )
    for projectType, projectParams := range updatedProjectsPerType {
        if projectParams != nil && len(projectParams) > 0 {
            if projectType == utils.ProjectTypeApi {
                for i, projectParam := range projectParams {
                    importParams := projectParam.ApiParams.Import
                    fmt.Println("\n" + strconv.Itoa(i + 1) + ": " + projectParam.Name + ": \t")
                    err := impl.ImportAPI(accessToken, environment, projectParam.BasePath, "",
                        importParams.Update, importParams.PreserveProvider, false)
                    if err != nil {
                        utils.Logln("\terror... ", err)
                    }
                }
            }
        }
    }
}

func logChangedFiles(changedFileList []string) {
    utils.Logln("Total changed files: " + strconv.Itoa(len(changedFileList)))
    for i, changedFile := range changedFileList {
        utils.Logln(strconv.Itoa(i+1) + ": " + changedFile)
    }
    utils.Logln()
}

func getProjectInfoFromProjectFile(basePath string, subPath string, pathInfoMap map[string]*params.ProjectParams) *params.ProjectParams {
    subPaths := getSubPaths(basePath, subPath)
    for _, s := range subPaths {
        projectParams := checkProjectTypeOfSpecificPath(s, pathInfoMap)
        if projectParams.Type != utils.ProjectTypeNone {
            return projectParams
        }
    }
    return &params.ProjectParams{
        Type: utils.ProjectTypeNone,
    }
}

func checkProjectTypeOfSpecificPath(fullPath string, pathInfoMap map[string]*params.ProjectParams) *params.ProjectParams {
    if pathInfoMap[fullPath] != nil {
        return pathInfoMap[fullPath]
    }

    files, err := ioutil.ReadDir(fullPath)
    if err != nil {
        utils.HandleErrorAndExit("cannot open path " + fullPath + " for checking project type", err)
    }
    var projectParams = &params.ProjectParams{
        Type:     utils.ProjectTypeNone,
        BasePath: fullPath,
        Name:     filepath.Base(fullPath),
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
                utils.HandleErrorAndExit("Error while parsing " + utils.ParamFileAPI + " file:" +fullPathWithFileName, err)
            }
            projectParams.ApiProductParams = apiProductParams
            break
        case utils.ParamFileApplication:
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
    folderPath, _ := filepath.Split(path)
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
