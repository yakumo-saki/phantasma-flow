package objects

import "fmt"

const NODE_LOCAL = "local"
const NODE_SSH = "ssh"
const NODE_WINRM = "winrm"

type NodeDefinition struct {
	ObjectBase
	Meta        ObjectMetaBase `json:"meta"`
	Id          string         `json:"id"`          // key
	DisplayName string         `json:"displayName"` // display name
	NodeType    string         `json:"nodeType"`    // SSH / WinRM / local / internal
	Capacity    int            `json:"capacity"`    // max concurrent job point.
}

func (nd NodeDefinition) String() string {
	return fmt.Sprintf("ID: %s, Name: %s, Type: %s, Cap: %v Meta: %v",
		nd.Id, nd.DisplayName, nd.NodeType, nd.Capacity, nd.Meta)
}
