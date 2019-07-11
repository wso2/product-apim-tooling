package utils

import "github.com/Jeffail/gabs"

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
