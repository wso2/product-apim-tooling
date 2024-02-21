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

package config

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/wso2/apk/adapter/pkg/utils/stringutils"
)

// WireLogConfig is the config holder for wire logs
type WireLogConfig struct {
	LogBodyEnabled     bool
	LogHeadersEnabled  bool
	LogTrailersEnabled bool
}

var (
	onceTrustedCertsRead sync.Once
	onceKeyCertsRead     sync.Once
	certificate          tls.Certificate
	certReadErr          error
	caCertPool           *x509.CertPool
)

const (
	pemExtension string = ".pem"
	crtExtension string = ".crt"
)

// GetWireLogConfig converts the wire log configs read from the config file
func GetWireLogConfig() *WireLogConfig {
	logConf := ReadLogConfigs()
	wireLogsEnabled := logConf.WireLogs.Enable
	includes := logConf.WireLogs.Include
	return &WireLogConfig{
		LogBodyEnabled:     wireLogsEnabled && stringutils.StringInSlice("Body", includes),
		LogHeadersEnabled:  wireLogsEnabled && stringutils.StringInSlice("Headers", includes),
		LogTrailersEnabled: wireLogsEnabled && stringutils.StringInSlice("Trailers", includes),
	}
}

// GetKeyLocations function returns the public key path and private key path
func GetKeyLocations() (string, string, string) {
	conf, _ := ReadConfigs()
	publicKeyLocation := conf.Agent.Keystore.CertPath
	privateKeyLocation := conf.Agent.Keystore.KeyPath
	truststoreLocation := conf.Agent.TrustStore.Location
	return publicKeyLocation, privateKeyLocation, truststoreLocation
}

// GetServerCertificate returns the certificate (used for the restAPI server and GRPC server) created based on configuration values.
// Move to pkg. remove config and read from a file path
func GetServerCertificate(tlsCertificate string, tlsCertificateKey string) (tls.Certificate, error) {
	certReadErr = nil
	onceKeyCertsRead.Do(func() {
		cert, err := tls.LoadX509KeyPair(string(tlsCertificate), string(tlsCertificateKey))
		if err != nil {
			loggerConfig.Logger.Fatal("Error while loading the tls keypair.", err)
			certReadErr = err
		}
		certificate = cert
	})
	return certificate, certReadErr
}

// GetTrustedCertPool returns the trusted certificate (used for the restAPI server and xds server) created based on
// the provided directory/file path.
// Move to pkg
func GetTrustedCertPool(truststoreLocation string) *x509.CertPool {
	onceTrustedCertsRead.Do(func() {
		caCertPool = x509.NewCertPool()
		err := filepath.Walk(truststoreLocation, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				loggerConfig.Logger.Warn("Error while reading the trusted certificates directory/file.", err)
			} else {
				if !info.IsDir() && (filepath.Ext(info.Name()) == pemExtension ||
					filepath.Ext(info.Name()) == crtExtension) {
					caCert, caCertErr := ioutil.ReadFile(filepath.Clean(path))
					if caCertErr != nil {
						loggerConfig.Logger.Warn("Error while reading the certificate file.", info.Name())
					}
					if IsPublicCertificate(caCert) {
						caCertPool.AppendCertsFromPEM(caCert)
						loggerConfig.Logger.Debugf("%v : Certificate is added as a trusted certificate.", info.Name())
					}
				}
			}
			return nil
		})
		if err != nil {
			loggerConfig.Logger.Warnf("Error walking the path %s. Error: %s", truststoreLocation, err)
		}
	})
	return caCertPool
}

// IsPublicCertificate checks if the file content represents valid public certificate in PEM format.
// Move to pkg
func IsPublicCertificate(certContent []byte) bool {
	certContentPattern := `\-\-\-\-\-BEGIN\sCERTIFICATE\-\-\-\-\-((.|\n)*)\-\-\-\-\-END\sCERTIFICATE\-\-\-\-\-`
	regex := regexp.MustCompile(certContentPattern)
	if regex.Match(certContent) {
		return true
	}
	return false
}
