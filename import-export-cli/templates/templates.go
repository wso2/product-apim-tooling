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
package templates

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"
)

// basicFuncs are used for common data printing
var basicFuncs = template.FuncMap{
	"json": func(v interface{}) string {
		buf := &bytes.Buffer{}
		encoder := json.NewEncoder(buf)
		encoder.SetEscapeHTML(false)
		_ = encoder.Encode(v)
		return strings.TrimSpace(buf.String())
	},
	"jsonPretty": func(v interface{}) string {
		buf := &bytes.Buffer{}
		encoder := json.NewEncoder(buf)
		encoder.SetEscapeHTML(false)
		encoder.SetIndent("", "  ")
		_ = encoder.Encode(v)
		return strings.TrimSpace(buf.String())
	},
	"split": strings.Split,
	"upper": strings.ToUpper,
	"lower": strings.ToLower,
	"title": strings.Title,
	"join":  strings.Join,
}

// HeaderFuncs are used to format headers in a table
// Some of functions in basicFuncs are overridden
var HeaderFuncs = template.FuncMap{
	"json":       func(s string) string { return s },
	"jsonPretty": func(s string) string { return s },
	"join":       func(s string) string { return s },
}

// NewBasicFormatter creates a new template engine with name
func NewBasicFormatter(name string) *template.Template {
	tmpl := template.New(name).Funcs(basicFuncs)
	return tmpl
}
