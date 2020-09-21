package utils

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
)

// JsonToYaml converts a json string to yaml
func JsonToYaml(jsonData []byte) ([]byte, error) {
	return yaml.JSONToYAML(jsonData)
}

// YamlToJson converts a yaml string to json
func YamlToJson(yamlData []byte) ([]byte, error) {
	return yaml.YAMLToJSON(yamlData)
}

// LoadYamlAsJson is acting as a wrapper for load a yaml file in json
func LoadYamlAsJson(fp string) ([]byte, error) {
	yamlData, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil, err
	}
	jsonData, err := YamlToJson(yamlData)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}
