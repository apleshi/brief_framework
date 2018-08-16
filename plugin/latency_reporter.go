package plugin

import (
	"os"
	"github.com/Unknwon/goconfig"
	"brief_framework/config"
	"time"
	"brief_framework/util/schedule"
	"brief_framework/util"
	"net/http"
	"encoding/json"
	"strings"
	"brief_framework/logger"
)

type IndicatorDef struct {
	Endpoint string			`json:"endpoint"`
	Metric string			`json:"metric"`
	Timestamp int64			`json:"timestamp"`
	Step int64				`json:"step"`
	Value interface{}		`json:"value"`
	CounterType string		`json:"counterType"`
	Tags string				`json:"tags"`
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
	thresholdTime = conf.MustFloat64("report", "response_time_thres", 1000)

	gap = conf.MustInt64("report", "interval", 60)
	schedule.DoFuncWithTimer(doReport, time.Duration(gap) * time.Second)
}

func collectData(elapseTime time.Duration) {
	costInMs := float64(elapseTime / time.Millisecond)
	if costInMs > maxTime {
		maxTime = costInMs
	}

	if costInMs > thresholdTime {
		overThresCount++
	}

	sumTime += costInMs
	count++

	logger.Instance().Debug("collectData one data elapseTime %.3f, count %d", costInMs, count)
}

func clearIndicator() {
	count = 0
	sumTime = 0
	overThresCount = 0
}

func doReport() error {
	var err error
	var timestamp int64
	var endPoint string
	var bodyStr []byte
	var reportData []IndicatorDef
	var reportMetrics map[string]interface{}

	defer clearIndicator()

	reportMetrics = map[string]interface{} {
		"avgTime" : sumTime / (float64(count)),
		"maxTime" : maxTime,
		"overThresCount" : overThresCount,
		"qps" : count / gap,
	}

	timestamp = time.Now().Unix()
	endPoint = "onlineFeature-" + util.GetIntranetIp()

	for key, val := range reportMetrics {
		reportData = append(reportData, IndicatorDef{endPoint, key, timestamp, gap, val, "GAUGE", ""})
	}

	bodyStr, err = json.Marshal(reportData)
	if err != nil {
		logger.Instance().Error("latency report err marshal data %+v", reportData)
		return nil
	}

	logger.Instance().Debug("latency report data %s", bodyStr)

	resp, e := http.Post(reportUrl, contentType, strings.NewReader(string(bodyStr)))
	if e != nil {
		logger.Instance().Error("latency report error on post data %s", bodyStr)
		return nil
	}

	resp.Body.Close()

	return nil
}