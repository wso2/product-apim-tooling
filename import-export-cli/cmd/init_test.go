package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_createDirectories(t *testing.T) {
	name, err := ioutil.TempDir("", "")
	assert.Nil(t, err, "Temp directory should be created")
	err = os.Chdir(name)
	assert.Nil(t, err, "Should be able to change directory")
	err = createDirectories()
	assert.Nil(t, err, "Should be no errors when creating directory structure")

	for _, dir := range dirs {
		dirPath := filepath.FromSlash(dir)
		assert.DirExists(t, dirPath, "Directory "+dirPath+" should be created")
	}

	_ = os.Chdir("..")
	_ = os.RemoveAll(name)
}

func Test_getDefaultTiers(t *testing.T) {
	tiers := getDefaultTiers()
	assert.Equal(t, 4, len(tiers), "Should load four default tiers")
	assert.Equal(t, "Bronze", tiers[0].Name, "Should get correct name")
	assert.Equal(t, "Allows 1000 requests per minute", tiers[0].Description, "Should get correct description")
}

func Test_getDefaultCORS(t *testing.T) {
	cors := getDefaultCORS()
	assert.Equal(t, false, cors.CorsConfigurationEnabled, "Should load default")
	assert.Equal(t, 6, len(cors.AccessControlAllowMethods), "Should load correct values")
}

func Test_loadSwagger2(t *testing.T) {
	sw, _, err := loadSwagger("testdata/swaggers/swagger-2.json")
	assert.Nil(t, err, "Loads correct swagger without errors")
	assert.Equal(t, "Simple API overview", sw.Info.Title, "Loads correct title")
	assert.NotNil(t, sw.Paths, "Paths should not be nil")
}

func Test_loadSwagger2YAML(t *testing.T) {
	sw, _, err := loadSwagger("testdata/swaggers/swagger-2.yaml")
	assert.Nil(t, err, "Loads correct swagger without errors")
	assert.Equal(t, "Simple API overview", sw.Info.Title, "Loads correct title")
	assert.NotNil(t, sw.Paths, "Paths should not be nil")
}

func Test_loadSwagger3(t *testing.T) {
	sw, _, err := loadSwagger("testdata/swaggers/swagger-3.json")
	assert.Nil(t, err, "Loads correct swagger without errors")
	assert.Equal(t, "Swagger Petstore", sw.Info.Title, "Loads correct title")
	assert.NotNil(t, sw.Paths, "Paths should not be nil")
}

func Test_loadSwagger3YAML(t *testing.T) {
	sw, _, err := loadSwagger("testdata/swaggers/swagger-3.yaml")
	assert.Nil(t, err, "Loads correct swagger without errors")
	assert.Equal(t, "Swagger Petstore", sw.Info.Title, "Loads correct title")
	assert.NotNil(t, sw.Paths, "Paths should not be nil")
}

func Test_APIDefinition_generateFieldsFromSwagger(t *testing.T) {
	sw, _, err := loadSwagger("testdata/swaggers/swagger-3.json")
	assert.Nil(t, err, "Loads correct swagger without errors")
	def := newApiDefinitionWithDefaults()
	def.generateFieldsFromSwagger(sw)

	assert.Equal(t, "SwaggerPetstore", def.ID.APIName, "Should correctly output name")
	assert.Equal(t, "/SwaggerPetstore/1.0.0", def.Context, "Should return correct context")
	assert.Equal(t, 14, len(def.URITemplates), "Should return correct number of uri templates")
}
