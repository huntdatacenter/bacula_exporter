package rdbms

// ////////////////////////////////////////////////////////////////////////////////// //

// GetLatestJobs return all jobs from the last 30 days, one row per job execution
func (db *DB) GetLatestJobs() ([]*BaculaJob, error) {
        baculaJobs := make([]*BaculaJob, 0)

        sqlState := `
          SELECT
                j.Name as name,
                p.Name as pool,
                j.Level as level,
                j.JobId as jobid,
                j.JobStatus as jobstatus,
                coalesce(extract(epoch from j.SchedTime), 0)::integer as SchedTime,
                coalesce(extract(epoch from j.StartTime), 0)::integer as StartTime,
                coalesce(extract(epoch from j.EndTime), 0)::integer as EndTime,
                j.JobBytes::bigint as jobbytes,
                j.JobFiles::bigint as jobfiles
          FROM
                Job j
                JOIN Pool p ON j.PoolId = p.PoolId
          WHERE
                j.Type = 'B'
                AND j.StartTime > NOW() - INTERVAL '30 days'
          ORDER BY
                j.JobId`

        err := db.Select(&baculaJobs, sqlState)

        return baculaJobs, err
}

// GetStoredData returns the sum of bytes and files for currently stored jobs,
// grouped by job name and pool. Jobs on recycled or purged volumes are excluded.
func (db *DB) GetStoredData() ([]*BaculaStoredData, error) {
        storedData := make([]*BaculaStoredData, 0)

        sqlState := `
          SELECT
                j.Name as name,
                p.Name as pool,
                SUM(j.JobBytes)::bigint as stored_bytes,
                SUM(j.JobFiles)::bigint as stored_files
          FROM
                Job j
                JOIN Pool p ON j.PoolId = p.PoolId
          WHERE
                j.Type = 'B'
                AND j.JobStatus = 'T'
                AND j.PurgedFiles = 0
                AND EXISTS (
                      SELECT 1
                      FROM JobMedia jm
                      JOIN Media m ON jm.MediaId = m.MediaId
                      WHERE jm.JobId = j.JobId
                            AND m.VolStatus NOT IN ('Purge', 'Recycle', 'Error', 'Missing')
                )
          GROUP BY
                j.Name, p.Name`

        err := db.Select(&storedData, sqlState)

        return storedData, err
}

// GetJobsSummary return summary of all jobs
func (db *DB) GetJobsSummary() ([]*BaculaJobSummary, error) {
        jobsSummary := make([]*BaculaJobSummary, 0)

        sqlState := `
          SELECT
                j.Name as name,
                p.Name as pool,
                j.Level as level,
                SUM(j.JobBytes)::bigint as TotalJobBytes,
                SUM(j.JobFiles)::bigint as TotalJobFiles
          FROM
                Job j
                JOIN Pool p ON j.PoolId = p.PoolId
          WHERE
                j.Name IN (
                      SELECT DISTINCT
                            Name
                      FROM
                            Job
                      WHERE
                            SchedTime::date = DATE(NOW())
                )
          GROUP BY
                j.Name, p.Name, j.Level`

        err := db.Select(&jobsSummary, sqlState)

        return jobsSummary, err
}
