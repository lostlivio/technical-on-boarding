package jobs

import (
	"time"
)

type Event struct {
	SessionID int    // The user session id
	Type      string // "start", "progress", "complete-success", and "complete-fail"
	Timestamp int    // Unix timestamp (secs)
	Text      string // What the job progress is (if Type == "progress")
}

func NewEvent(sid int, typ string, msg string) Event {
	return Event{sid, typ, int(time.Now().Unix()), msg}
}

/*
func StartWorkload(user *models.User) (int, error) {
func StartWorkload(job *GenerateProject) (int, error) {
	jobs.Now(job)
	return 0, nil
}
*/
