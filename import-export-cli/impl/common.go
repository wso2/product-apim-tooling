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
