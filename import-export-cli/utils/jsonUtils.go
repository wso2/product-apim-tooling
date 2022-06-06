package utils

import (
	"encoding/json"
	"fmt"
	"github.com/Jeffail/gabs"
	"github.com/savaki/jq"
)

// MergeJSON secondSource with firstSource and returns merged JSON string
// Note: Fields in firstSource are merged with secondSource.
// If a field is not presented in secondSource, the one in firstSource will be preserved.
// If not a field from secondSource will replace it.
func MergeJSON(firstSource, secondSource []byte) ([]byte, error) {
	secondSourceJSON, err := gabs.ParseJSON(secondSource)
	if err != nil {
		return nil, err
	}

	firstSourceJSON, err := gabs.ParseJSON(firstSource)
	if err != nil {
		return nil, err
	}

	err = firstSourceJSON.MergeFn(secondSourceJSON, func(destination, source interface{}) interface{} {
		if source == nil {
			return destination
		}
		if s, ok := source.(string); ok && s == "" {
			return destination
		}
		return source
	})

	return firstSourceJSON.Bytes(), nil
}

// ListArtifactsInJsonArrayFormat : This function will return the output of list apis/apiProducts/apps command in
// JsonObject format
func ListArtifactsInJsonArrayFormat(artifacts interface{}, artifactType string) {

	data, err := json.Marshal(artifacts)
	if err != nil {
		fmt.Println("Error executing template:", err.Error())
	}

	op, err := jq.Parse("")
	if err != nil {
		fmt.Println("Error executing template:", err.Error())
	}

	formattedData, _ := selectTypeOfOutputEntry(data, artifactType)
	output, err := op.Apply(formattedData)
	if err != nil {
		fmt.Println("Error executing template:", err.Error())
	}

	// Return JsonArray format output to CLI
	fmt.Println(string(output))
}

// Get formatted output based on the type of artifact
func selectTypeOfOutputEntry(data []byte, artifactType string) ([]byte, error) {

	if artifactType == ProjectTypeApi {
		var apiEntries []APIEntry
		// Map API information to APIEntry struct
		json.Unmarshal(data, &apiEntries)

		// Formatting data to get the JsonArray object in prettyPrint format
		return json.MarshalIndent(apiEntries, "", " ")
	} else if artifactType == ProjectTypeApiProduct {
		var apiProductEntries []APIProductEntry
		// Map API information to APIEntry struct
		json.Unmarshal(data, &apiProductEntries)

		// Formatting data to get the JsonArray object in prettyPrint format
		return json.MarshalIndent(apiProductEntries, "", " ")
	} else if artifactType == ProjectTypeApplication {
		var applicationEntries []ApplicationEntry
		// Map API information to APIEntry struct
		json.Unmarshal(data, &applicationEntries)

		// Formatting data to get the JsonArray object in prettyPrint format
		return json.MarshalIndent(applicationEntries, "", " ")
	} else if artifactType == ProjectTypePolicy {
		var policyEntries []PolicyEntry
		json.Unmarshal(data, &policyEntries)

		// Formatting data to get the JsonArray object in prettyPrint format
		return json.MarshalIndent(policyEntries, "", " ")
	} else {
		var revisionEntries []RevisionEntry
		// Map API information to APIEntry struct
		json.Unmarshal(data, &revisionEntries)

		// Formatting data to get the JsonArray object in prettyPrint format
		return json.MarshalIndent(revisionEntries, "", " ")
	}
}
