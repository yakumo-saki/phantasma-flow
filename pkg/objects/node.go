package objects

import (
	"fmt"

	"github.com/yakumo-saki/phantasma-flow/global/consts"
)

type NodeDefinition struct {
	ObjectBase  `yaml:",inline"`
	Meta        ObjectMetaBase `yaml:"meta"`
	Id          string         `yaml:"id"`          // key
	DisplayName string         `yaml:"displayName"` // display name
	NodeType    string         `yaml:"nodeType"`    // NODE_TYPE_* SSH / WinRM / local / internal
	Capacity    int            `yaml:"capacity"`    // max concurrent job point. default=2
	Ssh         SSHDefinition  `yaml:"ssh"`         // nodetype=SSH only
}

type SSHDefinition struct {
	Port          int    `yaml:"port"`          // default 22
	Host          string `yaml:"host"`          // SSH hostname or ip address
	HostAuthType  string `yaml:"hostAuthType"`  // Optional HOST_AUTHTYPE_*
	HostKey       string `yaml:"hostKey"`       // Optional host fingerprint
	User          string `yaml:"user"`          // SSH username
	Password      string `yaml:"password"`      // Password (authtype=password)
	AuthType      string `yaml:"authType"`      // Optional USER_AUTHTYPE_* (keyfile -> key -> password)
	Key           string `yaml:"key"`           // SSH key string (authtype=key) begin with ----------BEGIN OPENSSH PRIVATE KEY-----
	Keyfile       string `yaml:"keyfile"`       // SSH keyfile path (authtype=keyfile)
	KeyPassphrase string `yaml:"keyPassphrase"` // passphrase for ssh key
}

func (nd NodeDefinition) String() string {

	msg := fmt.Sprintf("ID: %s, Name: %s, Type: %s, Cap: %v Meta: %v",
		nd.Id, nd.DisplayName, nd.NodeType, nd.Capacity, nd.Meta)

	if nd.NodeType == consts.NODE_TYPE_SSH {
		key := ""
		if len(nd.Ssh.Key) > 10 {
			key = nd.Ssh.Key[:10] + "..."
		} else if len(nd.Ssh.Key) == 0 {
			key = "(none)"
		} else {
			key = "(hidden)"
		}
		msg = fmt.Sprintf("%s, Host:%s:%v, User:%s, Auth:%s Key:%s",
			msg, nd.Ssh.Host, nd.Ssh.Port, nd.Ssh.User, nd.Ssh.AuthType,
			key)
	}

	return msg
}
