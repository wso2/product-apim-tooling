package utils

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
)

const PlainTextWarnMessage = "WARNING: Error importing the certificate %s\n"

func ReadFromUrl(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	_ = response.Body.Close()
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, errors.New(response.Status)
	}
	return body, nil
}

func GetTlsConfigWithCertificate() *tls.Config {

	certs := ReadCertsFromDir()
	certs.AppendCertsFromPEM(WSO2PublicCertificate)

	return &tls.Config{
		InsecureSkipVerify: false,
		RootCAs:            certs,
		Renegotiation:      TLSRenegotiationMode,
	}
}

// IsValidUrl tests a string to determine if it is a well-structured url or not.
func IsValidUrl(urlStr string) bool {
	_, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return false
	}

	u, err := url.Parse(urlStr)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

func ReadCertsFromDir() *x509.CertPool {
	certs, err := x509.SystemCertPool()
	if err != nil || certs == nil {
		//if the OS is windows, systemCertPool will return an error. For windows, CA certificates has to be added
		//to the .wso2apictl/certs directory.
		certs = x509.NewCertPool()
	}

	certificates, err := ioutil.ReadDir(DefaultCertDirPath)
	if err == nil {
		for _, certificate := range certificates {
			extension := filepath.Ext(certificate.Name())
			if extension == ".pem" || extension == ".crt" || extension == ".cer" {
				certFilePath := filepath.Join(DefaultCertDirPath, certificate.Name())
				fileData, err := ioutil.ReadFile(certFilePath)
				if fileData != nil && err == nil {
					if c, err := x509.ParseCertificate(fileData); err == nil {
						//if the certificate is DER encoded, add it directly to the cert pool.
						certs.AddCert(c)
					} else {
						//if the certificate is PEM encoded.
						certs.AppendCertsFromPEM(fileData)
					}
				} else {
					fmt.Printf(PlainTextWarnMessage, certificate.Name())
				}
			}
		}
	}
	return certs
}
