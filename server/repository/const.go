package repository

type objectType int

const (
	JOB objectType = iota
	NODE
	CONFIG
)

func (ot objectType) String() string {
	switch ot {
	case JOB:
		return "JobDefinition"
	case NODE:
		return "NodeDefinition"
	case CONFIG:
		return "Config"
	default:
		return "Unknown"
	}
}
