package entity

type JobStatus string

const (
	JobPending JobStatus = "PENDING"
	JobRunning JobStatus = "RUNNING"
	JobDone    JobStatus = "DONE"
	JobFailed  JobStatus = "FAILED"
)

type Job struct {
	ID         string
	InputPath  string
	OutputPath string
	MimeType   string
	Status     JobStatus
	Progress   int
	Error      string
}
