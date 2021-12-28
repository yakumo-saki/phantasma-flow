package repository

import (
	"fmt"
	"strings"

	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
	"gopkg.in/yaml.v2"
)

// parseConfig returns kind and object. objects are objects.XxxConfig
func parseConfig(bytes []byte, filepath string) (string, interface{}) {
	obj := objects.Config{}
	err := yaml.Unmarshal(bytes, &obj)
	if err != nil {
		parseYamlPanic("config", obj.Kind, "", err.Error(), filepath)
	}
	if !strings.HasSuffix(obj.Kind, "config") {
		parseYamlPanic("config", obj.Kind, "", "config kind must be suffixed by 'config'", filepath)
	}

	parseYamlOrPanic := func(bytes []byte, out interface{}) {
		err := yaml.Unmarshal(bytes, out)
		if err != nil {
			parseYamlPanic("config", obj.Kind, "", err.Error(), filepath)
		}
	}

	var ret interface{}
	switch obj.Kind {
	case objects.KIND_PHFLOW_CFG:
		o := objects.PhantasmaFlowConfig{}
		parseYamlOrPanic(bytes, &o)
		ret = o
	case objects.KIND_JOBLOG_CFG:
		o := objects.JoblogConfig{}
		parseYamlOrPanic(bytes, &o)
		ret = o
	case objects.KIND_PPROF_SERVER_CFG:
		o := objects.PprofServerConfig{}
		parseYamlOrPanic(bytes, &o)
		ret = o
	default:
		msg := fmt.Sprintf("Unknown config kind: %s", obj.Kind)
		panic(msg)
	}

	return obj.Kind, ret
}
