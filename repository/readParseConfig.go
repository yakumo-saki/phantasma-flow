package repository

import (
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
	"gopkg.in/yaml.v2"
)

func parseConfig(bytes []byte) objects.Config {
	obj := objects.Config{}
	err := yaml.Unmarshal(bytes, &obj)
	if err != nil {
		parseYamlPanic("config", obj.Kind, "", err.Error())
	}

	return obj
}
