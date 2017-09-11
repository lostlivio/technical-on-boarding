package jobs

import (
	"time"

	"github.com/revel/cron"

	"github.com/revel/modules/jobs/app/jobs"
)

type Event struct {
	SessionID int    // The user session id
	Type      string // "start", "progress", "complete", and "error"
	Timestamp int    // Unix timestamp (secs)
	Text      string // What the job progress is (if Type == "progress")
	Error     error  // Source error (if Type == "error")
}

func NewEvent(sid int, typ string, msg string) Event {
	return Event{sid, typ, int(time.Now().Unix()), msg, nil}
}

func NewError(sid int, msg string, err error) Event {
	return Event{sid, "error", int(time.Now().Unix()), msg, err}
}

func StartJob(job cron.Job) {
	jobs.Now(job)
}
