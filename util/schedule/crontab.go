package schedule

import (
	"github.com/robfig/cron"
)

var c *cron.Cron
func initCron() {
	if c == nil {
		c = cron.New()
		c.Start()
	}
}

//func CronRun() {
//	c := cron.New()
//	spec, err := config.Instance().GetValue("schedule", "log_maintenance")
//	if err != nil {
//		logger.Instance().Warn("GetValue schedule, log_maintenance error, err = %v", err)
//	} else {
//		c.AddFunc(spec, func() {
//			logger.Instance().Info("Timed task for delete log files start.")
//			logger.Clean()
//		})
//	}
//
//	c.Start()
//}

func CronAdd(spec string, f func()) {
	if c == nil {
		c = cron.New()
		c.Start()
	}

	c.AddFunc(spec, func() {
		//logger.Instance().Info("Timed task for delete log files start.")
		f()
	})
}