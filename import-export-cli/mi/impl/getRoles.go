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
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/wso2/product-apim-tooling/import-export-cli/formatter"
	"github.com/wso2/product-apim-tooling/import-export-cli/mi/utils/artifactutils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const (
	defaultRoleListTableFormat = "table {{.Role}}"
	defaultRoleDetailedFormat  = "Role Name - {{.Role}}\n" +
		"Users  - " +
		"{{range $index, $user := .Users}}" +
		"{{if $index}}, {{end}}" +
		"{{$user}}" +
		"{{end}}"
)

// GetRoleList returns a list of roles in the micro integrator in a given environment
func GetRoleList(env string) (*artifactutils.RoleList, error) {
	resp, err := callMIManagementEndpointOfResource(utils.MiManagementRoleResource, nil, env, &artifactutils.RoleList{})
	if err != nil {
		return nil, err
	}
	return resp.(*artifactutils.RoleList), nil
}

// PrintRoleList print a list of mi roles according to the given format
func PrintRoleList(roleList *artifactutils.RoleList, format string) {
	if roleList.Count > 0 {
		roles := roleList.Roles
		roleListContext := getContextWithFormat(format, defaultRoleListTableFormat)

		renderer := func(w io.Writer, t *template.Template) error {
			for _, role := range roles {
				if err := t.Execute(w, role); err != nil {
					return err
				}
				_, _ = w.Write([]byte{'\n'})
			}
			return nil
		}
		roleListTableHeaders := map[string]string{
			"Role": roleHeader,
		}
		if err := roleListContext.Write(renderer, roleListTableHeaders); err != nil {
			fmt.Println("Error executing template:", err.Error())
		}
	} else {
		fmt.Println("No roles found")
	}
}

// GetRoleInfo returns a information about a specific role in the micro integrator in a given environment
func GetRoleInfo(env, role, domain string) (*artifactutils.RoleSummary, error) {
	var roleInfoResource = utils.MiManagementRoleResource + "/" + role
	params := make(map[string]string)
	putNonEmptyValueToMap(params, "domain", domain)
	resp, err := callMIManagementEndpointOfResource(roleInfoResource, params, env, &artifactutils.RoleSummary{})
	if err != nil {
		return nil, err
	}
	return resp.(*artifactutils.RoleSummary), nil
}

// PrintRoleDetails prints details about a role according to the given format
func PrintRoleDetails(roleInfo *artifactutils.RoleSummary, format string) {
	if format == "" || strings.HasPrefix(format, formatter.TableFormatKey) {
		format = defaultRoleDetailedFormat
	}

	roleInfoContext := formatter.NewContext(os.Stdout, format)
	renderer := getItemRendererEndsWithNewLine(roleInfo)

	if err := roleInfoContext.Write(renderer, nil); err != nil {
		fmt.Println("Error executing template:", err.Error())
	}
}
