package repository

import (
	"strings"

	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
	"gopkg.in/yaml.v2"
)

func parseConfig(bytes []byte, filepath string) objects.Config {
	obj := objects.Config{}
	err := yaml.Unmarshal(bytes, &obj)
	if err != nil {
		parseYamlPanic("config", obj.Kind, "", err.Error(), filepath)
	}
	if !strings.HasSuffix(obj.Kind, "config") {
		parseYamlPanic("config", obj.Kind, "", "config kind must be suffixed by 'config'", filepath)
	}

	return obj
}
