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

package base

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	// RelativeBinaryPath : Relative apictl binary path used by integration tests
	RelativeBinaryPath string

	// BinaryName : Name of apictl binary
	BinaryName = "apictl"

	archiveFileName string
)

func init() {
	flag.StringVar(&archiveFileName, "archive", "", "Archive file name of apictl distribution")
}

// ExtractArchiveFile : Extract apictl distribution archive file
func ExtractArchiveFile(path string) {
	if archiveFileName == "" {
		Fatal("apictl archive not provided as a command line argument '-archive <archive_file_name>'")
	}

	relativePath := path
	extractedFolder := "extracted"

	destPath := filepath.FromSlash(relativePath + extractedFolder)
	srcPath := filepath.FromSlash(relativePath + archiveFileName)

	os.RemoveAll(destPath)

	var err error
	if strings.Contains(archiveFileName, ".zip") {
		err = Unzip(destPath, srcPath)
		BinaryName = "apictl.exe" // Windows binaries are archived using zip
	} else { // tar.gz
		err = untar(destPath, srcPath)
	}

	if err != nil {
		Fatal(err)
	}

	subFolder := "/apictl/"
	RelativeBinaryPath = filepath.FromSlash(destPath + subFolder)
}

func untar(dstPath string, srcPath string) error {
	reader, err := os.Open(srcPath)
	if err != nil {
		Fatal(err)
	}
	defer reader.Close()

	gzr, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dstPath, header.Name)

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			f.Close()
		}
	}
}

func Unzip(destPath string, srcPath string) error {
	reader, err := zip.OpenReader(srcPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, f := range reader.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(destPath, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(destPath)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", fpath)
		}

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}
