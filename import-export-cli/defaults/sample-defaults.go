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

package defaults

var SampleConfig = []byte(`config: 
  export_directory: /home/kasun/.wso2apimcli/exported
  http_request_timeout: 10000
environments: 
  sample-env1: 
    admin_endpoint: "https://localhost:9443/api/am/admin/v0.14"
    api_import_export_endpoint: "https://localhost/api-import-export"
    api_list_endpoint: "https://localhost/publisher/apis"
    api_manager_endpoint: "https://localhost/apim"
    application_list_endpoint: "https://localhost:9443/api/am/admin/v0.14/applications"
    registration_endpoint: "https://localhost/register"
    token_endpoint: "https://localhost/token"
  sample-env2: 
    admin_endpoint: ""
    api_import_export_endpoint: ""
    api_list_endpoint: ""
    api_manager_endpoint: "https://localhost/apim"
    application_list_endpoint: ""
    registration_endpoint: "https://localhost/register"
    token_endpoint: "https://localhost/token"
`)
