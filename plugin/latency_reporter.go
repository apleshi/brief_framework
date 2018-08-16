package plugin

import (
	"os"
	"github.com/Unknwon/goconfig"
	"brief_framework/config"
	"time"
	"brief_framework/util/schedule"
	"brief_framework/util"
	"net/http"
	"github.com/gin-gonic/gin/json"
	"strings"
	"brief_framework/logger"
)

type indicatorDef struct {
	endpoint string
	metric string
	timestamp int64
	step int64
	value interface{}
	counterType string
	tags string
}

var (
	maxTime, sumTime, thresholdTime float64
	gap, count, overThresCount int64
	conf *goconfig.ConfigFile
	reportUrl, contentType, reportMethod string
)

func init() {
	var err error
	configFile := config.Instance().MustValue("latency", "report_conf", "./conf/monitor.ini")

	_, err = os.Stat(configFile)
	if err != nil {
		conf, err = goconfig.LoadConfigFile("../" + configFile)
	} else {
		conf, err = goconfig.LoadConfigFile(configFile)
	}

	if err != nil {
		panic(err)
	}

	reportUrl = conf.MustValue("report", "address", "http://127.0.0.1:1988/v1/push")
	contentType = conf.MustValue("report", "content-type", "application/json")
	reportMethod = conf.MustValue("report", "method", "POST")

	gap := conf.MustInt("report", "interval", 60)
	schedule.DoFuncWithTimer(doReport, time.Duration(gap) * time.Second)
}

func collectData(elapseTime float64) {
	if elapseTime > maxTime {
		maxTime = elapseTime
	}

	if elapseTime > thresholdTime {
		overThresCount++
	}

	sumTime += elapseTime
	count++
}

func doReport() error {
	var err error
	var timestamp int64
	var endPoint string
	var bodyStr []byte
	var reportData []indicatorDef
	var reportMetrics =
		map[string]interface{}{
		"avgTime" : sumTime / (float64(count)),
		"maxTime" : maxTime,
		"overThresCount" : overThresCount,
		"qps" : count / gap,
	}

	timestamp = time.Now().Unix()
	endPoint = "onlineFeature-" + util.GetIntranetIp()

	for key, val := range reportMetrics {
		reportData = append(reportData, indicatorDef{endPoint, key, timestamp, gap, val, "GAUGE", ""})
	}

	bodyStr, err = json.Marshal(reportData)
	if err != nil {
		logger.Instance().Error("latency report err marshal data %+v", reportData)
		return nil
	}

	resp, e := http.Post(reportUrl, contentType, strings.NewReader(string(bodyStr)))
	if e != nil {
		logger.Instance().Error("latency report error on post data %s", bodyStr)
		return nil
	}

	resp.Body.Close()
	return nil
}