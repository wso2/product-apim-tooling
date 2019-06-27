package v2

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"testing"
)

var petstoreProdUrls = []string{"https://petstore.swagger.io/v2", "https://petstore.swagger.io/v2/1", "https://petstore.swagger.io/v2/2"}

func Test_oai3WSO2Basepath(t *testing.T) {
	sw, err := openapi3.NewSwaggerLoader().LoadSwaggerFromFile("testdata/petstore_basic.yaml")
	assert.Nil(t, err, "err should be nil")
	basepath, ok, err := oai3WSO2Basepath(sw.Extensions)
	assert.Nil(t, err, "err should be nil")
	assert.True(t, ok, "should have vendor extension")
	assert.Equal(t, "/petstore/v1", basepath, "should return correct basepath")
}

func Test_oai3WSO2ProductionEndpoints(t *testing.T) {
	sw, err := openapi3.NewSwaggerLoader().LoadSwaggerFromFile("testdata/petstore_basic.yaml")
	assert.Nil(t, err, "err should be nil")
	ep, ok, err := oai3XWSO2ProductionEndpoints(sw.Extensions)
	assert.Nil(t, err, "err should be nil")
	assert.True(t, ok, "should have vendor extension")
	assert.ElementsMatch(t, petstoreProdUrls, ep.Urls, "should have same elements")
}

func Test_oai3WSO2SandboxEndpoints(t *testing.T) {
	sw, err := openapi3.NewSwaggerLoader().LoadSwaggerFromFile("testdata/petstore_basic.yaml")
	assert.Nil(t, err, "err should be nil")
	ep, ok, err := oai3XWso2SandboxEndpoints(sw.Extensions)
	assert.Nil(t, err, "err should be nil")
	assert.True(t, ok, "should have vendor extension")
	assert.ElementsMatch(t, petstoreProdUrls, ep.Urls, "should have same elements")
}

func Test_oai3Tags(t *testing.T) {
	sw, err := openapi3.NewSwaggerLoader().LoadSwaggerFromFile("testdata/petstore_basic.yaml")
	assert.Nil(t, err, "err should be nil")
	tags := oai3Tags(sw.Extensions)
	assert.Nil(t, err, "err should be nil")
	assert.ElementsMatch(t, []string{"pet", "user", "store"}, tags, "should have same elements")
}

func Test_oai3WSO2Cors(t *testing.T) {
	sw, err := openapi3.NewSwaggerLoader().LoadSwaggerFromFile("testdata/petstore_basic.yaml")
	assert.Nil(t, err, "err should be nil")
	cors, ok, err := oai3XWSO2Cors(sw.Extensions)
	assert.Nil(t, err, "err should be nil")
	assert.True(t, ok, "should have vendor extension")
	assert.ElementsMatch(t, []string{"GET", "PUT", "POST"}, cors.AccessControlAllowMethods, "should have same elements for access control")
	assert.ElementsMatch(t, []string{"test.com", "example.com"}, cors.AccessControlAllowOrigins, "should have same elements for origins")
}
