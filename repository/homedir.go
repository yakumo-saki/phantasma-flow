package repository

import (
	"os"
	"path"

	"github.com/yakumo-saki/phantasma-flow/util"
)

// External configable env value
const ENV_HOME_DIR = "PHFLOW_HOME"
const ENV_DEF_DIR = "PHFLOW_DEF_DIR"
const ENV_DATA_DIR = "PHFLOW_DATA_DIR"
const ENV_TEMP_DIR = "PHFLOW_TEMP_DIR"

type phflowPath struct {
	Home      string
	Def       string
	Data      string
	Temp      string
	NodeDef   string
	ConfigDef string
	JobDef    string
	JobLog    string
	JobMeta   string
}

// Get phantasma-flow home path and set it to ENV values
// ENV or ~/.config/phantasma-flow
// if fail, cause PANIC
func aquirePhflowPath() phflowPath {
	util.GetLoggerWithSource(myname, "phflowPath")

	p := phflowPath{}

	p.Home = os.Getenv(ENV_HOME_DIR)
	p.Def = os.Getenv(ENV_DEF_DIR)
	p.Data = os.Getenv(ENV_DATA_DIR)
	p.Temp = os.Getenv(ENV_TEMP_DIR)
	if p.Home == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			panic("Get home fail, Please set PHFLOW_HOME environment value.")
		}
		p.Home = path.Join(home, ".config", "phantasma-flow")
	}
	if p.Def == "" {
		p.Def = path.Join(p.Home, "definitions")
	}
	if p.Data == "" {
		p.Data = path.Join(p.Home, "data")
	}
	if p.Temp == "" {
		p.Temp = path.Join(p.Home, "temp")
	}

	p.JobDef = path.Join(p.Def, "job")
	p.NodeDef = path.Join(p.Def, "node")
	p.ConfigDef = path.Join(p.Def, "config")
	p.JobLog = path.Join(p.Data, "log")
	p.JobMeta = path.Join(p.Data, "meta")

	isNotGoodDir(p.Home, ENV_HOME_DIR)
	isNotGoodDir(p.Def, ENV_DEF_DIR)
	isNotGoodDir(p.Data, ENV_DATA_DIR)
	isNotGoodDir(p.Temp, ENV_TEMP_DIR)
	isNotGoodDir(p.NodeDef, p.NodeDef)
	isNotGoodDir(p.ConfigDef, p.ConfigDef)
	isNotGoodDir(p.JobDef, p.JobDef)
	isNotGoodDir(p.JobLog, p.JobLog)
	isNotGoodDir(p.JobMeta, p.JobMeta)

	makeSureDirExists(p)

	return p

}

func isNotGoodDir(dirname string, name string) {

	if dirname != "" {
		st, err := os.Stat(dirname)
		if os.IsNotExist(err) {
			// not exist is ok. try to create after this
			return
		}
		if st == nil {
			return
		}
		if !st.IsDir() {
			panic(name + " must be directory:" + dirname)
		}
	}

}

func makeSureDirExists(p phflowPath) {
	log := util.GetLoggerWithSource(myname, "phflowPath")

	util.MkdirAll(p.Home, &log)
	util.MkdirAll(p.Def, &log)
	util.MkdirAll(p.Data, &log)
	util.MkdirAll(p.Temp, &log)
	util.MkdirAll(p.NodeDef, &log)
	util.MkdirAll(p.ConfigDef, &log)
	util.MkdirAll(p.JobDef, &log)
	util.MkdirAll(p.JobLog, &log)
	util.MkdirAll(p.JobMeta, &log)
}
