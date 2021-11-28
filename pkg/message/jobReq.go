package message

import "github.com/yakumo-saki/phantasma-flow/pkg/objects"

type definition struct {
	Reason string // ADD , CHANGE , INITIAL
}

type NodeDefinitionMsg struct {
	definition
	NodeDefinition objects.NodeDefinition
}

type JobDefinitionMsg struct {
	definition
	JobDefinition objects.JobDefinition
}
