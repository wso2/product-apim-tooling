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

package formatter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
	"text/template"
)

// MarshalJSON marshals x into json
// But it creates json fields with Title case like Id, Name
func MarshalJSON(x interface{}) ([]byte, error) {
	m, err := marshalMap(x)
	if err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

// marshalMap marshals x to map[string]interface{}
// NOTE: this method only work for plain structs, nested structs are not marshaled correctly
func marshalMap(x interface{}) (map[string]interface{}, error) {
	val := reflect.ValueOf(x)
	if val.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("expected a pointer to a struct, got %v", val.Kind())
	}
	if val.IsNil() {
		return nil, fmt.Errorf("expected a pointer to a struct, got nil pointer")
	}
	values := val.Elem()
	if values.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a pointer to a struct, got a pointer to %v", values.Kind())
	}
	fields := reflect.TypeOf(x)
	num := fields.Elem().NumField()
	m := make(map[string]interface{})
	for i := 0; i < num; i++ {
		fieldName := fields.Elem().Field(i).Name
		fieldValue := values.Field(i).Interface()
		m[fieldName] = fieldValue
	}
	return m, nil
}

// contains helper functions for common printing
var basicFunc = template.FuncMap{
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

// NewBasicFormatter creates a new template engine with name
func NewBasicFormatter(name string) *template.Template {
	tmpl := template.New(name).Funcs(basicFunc)
	return tmpl
}

// Execute template on provided writer
func Execute(w io.Writer, t *template.Template, format string, data interface{}) error {
	apiTmpl, err := t.Parse(format)
	if err != nil {
		return err
	}
	err = apiTmpl.Execute(w, data)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte{'\n'})
	return err
}
