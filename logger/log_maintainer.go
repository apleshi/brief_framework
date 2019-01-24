package logger

import (
"regexp"
"time"
"path/filepath"
"os"
"brief_framework/config"
"brief_framework/util"
"brief_framework/util/schedule"
)

const DEFAULT_LOG_RETENTION_DAYS = 30
const DEFAULT_LOG_DIR = "../logs/"

func init() {
	spec := config.Instance().MustValue("schedule", "log_maintenance", "0 0 2 * * ?")
	schedule.CronAdd(spec, Clean)
}

func Clean() {
	retentionDays := config.Instance().MustInt(config.RunningMode() + ".log", "retention_days", DEFAULT_LOG_RETENTION_DAYS)

	Instance().Info("Now will delete log files which older than %d days.", retentionDays)
	DeleteLogFiles(retentionDays)
}

func DeleteLogFiles(days int) {
	var (
		logDir string
		layout = "2006-01-02"
		loc, _ = time.LoadLocation("Local")
		dateRegexp = regexp.MustCompile(`\.log\.([0-9\-]+)\.`)
	)

	logDir = config.Instance().MustValue(config.RunningMode() + ".log", "log_dir", DEFAULT_LOG_DIR)
	now := time.Now()
	err := filepath.Walk(util.GetBaseDirectory() + "/" + logDir, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		dateSlice := dateRegexp.FindStringSubmatch(path)

		if len(dateSlice) > 1 {
			fileDate := dateSlice[1]
			t, err := time.ParseInLocation(layout, fileDate, loc)
			if err != nil {
				Instance().Warn("time.Parse err = %v", err)
			}
			if t.Add(time.Hour * 24 * time.Duration(days)).Before(now) {
				Instance().Info("delete file: " + path)
				os.Remove(path)
			}
		} else {
			return nil
		}
		return nil
	})

	if err != nil {
		Instance().Warn("filepath.Walk() returned err, err = %v", err)
	}
}
