package util

import (
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

var fieldNamePattern = regexp.MustCompile("field ([^ ]+)")

// YamlUnmarshalStrictIgnoringFields does UnmarshalStrict but ignores type errors for given fields
func YamlUnmarshalStrictIgnoringFields(in []byte, out interface{}, ignore []string) (err error) {
	err = yaml.UnmarshalStrict(in, out)
	if err == nil {
		return nil
	}
	errYaml, isTypeError := err.(*yaml.TypeError)
	if !isTypeError {
		return err
	}
	for _, fieldErr := range errYaml.Errors {
		if !strings.Contains(fieldErr, "not found in type") {
			// we have some other error, just return error message
			return errYaml
		}
		match := fieldNamePattern.FindStringSubmatch(fieldErr)
		if match == nil {
			// again some other error
			return errYaml
		}

		if StringSliceContains(ignore, match[1]) {
			continue
		}
		// we have type error but not for the masked fields, return error
		return errYaml
	}

	return nil
}

// convert unstructuredYamlObject to Yaml.MapSlice
func UnstructuredYamlToYamlSlice(object map[string]interface{}) yaml.MapSlice {
	m := yaml.MapSlice{}

	for k, v := range object {
		m = append(m, yaml.MapItem{Key: k, Value: v})
	}
	return m
}

// convert Yaml.MapSlice to unstructuredYamlObject
func YamlSliceToUnstructuredYaml(s yaml.MapSlice) map[string]interface{} {
	obj := make(map[string]interface{})

	for _, v := range s {
		s := v.Key.(string)
		obj[s] = v.Value
	}
	return obj
}
