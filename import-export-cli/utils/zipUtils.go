/*
*  Copyright (c) 2005-2017, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
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

package utils

import (
	"github.com/jhoonb/archivex"
	"os"
)

func ZipDir(source, target string) error {
	err := os.Chdir(source)
	if err == nil {
		Logln(LogPrefixInfo + "Directory " + source + " exists")
		Logln(LogPrefixInfo + "Starting Compression...")
		zip := new(archivex.ZipFile)
		zip.Create(target)
		zip.AddAll(source, true)
		zip.Close()
		Logln(LogPrefixInfo + "Compression completed: Find file " + target)
		return nil
	} else {
		Logln(LogPrefixError + "Compressing " + source)
		return err
	}
}
