package llm

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/invopop/jsonschema"
)

// GenerateSchema produces a JSON schema map from a Go struct using
// invopop/jsonschema. OpenAI structured outputs require all properties
// in "required" and use nullable types for optional fields.
func GenerateSchema(v any) map[string]any {
	r := new(jsonschema.Reflector)
	r.DoNotReference = true
	s := r.Reflect(v)

	b, _ := json.Marshal(s)
	var schema map[string]any
	_ = json.Unmarshal(b, &schema)

	delete(schema, "$schema")
	delete(schema, "$id")

	MakeNullable(schema, reflect.TypeOf(v))

	return schema
}

// MakeNullable walks a JSON schema and its corresponding Go type, converting
// pointer fields from type:"string" to type:["string","null"] for OpenAI
// structured output compatibility.
func MakeNullable(schema map[string]any, t reflect.Type) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return
	}

	props, _ := schema["properties"].(map[string]any)
	if props == nil {
		return
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		jsonTag := f.Tag.Get("json")
		name := strings.Split(jsonTag, ",")[0]
		if name == "" || name == "-" {
			continue
		}

		prop, ok := props[name].(map[string]any)
		if !ok {
			continue
		}

		if f.Type.Kind() == reflect.Ptr {
			if typeStr, ok := prop["type"].(string); ok {
				prop["type"] = []any{typeStr, "null"}
			}
		}

		// Recurse into array items
		if f.Type.Kind() == reflect.Slice {
			elemType := f.Type.Elem()
			if elemType.Kind() == reflect.Ptr {
				elemType = elemType.Elem()
			}
			if items, ok := prop["items"].(map[string]any); ok && elemType.Kind() == reflect.Struct {
				MakeNullable(items, elemType)
			}
		}

		// Recurse into nested structs
		if f.Type.Kind() == reflect.Struct {
			MakeNullable(prop, f.Type)
		}
	}
}
