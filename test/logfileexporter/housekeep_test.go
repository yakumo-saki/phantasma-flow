package logfileexporter_test

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/huandu/go-assert"
	"github.com/yakumo-saki/phantasma-flow/logexporter/logfileexporter"
	"github.com/yakumo-saki/phantasma-flow/repository"
	"github.com/yakumo-saki/phantasma-flow/test/internal/testutils"
)

func TestHouseKeepLogDelete(t *testing.T) {
	const jobId = "TestHouseKeepLogDelete"

	testutils.SetupTestLogger(t)

	a := assert.New(t)
	testutils.StartRepository()

	logpath := path.Join(repository.GetLogDirectory(), jobId)

	t.Logf("logpath=%s\n", logpath)
	testutils.PrepareEmptyDir(logpath)

	for i := 0; i < 30; i++ {
		runId := fmt.Sprintf("RunId!%02v", i)
		timestamp := 20001122034501 + i

		filename := fmt.Sprintf("%v_%s_%s", timestamp, runId, jobId)

		filepath := path.Join(logpath, filename)
		f, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, 0644)
		testutils.ErrPanic(err)
		f.Close()
	}

	lfe := logfileexporter.LogFileExporter{}
	deleted := lfe.HouseKeep(logpath, 10)

	a.Equal(30-20, deleted)
	t.Log("end")
}
