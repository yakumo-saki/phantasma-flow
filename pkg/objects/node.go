package objects

import "fmt"

const NODE_TYPE_LOCAL = "local"
const NODE_TYPE_SSH = "ssh"
const NODE_TYPE_WINRM = "winrm"

type NodeDefinition struct {
	ObjectBase  `yaml:",inline"`
	Meta        ObjectMetaBase `yaml:"meta"`
	Id          string         `yaml:"id"`          // key
	DisplayName string         `yaml:"displayName"` // display name
	NodeType    string         `yaml:"nodeType"`    // SSH / WinRM / local / internal
	Capacity    int            `yaml:"capacity"`    // max concurrent job point.
}

func (nd NodeDefinition) String() string {
	return fmt.Sprintf("ID: %s, Name: %s, Type: %s, Cap: %v Meta: %v",
		nd.Id, nd.DisplayName, nd.NodeType, nd.Capacity, nd.Meta)
}
