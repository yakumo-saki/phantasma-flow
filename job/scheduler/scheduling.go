package scheduler

import (
	"time"

	"github.com/robfig/cron/v3"
	"github.com/thanhpk/randstr"
	"github.com/yakumo-saki/phantasma-flow/pkg/objects"
	"github.com/yakumo-saki/phantasma-flow/util"
)

// Schedule next run, single job with lock
func (js *JobScheduler) schedule(jobId string, now time.Time) {
	js.mutex.Lock()
	defer js.mutex.Unlock()
	js.scheduleWithoutLock(jobId, now)
}

// Schedule next run, single job without lock
func (js *JobScheduler) scheduleWithoutLock(jobId string, now time.Time) {
	log := util.GetLoggerWithSource(js.GetName(), "schedule")

	// Delete from exisiting schedule where jobId = thisJobId
	for e := js.schedules.Front(); e != nil; e = e.Next() {
		sche := e.Value.(schedule)
		if sche.jobId == jobId {
			js.schedules.Remove(e)
		}
	}

	job := js.jobs[jobId]

	// Check and calc next schedule or no schedule
	nextSchedule := js.calcNextSchedule(now, job.jobMeta.Schedules)
	if nextSchedule == -1 {
		// it is valid, running from run immediate feature or disabled job
		log.Debug().Str("JobId", jobId).Msg("Has no schedule. only register.")
	} else {
		log.Debug().Str("JobId", jobId).Int64("Next", nextSchedule).Msg("New schedule")
	}

	// push next run schedule
	newSchedule := schedule{}
	newSchedule.runId = js.generateRunId()
	newSchedule.jobId = jobId
	newSchedule.scheduledAt = now.Unix()
	newSchedule.reason = SC_TYPE_SCHEDULE

	newSchedule.time = nextSchedule
	js.schedules.PushFront(newSchedule)
}

func (js *JobScheduler) generateRunId() string {
	return randstr.String(8)
}

// return unixtime
func (js *JobScheduler) calcNextSchedule(now time.Time, schedules []objects.JobSchedule) int64 {
	log := util.GetLoggerWithSource(js.GetName(), "schedule")
	nextRun := []int64{}

	if len(schedules) == 0 {
		return -1
	}

	for _, sche := range schedules {
		parser := js.getCronParser()

		sc, err := parser.Parse(sche.Param)
		if err != nil {
			log.Err(err).Msg("Schedule parse error")
			panic("schedule error " + sche.Param)
		}

		now := time.Now()
		ret := sc.Next(now)

		// log.Debug().Str("next", ret.Format("2006-01-02 15:04:05")).Int64("nextunix", ret.Unix()).Msg("Next schedule")
		nextRun = append(nextRun, ret.Unix())
	}

	return util.MinInt64(nextRun...)
}

func (js *JobScheduler) getCronParser() cron.Parser {
	return cron.NewParser(cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
}
