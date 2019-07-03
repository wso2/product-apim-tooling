package utils

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
	yaml2 "gopkg.in/yaml.v2"
)

// JsonToYaml converts a json string to yaml
func JsonToYaml(jsonData []byte) ([]byte, error) {
	var m yaml2.MapSlice
	err := yaml2.Unmarshal(jsonData, &m)
	if err != nil {
		return nil, err
	}

	data, err := yaml2.Marshal(m)
	if err != nil {
		return nil, err
	}

	return data, nil
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
