package logfileexporter

import (
	"context"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
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
			log.Debug().Msgf("Not implemented config change %v", cfg)
			// TODO change config (maxlogfile count changed)
		case msg, ok := <-repoCh:
			if !ok {
				goto shutdown // channel closed
			}

			exeMsg := msg.Body.(*message.ExecuterMsg)
			if exeMsg.Subject != message.JOB_END {
				continue
			}

			// do cleanup by jobId
			m.HouseKeep(path.Join(logDir, exeMsg.JobId), cfg.MaxLogFileCount)
		}
	}

shutdown:
	messagehub.Unsubscribe(messagehub.TOPIC_JOB_REPORT, NAME)
	messagehub.Unsubscribe(messagehub.TOPIC_CONFIG_CHANGE, NAME)

	log.Debug().Msgf("%s/%s Stopped", m.GetName(), NAME)

}

// Housekeep is deleting files in directory from older (from filename)
func (m *LogFileExporter) HouseKeep(logpath string, count int) int {
	log := util.GetLoggerWithSource(m.GetName(), "HouseKeep")

	// files => alphabetical order (by os.ReadDir 's manual)
	files, err := os.ReadDir(logpath)
	if err != nil {
		log.Error().Err(err).Msg("Directory listing failed. Abort housekeeping.")
	}

	// create datetime -> filename map
	filemap := make(map[string]string)
	for _, file := range files {
		_, datetime, _, _, ok := parseFilename(file.Name())
		if ok {
			filemap[datetime] = file.Name()
		} else {
			log.Warn().Msgf("Ingored invalid filename %s", file.Name())
		}
	}

	keys := make([]string, 0, len(filemap))
	for k := range filemap {
		keys = append(keys, k)
	}
	sortedDates := sort.StringSlice(keys)

	deleteIdx := len(files) - int(count)
	deleted := 0
	for idx, key := range sortedDates {
		if idx < deleteIdx {
			fullpath := path.Join(logpath, filemap[key])
			log.Debug().Msgf("deleting %s", fullpath)
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

func getConfigFromRepository() objects.JoblogConfig {
	const KIND = "logfileexporter-config"
	bareConfig := repository.GetRepository().GetConfigByKind(KIND)
	if bareConfig != nil {
		// config exist
		return bareConfig.(objects.JoblogConfig)
	}

	ret := objects.JoblogConfig{}
	ret.Kind = KIND
	ret.MaxLogFileCount = 30
	ret.Meta.Version = objects.ObjectVersion{Major: 1, Minor: 0}

	return ret
}

func parseFilename(filename string) (JobNumber int, DatetimeStr, RunId, JobId string, ok bool) {
	strs := strings.Split(filename, FILENAME_SEP)

	if len(strs) < 4 {
		ok = false
		return
	}

	JobNumber, err := strconv.Atoi(strs[0])
	if err != nil {
		ok = false
		return
	}
	DatetimeStr = strs[1]
	if len(DatetimeStr) != 14 {
		ok = false
		return
	}
	RunId = strs[2]
	JobId = strs[3]

	ok = true
	return
}
