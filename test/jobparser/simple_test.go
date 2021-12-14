package jobparser

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yakumo-saki/phantasma-flow/job/jobparser"
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
	"github.com/yakumo-saki/phantasma-flow/test/internal/testutils"
	"gopkg.in/yaml.v3"
)

func TestJobParserSimple(t *testing.T) {
	fp := path.Join(testutils.GetTestJobDefDir(), "jobparser_simple.yaml")
	yamlStr := testutils.GetYamlBytes(fp)

	jobDef := objects.JobDefinition{}
	yaml.Unmarshal(yamlStr, &jobDef)

	execJobs, err := jobparser.BuildFromJobDefinition(&jobDef, "jobId", "runId")
	if err != nil {
		panic(err)
	}

	// asserts
	elem := execJobs.Front()
	assert.NotNil(t, elem, "step1")
	step := elem.Value.(jobparser.ExecutableJobStep)
	step1 := step
	assert.Equal(t, "step1", step1.Name, "step1")

	elem = elem.Next()
	assert.NotNil(t, elem)
	step = elem.Value.(jobparser.ExecutableJobStep)
	step2 := step
	assert.Equal(t, "step2", step2.Name, "step2")

	t.Log("ok")
}
