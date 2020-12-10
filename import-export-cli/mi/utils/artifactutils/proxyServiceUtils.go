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

package artifactutils

type Proxy struct {
	Name    string `json:"name"`
	Wsdl11  string `json:"wsdl1_1"`
	Wsdl20  string `json:"wsdl2_0"`
	Stats   string `json:"stats"`
	Tracing string `json:"tracing"`
}

type ProxyServiceList struct {
	Count   int32          `json:"count"`
	Proxies []ProxySummary `json:"list"`
}

type ProxySummary struct {
	Name   string `json:"name"`
	Wsdl11 string `json:"wsdl1_1"`
	Wsdl20 string `json:"wsdl2_0"`
}
