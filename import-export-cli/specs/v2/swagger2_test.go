package v2

import (
	"testing"

	"github.com/go-openapi/loads"
	"github.com/stretchr/testify/assert"
)

func Test_swagger2WSO2Cors(t *testing.T) {
	doc, err := loads.Spec("testdata/petstore_swagger2.yaml")
	assert.Nil(t, err, "err should be nil")
	cors, ok, err := swagger2XWSO2Cors(doc)
	assert.Nil(t, err, "err should be nil")
	assert.True(t, ok, "should have vendor extension")
	assert.ElementsMatch(t, []string{"GET", "PUT", "POST"}, cors.AccessControlAllowMethods, "should have same elements for access control")
	assert.ElementsMatch(t, []string{"test.com", "example.com"}, cors.AccessControlAllowOrigins, "should have same elements for origins")
}

func Test_swagger2Tags(t *testing.T) {
	doc, err := loads.Spec("testdata/petstore_swagger2.yaml")
	assert.Nil(t, err, "err should be nil")
	tags := swagger2Tags(doc)
	assert.ElementsMatch(t, []string{"pet", "user", "store"}, tags, "should have same elements")
}

func Test_swagger2WSO2ProductionEndpoints(t *testing.T) {
	doc, err := loads.Spec("testdata/petstore_swagger2.yaml")
	assert.Nil(t, err, "err should be nil")
	ep, ok, err := swagger2XWSO2ProductionEndpoints(doc)
	assert.Nil(t, err, "err should be nil")
	assert.True(t, ok, "should have vendor extension")
	assert.ElementsMatch(t, petstoreProdUrls, ep.Urls, "should have same elements")
}

func Test_swagger2WSO2SandboxEndpoints(t *testing.T) {
	doc, err := loads.Spec("testdata/petstore_swagger2.yaml")
	assert.Nil(t, err, "err should be nil")
	ep, ok, err := swagger2XWSO2SandboxEndpoints(doc)
	assert.Nil(t, err, "err should be nil")
	assert.True(t, ok, "should have vendor extension")
	assert.ElementsMatch(t, petstoreProdUrls, ep.Urls, "should have same elements")
}

func Test_swagger2HTTPVerbs(t *testing.T) {
	doc, err := loads.Spec("testdata/petstore_swagger2.yaml")
	assert.Nil(t, err, "err should be nil")
	item := doc.Spec().Paths.Paths["/pet/findByStatus"]
	verbs := swagger2GetHttpVerbs(item)
	assert.ElementsMatch(t, []string{"GET"}, verbs, "Should return correct values")
}

func TestSwagger2Populate(t *testing.T) {
	var def APIDefinition
	doc, err := loads.Spec("testdata/petstore_swagger2.yaml")
	assert.Nil(t, err, "err should be nil")
	err = Swagger2Populate(&def, doc)
	assert.Nil(t, err, "err should be nil")

	assert.Equal(t, "SwaggerPetstore", def.ID.APIName, "Should return correct api name")
	assert.Equal(t, "/petstore/v1/1.0.0", def.Context)
}
