package objects

import (
	"fmt"

	"github.com/yakumo-saki/phantasma-flow/util"
)

type ObjectBase struct {
	Kind string `json:"kind"`
}

type ObjectMetaBase struct {
	Version   ObjectVersion `json:"version"`   // Version
	Created   string        `json:"created"`   // ISO8601 yyyy-mm-dd hh:mm:ssZ
	CreatedBy string        `json:"createdBy"` // Username
}

func (omb ObjectMetaBase) String() string {
	created := util.Nvl(omb.Created, "(not set)")
	by := util.Nvl(omb.CreatedBy, "(not set)")

	return fmt.Sprintf("%s ,Created %s by %s",
		&omb.Version, created, by)
}

// Version of Objects
type ObjectVersion struct {
	Major uint
	Minor uint
}

func (ov ObjectVersion) String() string {
	return fmt.Sprintf("Version %d.%d", ov.Major, ov.Minor)
}
