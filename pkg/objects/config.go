package objects

import "fmt"

type Config struct {
	ObjectBase `json:",inline"`
	Meta       ObjectMetaBase `json:"meta"`
}

func (c Config) String() string {
	return fmt.Sprintf("Kind: %s, Meta: %v",
		c.Kind, c.Meta)
}

type LogCollecterConfig struct {
	Config
}

type JobSchedulerConfig struct {
	Config
}
