package repository

import (
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
	"gopkg.in/yaml.v2"
)

func parseNodeDef(bytes []byte) objects.NodeDefinition {
	obj := objects.NodeDefinition{}
	err := yaml.Unmarshal(bytes, &obj)
	if err != nil {
		parseYamlPanic("node definition", obj.Kind, obj.Id, err.Error())
	}
	if obj.Kind != objects.KIND_NODE_DEF {
		parseYamlPanic("node definition", obj.Kind, obj.Id, "")
	}
	if obj.DisplayName == "" {
		obj.DisplayName = obj.Id
	}

	return obj
}
