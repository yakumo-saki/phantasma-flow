package objects

import "fmt"

const KIND_PHFLOW_CFG = "phantasma-flow-config"
const KIND_JOBLOG_CFG = "joblog-config"
const KIND_PPROF_SERVER_CFG = "pprof-server-config"

// Config type is base struct of All config structs
// Config structs are serialized to yaml
type Config struct {
	ObjectBase `yaml:",inline"`
	Meta       ObjectMetaBase `yaml:"meta"`
}

func (c Config) String() string {
	return fmt.Sprintf("Kind: %s, Meta: %v",
		c.Kind, c.Meta)
}

// PhantasmaFlowConfig is config object for Phantasma-Flow.
//  kind is KIND_PHFLOW_CFG
type PhantasmaFlowConfig struct {
	Config `yaml:",inline"`

	Security struct {
		Mode string // none, (not impl -> apikey)
	}
}

// JoblogConfig is config object for job logs.
//  kind is KIND_LOGFILE_EXPORTER_CFG
type JoblogConfig struct {
	Config `yaml:",inline"`

	MaxLogFileCount int `yaml:"logFileCount"`   // Max logfile count per jobId (If set 0 => default => 30)
	JobResultCount  int `yaml:"jobResultCount"` // Max result count per jobId (metalog)
}

// LogFileExporterConfig is config object for pprofserver.
//  kind is KIND_PPROF_SERVER_CFG
type PprofServerConfig struct {
	Config `yaml:",inline"`

	Enabled           bool   `yaml:"enabled"`           // true to start server
	ListenAddrAndPort string `yaml:"listenAddrAndPort"` // ex) localhost:6060 (default)
}

// GeneralConfig is config object for fallback
//  Use of this is not recommended.
type GeneralConfig struct {
	Config `yaml:",inline"`

	ConfigMap map[string]interface{}
}
