package url

import (
	"net/url"
	"reflect"
)

// BuildQuery converts a struct to url.Values based on struct field tags.
// It checks for "json", "form", and "url" tags in that order.
func BuildQuery(s interface{}) url.Values {
	values := url.Values{}
	v := reflect.ValueOf(s)

	// Ensure we are working with a struct
	if v.Kind() != reflect.Struct {
		return values
	}

	t := v.Type()

	// Iterate over the struct fields
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Check for tags in the order of json, form, url
		jsonTag := fieldType.Tag.Get("json")
		formTag := fieldType.Tag.Get("form")
		urlTag := fieldType.Tag.Get("url")

		var key string
		if jsonTag != "" {
			key = jsonTag
		} else if formTag != "" {
			key = formTag
		} else if urlTag != "" {
			key = urlTag
		}

		// Only add to values if a key was found and the field is exported
		if key != "" && field.CanInterface() {
			// Convert the field value to a string
			value := reflect.ValueOf(field.Interface())
			values.Add(key, value.String())
		}
	}

	return values
}
