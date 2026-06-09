package rdbms

// ////////////////////////////////////////////////////////////////////////////////// //

type BaculaJob struct {
	Name      string `db:"name"`
	Pool      string `db:"pool"`
	Level     string `db:"level"`
	Status    string `db:"jobstatus"`
	SchedTime uint32 `db:"schedtime"`
	StartTime uint32 `db:"starttime"`
	EndTime   uint32 `db:"endtime"`
	JobBytes  uint64 `db:"jobbytes"`
	JobFiles  uint64 `db:"jobfiles"`
}

type BaculaJobSummary struct {
	Name          string `db:"name"`
	Pool          string `db:"pool"`
	Level         string `db:"level"`
	TotalJobBytes uint64 `db:"totaljobbytes"`
	TotalJobFiles uint64 `db:"totaljobfiles"`
}

type BaculaStoredData struct {
	Name        string `db:"name"`
	Pool        string `db:"pool"`
	StoredBytes uint64 `db:"stored_bytes"`
	StoredFiles uint64 `db:"stored_files"`
}

type BaculaSummary struct {
	ScheduledJobs uint32 `db:"scheduledjobs"`
}

// ////////////////////////////////////////////////////////////////////////////////// //
