package repository

import (
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
	"gopkg.in/yaml.v2"
)

func parseJobDef(bytes []byte) objects.JobDefinition {
	obj := objects.JobDefinition{}
	err := yaml.Unmarshal(bytes, &obj)
	if err != nil {
		parseYamlPanic("node definition", obj.Kind, obj.Id, err.Error())
	}
	if obj.Kind != objects.KIND_JOB_DEF {
		parseYamlPanic("job definition", obj.Kind, obj.Id, "")
	}

	return obj
}
