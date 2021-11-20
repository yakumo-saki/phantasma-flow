package objects

import "fmt"

type NodeDefinition struct {
	ObjectBase
	Meta     ObjectMetaBase `yaml:"meta"`
	Name     string         `yaml:"name"`
	NodeType string         `yaml:"nodeType"` // SSH / WinRM / internal
}

func (nd NodeDefinition) String() string {
	return fmt.Sprintf("Name: %s, Type: %s, Meta: %v",
		nd.Name, nd.NodeType, nd.Meta)
}
