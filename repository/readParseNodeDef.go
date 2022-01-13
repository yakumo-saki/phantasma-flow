package repository

import (
	"errors"
	"fmt"

	"github.com/yakumo-saki/phantasma-flow/global/consts"
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
	"github.com/yakumo-saki/phantasma-flow/util"
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

	err = validateNodeDef(&obj)
	if err != nil {
		msg := fmt.Sprintf("node definiton validate error %s %s %s:%s", obj.Kind, obj.Id, filepath, err.Error())
		panic(msg)
	}

	return obj
}

func validateNodeDef(def *objects.NodeDefinition) error {
	if def.Ssh.AuthType != "" {
		if !util.ContainsAny(def.Ssh.AuthType,
			consts.USER_AUTHTYPE_KEY,
			consts.USER_AUTHTYPE_KEYFILE,
			consts.USER_AUTHTYPE_PASSWORD) {
			return errors.New("Invalid AuthType " + def.Ssh.AuthType)
		}
	}

	return nil
}
