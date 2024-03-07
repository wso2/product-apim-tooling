/*
 *  Copyright (c) 2024, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

package auth

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBasicAuth(t *testing.T) {
	authData := map[string]string{
		"authUser":      "Is7ZOcq7EbYCf13",
		"randomUser":    "AM8MprkIukDEpgo",
		"TestDummyUser": "X8vDK4pgkxnUXmM",
	}
	for _, authPair := range authData {
		res := GetBasicAuth(authPair, authData[authPair])
		assert.IsType(t, string(""), res)
		assert.Equal(t, base64.StdEncoding.EncodeToString([]byte(authPair+":"+authData[authPair])), res)
		assert.NotEqual(t, base64.StdEncoding.EncodeToString([]byte(authPair+authData[authPair])), res)
	}
}
