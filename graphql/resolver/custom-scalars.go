package resolver

import "fmt"

type JSON map[string]interface{}

func (JSON) ImplementsGraphQLType(name string) bool {
	return name == "JSON"
}
func (j *JSON) UnmarshalGraphQL(input interface{}) error {
	switch input := input.(type) {
	case map[string]interface{}:
		*j = JSON(input)
		return nil
	default:
		return fmt.Errorf("wrong type")
	}
}
