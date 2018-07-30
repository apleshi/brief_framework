package schedule

import (
	"github.com/robfig/cron"
	"brief_framework/config"
	"brief_framework/logger"
	"brief_framework/util"
)

func CronRun() {
	c := cron.New()
	spec, err := config.Instance().GetValue("schedule", "log_maintenance")
	if err != nil {
		logger.Instance().Warn("GetValue schedule, log_maintenance error, err = %v", err)
	} else {
		c.AddFunc(spec, func() {
			logger.Instance().Info("Timed task for delete log files start.")
			util.Clean()
		})
	}

	c.Start()
}