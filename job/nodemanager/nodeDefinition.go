package nodemanager

import (
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/yakumo-saki/phantasma-flow/global/consts"
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
	"github.com/yakumo-saki/phantasma-flow/util"
)

//
func SetDefaultNodeDefinition(orgNodeDef objects.NodeDefinition) (*objects.NodeDefinition, error) {

	nodeDef := objects.NodeDefinition{}
	util.DeepCopy(orgNodeDef, &nodeDef)

	if nodeDef.Id == "" {
		return nil, errors.New("NodeDef ID is empty")
	}

	if nodeDef.DisplayName == "" {
		nodeDef.DisplayName = nodeDef.Id
	}

	if nodeDef.NodeType == "" {
		return nil, errors.New("NodeDef NodeType is empty")
	}

	if nodeDef.Capacity == 0 {
		nodeDef.Capacity = 2
	}

	// ssh
	if nodeDef.NodeType == consts.NODE_TYPE_SSH {
		err := setDefaultNodeDefinitionSSH(&nodeDef)
		if err != nil {
			return nil, err
		}
	}

	return &nodeDef, nil
}

func setDefaultNodeDefinitionSSH(nodeDef *objects.NodeDefinition) error {
	if nodeDef.Ssh.Port == 0 {
		nodeDef.Ssh.Port = 22
	}
	if nodeDef.Ssh.AuthType == "" {
		if nodeDef.Ssh.Keyfile != "" {
			nodeDef.Ssh.AuthType = consts.USER_AUTHTYPE_KEYFILE
		} else if nodeDef.Ssh.Key != "" {
			nodeDef.Ssh.AuthType = consts.USER_AUTHTYPE_KEY
		} else if nodeDef.Ssh.Password != "" {
			nodeDef.Ssh.AuthType = consts.USER_AUTHTYPE_PASSWORD
		}
	}

	// read ssh key file and set to KEY
	if nodeDef.Ssh.AuthType == consts.USER_AUTHTYPE_KEYFILE {
		nodeDef.Ssh.AuthType = consts.USER_AUTHTYPE_KEY
		keystr := util.ReadPublicKeyfile(nodeDef.Ssh.Keyfile)
		nodeDef.Ssh.Key = keystr
		log.Info().Msgf("Read ok -> %s", nodeDef.Ssh.Keyfile)
	}

	if nodeDef.Ssh.HostAuthType == "" {
		if nodeDef.Ssh.HostKey != "" {
			nodeDef.Ssh.HostAuthType = consts.HOST_AUTHTYPE_KEY
		} else {
			nodeDef.Ssh.HostAuthType = consts.HOST_AUTHTYPE_IGNORE
		}
	}

	return nil
}
