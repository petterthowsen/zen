package global

import (
	"fmt"
	"zen/runtime"
	"zen/runtime/types"
)

// Print prints the string representations of all parameters to stdout
// returns Zen Bool(true)
func Print(env *runtime.EnvironmentInterface, params map[string]types.Value) (types.Value, error) {
	// convert all parameters to their string representation
	strings := make([]any, 0)

	for _, paramVal := range params {
		strValue := paramVal.String()
		strings = append(strings, strValue)
	}

	_, err := fmt.Println(strings...)
	if err != nil {
		return nil, err
	}
	return types.NewBool(true), err
}
