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

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/impl"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"io"
	"os"
	"path/filepath"
)

var bundleDestination string
var bundleSource string

// Get command related usage Info
const BundleCmdLiteral = "bundle"
const BundleCmdShortDesc = "Archive any project to zip format"

const BundleCmdLongDesc = "Archive API, Application or API Product projects to a zip format. Bundle name will have " +
	"project name, version and revision number (if there is any)"

const BundleCmdExamples = utils.ProjectName + ` ` + BundleCmdLiteral + ` -s /home/prod/APIs/API1-1.0.0 -d /home/prod/Projects/
` + utils.ProjectName + ` ` + BundleCmdLiteral + ` -s /home/prod/APIs/API1-1.0.0 
NOTE: The flag (--source (-s)) is mandatory.`

// BundleCmd represents the bundle command
var BundleCmd = &cobra.Command{
	Use:     BundleCmdLiteral,
	Short:   BundleCmdShortDesc,
	Long:    BundleCmdLongDesc,
	Example: BundleCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(utils.LogPrefixInfo + BundleCmdLiteral + " called")

		if stat, err := os.Stat(bundleSource); !os.IsNotExist(err) {
			if !stat.IsDir() {
				fmt.Printf("%s is not a directory\n", bundleSource)
				os.Exit(1)
			}
		}

		err := executeBundleCmd()
		if err != nil {
			utils.HandleErrorAndContinue("Error archiving the " + bundleSource, err)
		}
	},
}

func executeBundleCmd() error {
	var bundleDirParent string

	// Check the validity of destination path when it is given. if not given, use the working directory
	if bundleDestination != "" {
		err := os.MkdirAll(bundleDestination, os.ModePerm)
		if err != nil {
			return err
		}
		p, err := filepath.Abs(bundleDestination)
		if err != nil {
			return err
		}
		bundleDirParent = p
	} else {
		pwd, err := os.Getwd()
		if err != nil {
			return err
		}
		bundleDirParent = pwd
	}

	bundleName, err := generateBundleName(bundleSource)
	if err != nil {
		return err
	}

	bundleLocation := filepath.Join(bundleDirParent, bundleName + utils.ZipFileSuffix)
	err = utils.Zip(bundleSource, bundleLocation)
	if err != nil {
		return err
	}

	fmt.Println("The bundle for the " + bundleName + " is generated at " + bundleLocation)
	return nil
}

func generateBundleName(SourceDir string) (string, error) {
	metaFileName, err := impl.GetFileLocationFromPattern(SourceDir, "*_meta.yaml")
	if err != nil && err != io.EOF {
		fmt.Println("Error reading the meta information.", err)
	}
	bundleName := filepath.Base(SourceDir)
	if metaFileName != "" {
		metaData, err := impl.LoadMetaInfoFromFile(metaFileName)
		if err != nil {
			return bundleName, err
		}

		bundleName = metaData.Name + "_" + metaData.Version
		if metaData.Revision != "" {
			bundleName += "_" + metaData.Revision
		}
	} else {
		fmt.Println( "Meta information for the Project is not found. Source directory name will be used as the bundle name.")
	}

	return bundleName, nil
}

// init using Cobra
func init() {
	RootCmd.AddCommand(BundleCmd)
	BundleCmd.Flags().StringVarP(&bundleDestination, "destination", "d", "", "Path of "+
		"the directory where the bundle should be generated")
	BundleCmd.Flags().StringVarP(&bundleSource, "source", "s", "", "Path of "+
		"the source directory to bundle")
	_ = BundleCmd.MarkFlagRequired("source")
}
