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
    "github.com/wso2/product-apim-tooling/import-export-cli/utils"
    "io"
    "mime/multipart"
    "net/http"
    "os"
    "path/filepath"
)

// newFileUploadRequest forms an HTTP request
// Helper function for forming multi-part form data
// Returns the formed http request and errors
func NewFileUploadRequest(uri string, method string, params map[string]string, paramName, path,
    accessToken string) (*http.Request, error) {
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

    request, err := http.NewRequest(method, uri, body)
    request.Header.Add(utils.HeaderAuthorization, utils.HeaderValueAuthBearerPrefix+" "+accessToken)
    request.Header.Add(utils.HeaderContentType, writer.FormDataContentType())
    request.Header.Add(utils.HeaderAccept, "*/*")
    request.Header.Add(utils.HeaderConnection, utils.HeaderValueKeepAlive)

    return request, err
}
