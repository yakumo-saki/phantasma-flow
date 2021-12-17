package scheduler

import (
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
)

// Create from jobdefinition. Filter out not needed for scheduling.
type job struct {
	id      string
	name    string
	lastRun int64
	jobMeta objects.JobMetaInfo
}

const SC_TYPE_SCHEDULE = "SC_TYPE_SCHEDULE"
const SC_TYPE_IMMEDIATE = "SC_TYPE_IMMEDIATE"

type schedule struct {
	time        int64  // Next run. unixtime
	runId       string // sha1 of uuid
	jobId       string // job ID
	reason      string // SC_TYPE_* scheduler or immediate
	scheduledAt int64  // Schedule created unixtime
	runAt       int64  // unixtime
}
