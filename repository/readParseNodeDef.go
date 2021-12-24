package repository

import (
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
	"gopkg.in/yaml.v2"
)

func parseNodeDef(bytes []byte, filepath string) objects.NodeDefinition {
	obj := objects.NodeDefinition{}
	err := yaml.Unmarshal(bytes, &obj)
	if err != nil {
		parseYamlPanic("node definition", obj.Kind, obj.Id, err.Error(), filepath)
	}
	if obj.Kind != objects.KIND_NODE_DEF {
		parseYamlPanic("node definition", obj.Kind, obj.Id, "", filepath)
	}
	if obj.DisplayName == "" {
		obj.DisplayName = obj.Id
	}

	return obj
}
