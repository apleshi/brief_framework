package plugin

import (
	"os"
	"github.com/Unknwon/goconfig"
	"brief_framework/config"
	"time"
	"brief_framework/util/schedule"
	"net/http"
	"encoding/json"
	"strings"
	"brief_framework/logger"
	"io/ioutil"
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
	gap, count, overThresCount int64
	conf *goconfig.ConfigFile
	matricPre, endPoint, reportUrl, contentType, reportMethod string
	maxTime, avgTime, sumTime, thresholdTime float64
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
	matricPre = conf.MustValue("report", "metric_pre", "metric.pre.def")

	gap = conf.MustInt64("report", "interval", 60)
	agentConfFile := conf.MustValue("report", "agent_config", "/usr/local/open-falcon/agent/cfg.json")
	endPoint = getHostName(agentConfFile)

	schedule.DoFuncWithTimer(doReport, time.Duration(gap) * time.Second)
}

func collectData(elapseTime time.Duration) {
	costInMs := float64(elapseTime) / float64(time.Millisecond)
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
	avgTime = 0
	sumTime = 0
	maxTime = 0
	overThresCount = 0
}

func doReport() error {
	var err error
	var timestamp, qps int64
	var bodyStr []byte
	var reportData []IndicatorDef
	var reportMetrics map[string]interface{}

	defer clearIndicator()

	if count != 0 {
		avgTime = sumTime / float64(count)
		qps = count / gap
	}
	reportMetrics = map[string]interface{} {
		".avgTime" : avgTime,
		".maxTime" : maxTime,
		".overThresCount" : overThresCount,
		".qps" : qps,
		".cnt" : count,
	}


	timestamp = time.Now().Unix()
	//endPoint, err := os.Hostname()
	//if err != nil {
	//	logger.Instance().Warn("doReport err on gethostname %s", err.Error())
	//	endPoint = "l25-248-35.lq2.autohome.cc"
	//}

	for key, val := range reportMetrics {
		reportData = append(reportData, IndicatorDef{endPoint, matricPre + key, timestamp, gap, val, "GAUGE", ""})
	}

	bodyStr, err = json.Marshal(reportData)
	if err != nil {
		logger.Instance().Error("latency report err marshal data %+v", reportData)
		return nil
	}

	logger.Instance().Debug("latency report data %s", bodyStr)

	resp, e := http.Post(reportUrl, contentType, strings.NewReader(string(bodyStr)))
	if e != nil {
		logger.Instance().Error("latency report error on post data %s, err is %s", bodyStr, e.Error())
		return nil
	}

	resp.Body.Close()

	return nil
}

func getHostName(filename string) string {
	/*
	#/usr/local/open-falcon/agent/cfg.json
	{
		"debug": false,
		"hostname":"l25-248-35.lq2.autohome.cc",
		"ip":"10.25.248.35",
		"plugin": {
			"enabled": true,
			"dir": "/usr/local/open-falcon/plugins",
			"git": "https://github.com/open-falcon/plugin.git",
			"logs": "./logs"
		},
		"heartbeat": {
			"enabled": true,
			"addr": "10.23.3.123:6030",
			"interval": 60,
			"timeout": 1000
		},
		"transfer": {
			"enabled": true,
			"addrs": [
				"10.33.3.5:8433",
			"10.20.2.45:8433",
				"10.20.2.46:8433",
				"10.20.2.47:8433"
		],
			"interval": 60,
			"timeout": 1000
		},
		"http": {
			"enabled": true,
			"listen": ":1988",
			"backdoor": false
		},
		"collector": {
			"ifacePrefix": ["eth", "em", "bond" ," vir", "p3p"]
		},
		"ignore": {
			"cpu.busy": true,
			"df.bytes.free": true,
			"df.bytes.total": true,
			"df.bytes.used": true,
			"df.bytes.used.percent": true,
			"df.inodes.total": true,
			"df.inodes.free": true,
			"df.inodes.used": true,
			"df.inodes.used.percent": true,
			"mem.memtotal": true,
			"mem.memused": true,
			"mem.memused.percent": true,
			"mem.memfree": true,
			"mem.swaptotal": true,
			"mem.swapused": true,
			"mem.swapfree": true
		}
	}
	 */
	if _ , err := os.Stat(filename); err != nil {
		logger.Instance().Warn("getHostName not find agent configfile %s, err %s", filename, err.Error())
		return "l25-248-35.lq2.autohome.cc"
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Instance().Warn("getHostName read agent configfile %s err %s.", filename, err.Error())
		return "l25-248-35.lq2.autohome.cc"
	}

	v := make(map[string]interface{})
	//读取的数据为json格式，需要进行解码
	err = json.Unmarshal(data, &v)
	if err != nil {
		logger.Instance().Warn("getHostName Unmarshal config %s err %s.", data, err.Error())
		return "l25-248-35.lq2.autohome.cc"
	}

	return v["hostname"].(string)
}