package metrics

import (
	"encoding/json"
	"regexp"
	"strings"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func Extract(message []byte) map[string]any {
	incoming := map[string]any{}
	normalized := map[string]any{}
	err := json.Unmarshal(message, &incoming)
	if err == nil {
		normalized = normalize(incoming)
	}
	return normalized
}

func normalize(incoming map[string]any) map[string]any {
	normalized := map[string]any{}
	for k, v := range incoming {
		switch x := v.(type) {
		default:
			//ignore
		case int64:
			normalized[toSnakeCase(k)] = float64(x)
		case float64:
			normalized[toSnakeCase(k)] = x
		case string:
			normalized[toSnakeCase(k)] = x
		case map[string]any:
			converted := normalize(x)
			for k1, v := range converted {
				normalized[toSnakeCase(k)+"_"+k1] = v
			}
		}
	}
	return normalized
}

func toSnakeCase(camelCase string) string {
	snake := matchFirstCap.ReplaceAllString(camelCase, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
