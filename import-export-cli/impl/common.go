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

package impl

import (
    "bytes"
    "errors"
    "github.com/go-resty/resty"
    "github.com/wso2/product-apim-tooling/import-export-cli/box"
    "github.com/wso2/product-apim-tooling/import-export-cli/utils"
    "io"
    "mime/multipart"
    "os"
    "path/filepath"
    "strings"
    "text/template"
)

// ExecuteNewFileUploadRequest forms an HTTP request
// Helper function for forming multi-part form data
// Returns the formed http request and errors
func ExecuteNewFileUploadRequest(uri string, params map[string]string, paramName, path,
    accessToken string) (*resty.Response, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)
    part, err := writer.CreateFormFile(paramName, filepath.Base(path))
    if err != nil {
        return nil, err
    }
    _, err = io.Copy(part, file)

    for key, val := range params {
        _ = writer.WriteField(key, val)
    }
    err = writer.Close()
    if err != nil {
        return nil, err
    }

    // Set headers
    headers := make(map[string]string)
    headers[utils.HeaderContentType] = writer.FormDataContentType()
    headers[utils.HeaderAuthorization] =  utils.HeaderValueAuthBearerPrefix+" "+accessToken
    headers[utils.HeaderAccept] = "application/json"
    headers[utils.HeaderConnection] = utils.HeaderValueKeepAlive

    resp, err := utils.InvokePOSTRequestWithBytes(uri, headers, body.Bytes())

    return resp, err
}

// Include x_params.yaml (api_params.yaml, application_params.yaml, .. ) into the sourceZipFile and create a new
//  new Zip file in the provided targetZipFile location. paramsFile needs to be one of the supported x_params.yaml.
//  Eg.: api_params.yaml, application_params.yaml, api_product_params.yaml
func IncludeParamsFileToZip(sourceZipFile, targetZipFile, paramsFile string) error {
    // Now, we need to extract the zip, copy x_params.yaml file inside and then create the zip again
    //	First, create a temp directory (tmpClonedLoc) by extracting the original zip file.
    tmpClonedLoc, err := utils.GetTempCloneFromDirOrZip(sourceZipFile)
    // Create the api_params.yaml file inside the cloned directory.
    tmpLocationForAPIParamsFile := filepath.Join(tmpClonedLoc, paramsFile)
    err = ScaffoldParams(tmpLocationForAPIParamsFile)
    if err != nil {
        utils.HandleErrorAndExit("Error creating api_params.yaml inside the exported zip archive", err)
    }

    err = utils.Zip(tmpClonedLoc, targetZipFile)
    if err != nil {
        utils.HandleErrorAndExit("Error creating the final zip archive", err)
    }
    return nil
}

// Creates the initial api_params.yaml/api_product_params.yaml/application_params.yaml in the given file path
//	The targetFile will be populated with environments and default import parameters for "vcs deploy".
func ScaffoldParams(targetFile string) error {
    envs := utils.GetMainConfigFromFile(utils.MainConfigFilePath)
    var tmpl []byte
    if strings.HasSuffix(targetFile, utils.ParamFileAPI) {
        tmpl, _ = box.Get("/init/api_params.tmpl")
    } else if strings.HasSuffix(targetFile, utils.ParamFileAPIProduct) {
        tmpl, _ = box.Get("/init/api_product_params.tmpl")
    } else if strings.HasSuffix(targetFile, utils.ParamFileApplication) {
        tmpl, _ = box.Get("/init/application_params.tmpl")
    } else {
        return errors.New("Unsupported target file: " + targetFile)
    }
    return WriteTargetFileFromTemplate(targetFile, tmpl, envs)
}

// From the template data (tmpl) writes the target file using the provided mainConfig
func WriteTargetFileFromTemplate(targetFile string, tmpl []byte, envs *utils.MainConfig) error {
    t, err := template.New("").Parse(string(tmpl))
    if err != nil {
        return err
    }

    f, err := os.Create(targetFile)
    if err != nil {
        return err
    }
    defer f.Close()

    err = t.Execute(f, envs.Environments)
    if err != nil {
        return err
    }
    return nil
}
