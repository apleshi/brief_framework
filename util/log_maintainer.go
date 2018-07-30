package util

import (
	"regexp"
	"time"
	"path/filepath"
	"os"
	"brief_framework/config"
	"brief_framework/logger"
)

const DEFAULT_LOG_RETENTION_DAYS = 30
const DEFAULT_LOG_DIR = "./logs/"

func Clean() {
	retentionDays := config.Instance().MustInt(config.RunningMode(), "log_retention_days", DEFAULT_LOG_RETENTION_DAYS)

	logger.Instance().Info("Now will delete log files which older than %d days.", retentionDays)
	DeleteLogFiles(retentionDays)
}

func DeleteLogFiles(days int) {
	var (
		logDir string
		layout = "2006-01-02"
		loc, _ = time.LoadLocation("Local")
		dateRegexp = regexp.MustCompile(`\.log\.([0-9\-]+)\.`)
	)

	logDir = config.Instance().MustValue(config.RunningMode(), "log_dir", DEFAULT_LOG_DIR)
	now := time.Now()
	err := filepath.Walk(logDir, func(path string, f os.FileInfo, err error) error {
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
				logger.Instance().Warn("time.Parse err = %v", err)
			}
			if t.Add(time.Hour * 24 * time.Duration(days)).Before(now) {
				logger.Instance().Info("delete file: " + path)
				os.Remove(path)
			}
		} else {
			return nil
		}
		return nil
	})

	if err != nil {
		logger.Instance().Warn("filepath.Walk() returned err, err = %v", err)
	}
}
