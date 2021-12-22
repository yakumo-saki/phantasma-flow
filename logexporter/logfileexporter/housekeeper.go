package logfileexporter

import (
	"context"
	"os"
	"path"
	"sync"

	"github.com/yakumo-saki/phantasma-flow/messagehub"
	"github.com/yakumo-saki/phantasma-flow/pkg/message"
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
	"github.com/yakumo-saki/phantasma-flow/repository"
	"github.com/yakumo-saki/phantasma-flow/util"
)

// repoからconfig取得
// configの数以上ログファイルがあったら古いものを削除（ファイル名から判断）
// logディレクトリ以下で、0ファイルなディレクトリは削除
// タイミングはJOB COMPLETEでOK

// main method of LogListener
func (m *LogFileExporter) HouseKeeper(ctx context.Context, startUp, shutdown *sync.WaitGroup) {
	const NAME = "HouseKeeper"
	log := util.GetLoggerWithSource(m.GetName(), NAME)

	cfg := getConfigFromRepository()
	logDir := repository.GetLogDirectory()

	repoCh := messagehub.Subscribe(messagehub.TOPIC_JOB_REPORT, NAME)
	cfgCh := messagehub.Subscribe(messagehub.TOPIC_CONFIG_CHANGE, NAME)

	startUp.Done()
	defer shutdown.Done()

	for {
		select {
		case <-ctx.Done():
			goto shutdown
		case cfg, ok := <-cfgCh:
			if !ok {
				goto shutdown // channel closed
			}
			log.Debug().Msgf("%v", cfg)
		case msg, ok := <-repoCh:
			if !ok {
				goto shutdown // channel closed
			}

			exeMsg := msg.Body.(*message.ExecuterMsg)
			if exeMsg.Reason != message.JOB_END {
				continue
			}

			// do cleanup by jobId
			m.HouseKeep(path.Join(logDir, exeMsg.JobId), cfg.MaxLogFileCount)
		}
	}

shutdown:
	log.Debug().Msgf("%s/%s Stopped", m.GetName(), NAME)

}

// Housekeep is deleting files in directory from older (from filename)
func (m *LogFileExporter) HouseKeep(logpath string, count uint) int {
	log := util.GetLoggerWithSource(m.GetName(), "HouseKeep")

	// files => alphabetical order (by os.ReadDir 's manual)
	files, err := os.ReadDir(logpath)
	if err != nil {
		log.Error().Err(err).Msg("Directory listing failed. Abort housekeeping.")
	}

	deleteIdx := len(files) - int(count)
	deleted := 0
	for idx, file := range files {
		if idx < deleteIdx {
			fullpath := path.Join(logpath, file.Name())
			err := os.Remove(fullpath)
			if err == nil {
				deleted++
			} else {
				log.Warn().Msgf("Failed to delete %s", fullpath)
				continue
			}
		}
	}

	return deleted
}

func getConfigFromRepository() *objects.LogFileExporterConfig {
	const KIND = "logfileexporter-config"
	bareConfig := repository.GetRepository().GetConfigByKind(KIND)
	if bareConfig != nil {
		// config exist
		return bareConfig.(*objects.LogFileExporterConfig)
	}

	ret := objects.LogFileExporterConfig{}
	ret.Kind = KIND
	ret.MaxLogFileCount = 30
	ret.Meta.Version = objects.ObjectVersion{Major: 1, Minor: 0}

	return &ret
}
