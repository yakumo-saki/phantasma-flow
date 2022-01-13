package repository

import (
	"os"
	"path"

	"github.com/yakumo-saki/phantasma-flow/global/consts"
	"github.com/yakumo-saki/phantasma-flow/util"
)

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

	p.Home = os.Getenv(consts.ENV_HOME_DIR)
	p.Def = os.Getenv(consts.ENV_DEF_DIR)
	p.Data = os.Getenv(consts.ENV_DATA_DIR)
	p.Temp = os.Getenv(consts.ENV_TEMP_DIR)
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

	isNotGoodDir(p.Home, consts.ENV_HOME_DIR)
	isNotGoodDir(p.Def, consts.ENV_DEF_DIR)
	isNotGoodDir(p.Data, consts.ENV_DATA_DIR)
	isNotGoodDir(p.Temp, consts.ENV_TEMP_DIR)
	isNotGoodDir(p.NodeDef, p.NodeDef)
	isNotGoodDir(p.ConfigDef, p.ConfigDef)
	isNotGoodDir(p.JobDef, p.JobDef)
	isNotGoodDir(p.JobLog, p.JobLog)
	isNotGoodDir(p.JobMeta, p.JobMeta)

	makeSureDirExists(p)

	os.Setenv(consts.ENV_HOME_DIR, p.Home)
	os.Setenv(consts.ENV_DEF_DIR, p.Def)
	os.Setenv(consts.ENV_DATA_DIR, p.Data)
	os.Setenv(consts.ENV_TEMP_DIR, p.Temp)

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
