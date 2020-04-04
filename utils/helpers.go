package utils

import (
	"encoding/json"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

// Bool is a helper routine that allocates a new bool value
// to store v and returns a pointer to it.

// GetBool is a helper routine that returns a boolean representing
// if a value was set, and if so, dereferences the pointer to it.
func GetBool(v *bool) (bool, bool) {
	if v != nil {
		return *v, true
	}

	return false, false
}

// GetInt is a helper routine that returns a boolean representing
// if a value was set, and if so, dereferences the pointer to it.
func GetIntOk(v *int) (int, bool) {
	if v != nil {
		return *v, true
	}

	return 0, false
}

// String is a helper routine that allocates a new string value
// to store v and returns a pointer to it.

// GetString is a helper routine that returns a boolean representing
// if a value was set, and if so, dereferences the pointer to it.
func GetStringOk(v *string) (string, bool) {
	if v != nil {
		return *v, true
	}

	return "", false
}

// JsonNumber is a helper routine that allocates a new string value
// to store v and returns a pointer to it.
func JsonNumber(v json.Number) *json.Number { return &v }

// GetJsonNumber is a helper routine that returns a boolean representing
// if a value was set, and if so, dereferences the pointer to it.
func GetJsonNumberOk(v *json.Number) (json.Number, bool) {
	if v != nil {
		return *v, true
	}

	return "", false
}

func convertToHelperSchema(descrs attrDescrs, in map[schemaAttr]*schema.Schema) map[string]*schema.Schema {
	out := make(map[string]*schema.Schema, len(in))
	for k, v := range in {
		if descr, ok := descrs[k]; ok {
			// NOTE(sean@): At some point this check needs to be uncommented and all
			// empty descriptions need to be populated.
			//
			// if len(descr) == 0 {
			// 	log.Printf("[WARN] PROVIDER BUG: Description of attribute %s empty", k)
			// }

			v.Description = string(descr)
		} else {
			log.Printf("[WARN] PROVIDER BUG: Unable to find description for attr %q", k)
		}

		out[string(k)] = v
	}

	return out
}

func StringPtr(v string) *string {
	return &v
}
