package jobparser

import (
	"path"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/huandu/go-assert"
	"github.com/yakumo-saki/phantasma-flow/job/jobparser"
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
	"github.com/yakumo-saki/phantasma-flow/test/internal/testutils"
)

func TestJobParserSimple(t *testing.T) {
	a := assert.New(t)

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
	a.NotEqual(elem, nil)
	step := elem.Value.(jobparser.ExecutableJobStep)
	step1 := step
	a.Equal("step1", step1.Name)

	elem = elem.Next()
	a.NotEqual(elem, nil)
	step = elem.Value.(jobparser.ExecutableJobStep)
	step2 := step
	a.Equal("step2", step2.Name)

	t.Log(step1.UseCapacity)
}
