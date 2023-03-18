package testutils

import (
	"io/ioutil"
)

// GetYamlBytes read yaml file and returen as string.
func GetYamlBytes(path string) []byte {

	yamlbytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return yamlbytes

}
