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

 package aws

 import (
	 "github.com/spf13/cobra"
 )
 
 const awsCmdShortDesc = "AWS Api-gateway related commands"
 const awsCmdLongDesc = `AWS Api-gateway related commands such as init.`
 const awsCmdLiteral = "aws"
 
 // MICmd represents the mi command
 var AWSCmd = &cobra.Command{
	 Use:   awsCmdLiteral,
	 Short: awsCmdShortDesc,
	 Long:  awsCmdLongDesc,
	 Example: awsInitCmdExamples,
	 Run: func(cmd *cobra.Command, args []string) {
		 cmd.Help()
	 },
 }
 
 func init() {
	 AWSCmd.AddCommand(InitCmd)
 }
 