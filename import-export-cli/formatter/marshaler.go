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
	"encoding/json"
	"fmt"
	"reflect"
	"unicode"
)

// MarshalJSON marshals x into json
// It differs a bit from encoding/json MarshalJSON function for formatter
// This method will use Methods instead of fields.
// Make sure that method are Exported
func MarshalJSON(x interface{}) ([]byte, error) {
	m, err := marshalMap(x)
	if err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

// marshalMap marshals x to map[string]interface{}
func marshalMap(x interface{}) (map[string]interface{}, error) {
	val := reflect.ValueOf(x)
	if val.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("expected a pointer to a struct, got %v", val.Kind())
	}
	if val.IsNil() {
		return nil, fmt.Errorf("expected a pointer to a struct, got nil pointer")
	}
	valElem := val.Elem()
	if valElem.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a pointer to a struct, got a pointer to %v", valElem.Kind())
	}
	typ := val.Type()
	m := make(map[string]interface{})
	for i := 0; i < val.NumMethod(); i++ {
		k, v, err := marshalForMethod(typ.Method(i), val.Method(i))
		if err != nil {
			return nil, err
		}
		if k != "" {
			m[k] = v
		}
	}
	return m, nil
}

// marshalForMethod returns the map key and the map value for marshalling the method.
// It returns ("", nil, nil) for valid but non-marshallable parameter. (e.g. "unexportedFunc()")
// This only works for methods without parameters and outputs only single value
func marshalForMethod(typ reflect.Method, val reflect.Value) (string, interface{}, error) {
	if val.Kind() != reflect.Func {
		return "", nil, fmt.Errorf("expected func, got %v", val.Kind())
	}
	name, numIn, numOut := typ.Name, val.Type().NumIn(), val.Type().NumOut()
	marshallable := unicode.IsUpper(rune(name[0])) && numIn == 0 && numOut == 1
	if !marshallable {
		return "", nil, nil
	}
	result := val.Call(make([]reflect.Value, numIn))
	intf := result[0].Interface()
	return name, intf, nil
}
