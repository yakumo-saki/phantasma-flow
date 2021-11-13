package objects

import "fmt"

type JobDefinition struct {
	ObjectBase
	Meta  ObjectMetaBase
	Steps []JobStepDefinition
	Name  string `yaml:"name"`
}

func (nd JobDefinition) String() string {
	ret := fmt.Sprintf("Name: %s Meta: %v", nd.Name, nd.Meta)

	ret = ret + "\n"
	for _, st := range nd.Steps {
		ret = ret + fmt.Sprintf("Step: %v\n", st)
	}
	return ret
}

type JobStepDefinition struct {
	Name string
}

func (nd JobStepDefinition) String() string {
	return nd.Name
}
