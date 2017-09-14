package utils

// File Names and Paths
const ApplicationRoot string = "/home/menuka/.go/src/github.com/menuka94/wso2apim-cli"
const EnvKeysAllFileName string = "env_keys_all.yaml"
const EnvKeysAllFilePath string = ApplicationRoot + "/" + EnvKeysAllFileName
const EnvEndpointsAllFileName string = "env_endpoints_all.yaml"
const EnvEndpointsAllFilePath string = ApplicationRoot +  "/" + EnvEndpointsAllFileName
const ExportedAPIsDirectoryName string = "exported"
const ExportedAPIsDirectoryPath string = ApplicationRoot + "/" + ExportedAPIsDirectoryName

// Headers and Header Values
const HeaderAuthorization string = "Authorization"
const HeaderContentType string = "Content-Type"
const HeaderConnection string = "Connection"
const HeaderAccept string = "Accept"
const HeaderProduces string = "Produces"
const HeaderConsumes string = "Consumes"
const HeaderContentEncoding string = "Content-Encoding"
const HeaderTransferEncoding string = "transfer-encoding"
const HeaderValueChunked string = "chunked"
const HeaderValueGZIP string = "gzip"
const HeaderValueKeepAlive string = "keep-alive"
const HeaderValueApplicationZip = "application/zip"
const HeaderValueApplicationJSON string = "application/json"
const HeaderValueXWWWFormUrlEncoded string = "application/x-www-form-urlencoded"
const HeaderValueAuthBearerPrefix string = "Bearer"
const HeaderValueAuthBasicPrefix string = "Basic"
const HeaderValueMultiPartFormData string = "multipart/form-data"

// Logging Prefixes
const LogPrefixInfo = "[INFO]: "
const LogPrefixWarning = "[WARN]: "
const LogPrefixError = "[ERROR]: "

// Other
const APICallTimeout = 5
const WSO2APIMUpdateTokenTimeout = 2
const DefaultTokenValidityPeriod string = "3600"

