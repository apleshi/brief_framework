package schedule

import (
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"brief_framework/logger"
)

var err_falg int32
var wg sync.WaitGroup

func init() {
	CronRun()
}

func doFunc(f func() error) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		i := 0
		for ; i < 3; i++ {
			err := f()
			if err != nil {
				logger.Instance().Error("schedule init %s, err, %v", runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(), err)
				time.Sleep(time.Millisecond * 300)
				continue
			}
			break
		}
		if i == 3 {
			logger.Instance().Error("schedule init %s, err, for 3 times", runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name())
			atomic.AddInt32(&err_falg, 1)
		}
	}()
}

func DoFuncWithTimer(f func() error, duration time.Duration) {
	go func() {
		var ticker = time.NewTicker(duration)
		for _ = range ticker.C {
			i := 0
			for ; i < 3; i++ {
				err := f()
				if err != nil {
					logger.Instance().Warn("schedule timer %s, err, %v", runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(), err)
					time.Sleep(time.Millisecond * 300)
					continue
				}
				break
			}
			if i == 3 {
				logger.Instance().Error("schedule timer %s, err, for 3 times", runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name())
			}
		}
	}()
}

