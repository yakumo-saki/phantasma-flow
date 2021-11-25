package messagehubObjects

import "github.com/yakumo-saki/phantasma-flow/pkg/objects"

// From repository to all

const DEF_REASON_ADD = "ADD"         // reason: new definition add
const DEF_REASON_CHANGE = "CHANGE"   // reason: new definition changed
const DEF_REASON_INITIAL = "INITIAL" // reason: startup all data send

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
