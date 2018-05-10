package codegen

import (
	"encoding/json"
	"os"
	"strings"
	"text/template"
)

func WriteTemplateFile(filename string, data interface{}, templateString string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer file.Close()

	tpl, _ := template.New("").Parse(templateString)

	return tpl.Execute(file, data)
}

func DeserializeJsonFile(path string, data interface{}) error {
	file, openErr := os.Open(path)
	if openErr != nil {
		return openErr
	}
	defer file.Close()

	if jsonErr := json.NewDecoder(file).Decode(data); jsonErr != nil {
		return jsonErr
	}

	return nil
}

func isUppercase(input string) bool {
	return strings.ToLower(input) != input
}

func isCustomType(typeName string) bool {
	return isUppercase(typeName[0:1])
}
