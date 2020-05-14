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

package integration

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
)

func TestVersion(t *testing.T) {
	output, _ := base.Execute(t, "version")

	assert.True(t, strings.Contains(output, "3.1.0"), "Expected version is 3.1.0, actual version is "+output)
}
