package utils

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

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
	certs := x509.NewCertPool()

	certs.AppendCertsFromPEM(WSO2PublicCertificate)

	return &tls.Config{
		InsecureSkipVerify: false,
		RootCAs:            certs,
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
