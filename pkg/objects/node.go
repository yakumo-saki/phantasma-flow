package objects

import "fmt"

type NodeDefinition struct {
	ObjectBase
	Meta            ObjectMetaBase `json:"meta"`
	Id              string         `json:"id"`              // key
	Name            string         `json:"name"`            // display name
	NodeType        string         `json:"nodeType"`        // SSH / WinRM / internal
	MaxParallelJobs int            `json:"maxParallelJobs"` // max concurrent job
}

func (nd NodeDefinition) String() string {
	return fmt.Sprintf("ID: %s, Name: %s, Type: %s, Meta: %v",
		nd.Id, nd.Name, nd.NodeType, nd.Meta)
}
