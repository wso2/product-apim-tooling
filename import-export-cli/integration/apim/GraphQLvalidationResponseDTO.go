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
package apim

// GraphQLValidationResponseDTO : GraphQL Validation Response DTO
type GraphQLValidationResponseDTO struct {
	IsValid      bool                                    `json:"isValid,omitempty" yaml:"isValid,omitempty"`
	ErrorMessage string                                  `json:"errorMessage,omitempty" yaml:"errorMessage,omitempty"`
	GraphQLInfo  GraphQLValidationResponseGraphQLInfoDTO `json:"graphQLInfo,omitempty" yaml:"graphQLInfo,omitempty"`
}

// GraphQLValidationResponseGraphQLInfoDTO: GraphQL Validation Response GraphQL Info DTO
type GraphQLValidationResponseGraphQLInfoDTO struct {
	Operations    []interface{} `json:"operations,omitempty" yaml:"operations,omitempty"`
	GraphQLSchema interface{}   `json:"graphQLSchema,omitempty" yaml:"graphQLSchema,omitempty"`
}
