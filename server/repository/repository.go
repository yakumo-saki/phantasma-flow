package repository

import (
	"io/ioutil"
	"path/filepath"
	"sync"

	"github.com/yakumo-saki/phantasma-flow/global/consts"
	"github.com/yakumo-saki/phantasma-flow/messagehub"
	"github.com/yakumo-saki/phantasma-flow/pkg/message"
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
	"github.com/yakumo-saki/phantasma-flow/util"
)

const myname = "Repository"

type Repository struct {
	Initialized bool

	mutex   sync.Mutex
	nodes   []objects.NodeDefinition
	jobs    []objects.JobDefinition
	configs map[string]interface{} // map[kind] -> struct

	paths phflowPath
}

func (r *Repository) Initialize() error {
	log := util.GetLoggerWithSource(myname, "initialize")
	log.Debug().Msg("Repository initialize start")

	r.paths = aquirePhflowPath()
	r.configs = make(map[string]interface{})

	dirType := map[objectType]string{
		NODE:   r.paths.NodeDef,
		CONFIG: r.paths.ConfigDef,
		JOB:    r.paths.JobDef,
	}

	log.Info().Msgf("%s=%s", consts.ENV_HOME_DIR, r.paths.Home)
	log.Info().Msgf("%s=%s", consts.ENV_DEF_DIR, r.paths.Def)
	log.Info().Msgf("%s=%s", consts.ENV_DATA_DIR, r.paths.Data)

	for typ, dirPath := range dirType {
		// log.Debug().Msgf("Reading %s from %s", typ, dirPath)
		err := r.readAllYaml(dirPath, typ)
		if err != nil {
			return err
		}
	}

	// if r.Dump() {
	// 	log.Debug().Msg("Repository initialized")
	// 	return errors.New("not impremented")
	// }

	r.Initialized = true
	log.Info().Msg("Repository initialized")

	return nil
}

func (repo *Repository) Dump() bool {
	log := util.GetLoggerWithSource(myname, "dump")
	log.Debug().Msg("Jobs")
	for _, v := range repo.jobs {
		log.Debug().Msgf("%s", v)
		log.Debug().Msg("")
	}
	log.Debug().Msg("-------------------------------------")
	log.Debug().Msg("Nodes")
	for _, v := range repo.nodes {
		log.Debug().Msgf("%s", v)
		log.Debug().Msg("")
	}
	log.Debug().Msg("-------------------------------------")
	log.Debug().Msg("Configs")
	for _, v := range repo.configs {
		log.Debug().Msgf("%s", v)
		log.Debug().Msg("")
	}
	log.Debug().Msg("-------------------------------------")
	return false
}

func (repo *Repository) readAllYaml(path string, objType objectType) error {
	log := util.GetLoggerWithSource(myname, "readYaml")
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Err(err)
	}

	for _, fileinfo := range files {
		if fileinfo.IsDir() {
			continue
		}

		fp := filepath.Join(path, fileinfo.Name())
		bytes, err := ioutil.ReadFile(fp)
		if err != nil {
			return err
		}

		switch objType {
		case NODE:
			obj := parseNodeDef(bytes, fp)
			repo.nodes = append(repo.nodes, obj)
		case JOB:
			obj := parseJobDef(bytes, fp)
			repo.jobs = append(repo.jobs, obj)
		case CONFIG:
			kind, obj := parseConfig(bytes, fp)
			repo.configs[kind] = obj
		}

	}

	return nil
}

func (repo *Repository) GetConfigByKind(kind string) interface{} {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	obj, ok := repo.configs[kind]
	if !ok {
		return nil
	}

	return obj
}

func (repo *Repository) GetJobById(jobId string) *objects.JobDefinition {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	for _, j := range repo.jobs {
		if j.Id == jobId {
			jCopy := &j
			return jCopy
		}
	}

	return nil
}

func (repo *Repository) SendAllNodes() int {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	sent := 0
	for _, v := range repo.nodes {
		nodeMsg := message.NodeDefinitionMsg{}
		nodeMsg.Reason = message.DEF_REASON_INITIAL
		nodeMsg.NodeDefinition = v
		messagehub.Post(messagehub.TOPIC_NODE_DEFINITION, nodeMsg)
		sent++
	}
	return sent
}

func (repo *Repository) SendAllJobs() int {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	sent := 0
	for _, v := range repo.jobs {
		jobMsg := message.JobDefinitionMsg{}
		jobMsg.Reason = message.DEF_REASON_INITIAL
		jobMsg.JobDefinition = v
		messagehub.Post(messagehub.TOPIC_JOB_DEFINITION, jobMsg)
		sent++
	}
	return sent
}
