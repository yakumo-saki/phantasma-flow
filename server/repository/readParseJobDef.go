package repository

import (
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
	"gopkg.in/yaml.v2"
)

func parseJobDef(bytes []byte, filepath string) objects.JobDefinition {
	obj := objects.JobDefinition{}
	err := yaml.Unmarshal(bytes, &obj)
	if err != nil {
		parseYamlPanic("node definition", obj.Kind, obj.Id, err.Error(), filepath)
	}
	if obj.Kind != objects.KIND_JOB_DEF {
		parseYamlPanic("job definition", obj.Kind, obj.Id, "", filepath)
	}

	return obj
}
