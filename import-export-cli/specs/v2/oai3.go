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

package v2

import (
	"encoding/json"

	"github.com/getkin/kin-openapi/openapi3"
)

func oai3XWSO2Cors(exts map[string]interface{}) (*CorsConfiguration, bool, error) {
	if v, ok := exts["x-wso2-cors"]; ok {
		ep, ok := v.(json.RawMessage)
		if ok {
			var cors CorsConfiguration
			err := json.Unmarshal(ep, &cors)
			if err != nil {
				return nil, true, err
			}
			cors.CorsConfigurationEnabled = true
			return &cors, true, nil
		}
	}
	return nil, true, nil
}

type Tag struct {
	Name string `json:"name"`
}

func oai3Tags(exts map[string]interface{}) []string {
	if v, ok := exts["tags"]; ok {
		b, ok := v.(json.RawMessage)
		if ok {
			var tags []Tag
			_ = json.Unmarshal(b, &tags)
			strs := make([]string, len(tags))
			for i, t := range tags {
				strs[i] = t.Name
			}
			return strs
		}
	}
	return nil
}

type Endpoints struct {
	Type                  string                        `yaml:"type"`
	Urls                  []string                      `yaml:"urls"`
	AdvanceEndpointConfig *AdvanceEndpointConfiguration `yaml:"advanceEndpointConfig,omitempty"`
}

type AdvanceEndpointConfiguration struct {
	TimeOutInMillis *int `yaml:"timeoutInMillis" json:"timeoutInMillis"`
}

func oai3XWSO2ProductionEndpoints(exts map[string]interface{}) (*Endpoints, bool, error) {
	if v, ok := exts["x-wso2-production-endpoints"]; ok {
		ep, ok := v.(json.RawMessage)
		if ok {
			var prodEp Endpoints
			err := json.Unmarshal(ep, &prodEp)
			if err != nil {
				return nil, true, err
			}
			return &prodEp, true, nil
		}
	}
	return nil, true, nil
}

func oai3XWso2SandboxEndpoints(exts map[string]interface{}) (*Endpoints, bool, error) {
	if v, ok := exts["x-wso2-sandbox-endpoints"]; ok {
		ep, ok := v.(json.RawMessage)
		if ok {
			var sandboxEp Endpoints
			err := json.Unmarshal(ep, &sandboxEp)
			if err != nil {
				return nil, true, err
			}
			return &sandboxEp, true, nil
		}
	}
	return nil, true, nil
}

func oai3XWSO2AuthHeader(exts map[string]interface{}) (string, bool) {
	if v, ok := exts["x-wso2-auth-header"]; ok {
		data, ok := v.(json.RawMessage)
		if ok {
			var auth string
			err := json.Unmarshal(data, &auth)
			if err != nil {
				return "", false
			}
			return auth, true
		}
	}
	return "", false
}

func oai3WSO2Basepath(exts map[string]interface{}) (string, bool, error) {
	if v, ok := exts["x-wso2-basePath"]; ok {
		data, ok := v.(json.RawMessage)
		if ok {
			var basepath string
			err := json.Unmarshal(data, &basepath)
			if err != nil {
				return "", false, err
			}
			return basepath, true, nil
		}
	}
	return "", false, nil
}

// oai3GetHttpVerbs generates verbs for api definition
func oai3GetHttpVerbs(item *openapi3.PathItem) (verbs []string) {
	if item.Get != nil {
		verbs = append(verbs, "GET")
	}
	if item.Post != nil {
		verbs = append(verbs, "POST")
	}
	if item.Put != nil {
		verbs = append(verbs, "PUT")
	}
	if item.Delete != nil {
		verbs = append(verbs, "DELETE")
	}
	if item.Patch != nil {
		verbs = append(verbs, "PATCH")
	}
	if item.Head != nil {
		verbs = append(verbs, "HEAD")
	}
	if item.Options != nil {
		verbs = append(verbs, "OPTIONS")
	}
	return
}
