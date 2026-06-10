package daemon

// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"net/http"

	"pkg.re/essentialkaos/ek.v12/log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// /////////////////////////////////////////////////////////////////////////////

type baculaMetrics struct {
	LatestJobFiles       *prometheus.Desc
	LatestJobBytes       *prometheus.Desc
	LatestJobSchedTime   *prometheus.Desc
	LatestJobStartTime   *prometheus.Desc
	LatestJobEndTime     *prometheus.Desc
	SummaryJobTotalFiles *prometheus.Desc
	SummaryJobTotalBytes *prometheus.Desc
	StoredBytes          *prometheus.Desc
	StoredFiles          *prometheus.Desc
	PoolAvailableTapes   *prometheus.Desc
}

// /////////////////////////////////////////////////////////////////////////////

// startHTTPServer start HTTP server
func startHTTPServer(ip, port, endpoint string) error {
	addr := ip + ":" + port

	log.Info("HTTP server is started on %s", addr)

	http.Handle(endpoint, promhttp.Handler())

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})

	return http.ListenAndServe(addr, nil)
}

// /////////////////////////////////////////////////////////////////////////////

// baculaCollector returns baculaMetrics struct for Prometheus
func baculaCollector() *baculaMetrics {
	return &baculaMetrics{
		LatestJobFiles: prometheus.NewDesc("bacula_latest_job_files_total",
			"Total files saved for server during latest backup for client combined",
			[]string{"name", "pool", "jobid", "level", "status"}, nil,
		),
		LatestJobBytes: prometheus.NewDesc("bacula_latest_job_bytes_total",
			"Total bytes saved for server during latest backup for client combined",
			[]string{"name", "pool", "jobid", "level", "status"}, nil,
		),
		LatestJobSchedTime: prometheus.NewDesc("bacula_latest_job_sched_time",
			"Timestamp when the latest job was scheduled",
			[]string{"name", "pool", "jobid", "level", "status"}, nil,
		),
		LatestJobStartTime: prometheus.NewDesc("bacula_latest_job_start_time",
			"Timestamp when the latest job was started",
			[]string{"name", "pool", "jobid", "level", "status"}, nil,
		),
		LatestJobEndTime: prometheus.NewDesc("bacula_latest_job_end_time",
			"Timestamp when the latest job was ended",
			[]string{"name", "pool", "jobid", "level", "status"}, nil,
		),
		SummaryJobTotalFiles: prometheus.NewDesc("bacula_summary_job_files_total",
			"Total files saved for server during all backups for client combined",
			[]string{"name", "pool", "level"}, nil,
		),
		SummaryJobTotalBytes: prometheus.NewDesc("bacula_summary_job_bytes_total",
			"Total bytes saved for server during all backups for client combined",
			[]string{"name", "pool", "level"}, nil,
		),
		StoredBytes: prometheus.NewDesc("bacula_stored_bytes_total",
				"Total bytes currently stored per job name and pool (excludes recycled/purged volumes)",
				[]string{"name", "pool"}, nil,
		),
		StoredFiles: prometheus.NewDesc("bacula_stored_files_total",
				"Total files currently stored per job name and pool (excludes recycled/purged volumes)",
				[]string{"name", "pool"}, nil,
		),
		PoolAvailableTapes: prometheus.NewDesc("bacula_pool_available_tapes",
				"Total number of available tapes (enabled volumes with Append, Purged or Recycle status) per pool",
				[]string{"pool"}, nil,
		),
	}
}

// /////////////////////////////////////////////////////////////////////////////

// Describe implements Describe() method using by the Prometheus registry
// when describing metrics
func (collector *baculaMetrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.LatestJobFiles
	ch <- collector.LatestJobBytes
	ch <- collector.LatestJobSchedTime
	ch <- collector.LatestJobStartTime
	ch <- collector.LatestJobEndTime
	ch <- collector.SummaryJobTotalFiles
	ch <- collector.SummaryJobTotalBytes
	ch <- collector.StoredBytes
	ch <- collector.StoredFiles
	ch <- collector.PoolAvailableTapes
}

// Collect implements Collect() method using by the Prometheus registry
// when collecting metrics
func (collector *baculaMetrics) Collect(ch chan<- prometheus.Metric) {
	latestJobs, err := env.DB.GetLatestJobs()

	if err != nil {
		log.Crit(err.Error())
		return
	}

	for _, job := range latestJobs {
		ch <- prometheus.MustNewConstMetric(
			collector.LatestJobFiles,
			prometheus.GaugeValue,
			float64(job.JobFiles),
			job.Name,
			job.Pool,
			job.JobId,
			job.Level,
			job.Status,
		)
		ch <- prometheus.MustNewConstMetric(
			collector.LatestJobBytes,
			prometheus.GaugeValue,
			float64(job.JobBytes),
			job.Name,
			job.Pool,
			job.JobId,
			job.Level,
			job.Status,
		)
		ch <- prometheus.MustNewConstMetric(
			collector.LatestJobSchedTime,
			prometheus.CounterValue,
			float64(job.SchedTime),
			job.Name,
			job.Pool,
			job.JobId,
			job.Level,
			job.Status,
		)
		ch <- prometheus.MustNewConstMetric(
			collector.LatestJobStartTime,
			prometheus.CounterValue,
			float64(job.StartTime),
			job.Name,
			job.Pool,
			job.JobId,
			job.Level,
			job.Status,
		)
		ch <- prometheus.MustNewConstMetric(
			collector.LatestJobEndTime,
			prometheus.CounterValue,
			float64(job.EndTime),
			job.Name,
			job.Pool,
			job.JobId,
			job.Level,
			job.Status,
		)
	}

	jobsSummary, err := env.DB.GetJobsSummary()

	if err != nil {
		log.Crit(err.Error())
		return
	}

	for _, job := range jobsSummary {
		ch <- prometheus.MustNewConstMetric(
			collector.SummaryJobTotalFiles,
			prometheus.GaugeValue,
			float64(job.TotalJobFiles),
			job.Name,
			job.Pool,
			job.Level,
		)
		ch <- prometheus.MustNewConstMetric(
			collector.SummaryJobTotalBytes,
			prometheus.GaugeValue,
			float64(job.TotalJobBytes),
			job.Name,
			job.Pool,
			job.Level,
		)
	}

	storedData, err := env.DB.GetStoredData()

	if err != nil {
		log.Crit(err.Error())
		return
	}

	for _, item := range storedData {
		ch <- prometheus.MustNewConstMetric(
			collector.StoredBytes,
			prometheus.GaugeValue,
			float64(item.StoredBytes),
			item.Name,
			item.Pool,
		)
		ch <- prometheus.MustNewConstMetric(
			collector.StoredFiles,
			prometheus.GaugeValue,
			float64(item.StoredFiles),
			item.Name,
			item.Pool,
		)
	}

	availableTapes, err := env.DB.GetAvailableTapes()

	if err != nil {
		log.Crit(err.Error())
		return
	}

	for _, item := range availableTapes {
		ch <- prometheus.MustNewConstMetric(
			collector.PoolAvailableTapes,
			prometheus.GaugeValue,
			float64(item.AvailableTapes),
			item.Pool,
		)
	}
}

// /////////////////////////////////////////////////////////////////////////////
