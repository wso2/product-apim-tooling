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

package notifier

import (
	"reflect"
	"testing"
)

func TestUpdateDeployedRevisions(t *testing.T) {
	testCases := []struct {
		apiID      string
		revisionID int
		envs       []string
		vhost      string
		expected   *DeployedAPIRevision
	}{
		{
			apiID:      "api1",
			revisionID: 1,
			envs:       []string{"dev", "prod"},
			vhost:      "example.wso2.com",
			expected: &DeployedAPIRevision{
				APIID:      "api1",
				RevisionID: 1,
				EnvInfo: []DeployedEnvInfo{
					{Name: "dev", VHost: "example.wso2.com"},
					{Name: "prod", VHost: "example.wso2.com"},
				},
			},
		},
		{
			apiID:      "wkd1w18uh2o0e22oi7bnc29ue902jd8ud38d",
			revisionID: 2,
			envs:       []string{},
			vhost:      "example.dev.org",
			expected: &DeployedAPIRevision{
				APIID:      "wkd1w18uh2o0e22oi7bnc29ue902jd8ud38d",
				RevisionID: 2,
				EnvInfo:    []DeployedEnvInfo{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.apiID, func(t *testing.T) {
			result := UpdateDeployedRevisions(tc.apiID, tc.revisionID, tc.envs, tc.vhost)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Unexpected result. Expected: %v, Got: %v", tc.expected, result)
			}
		})
	}
}
