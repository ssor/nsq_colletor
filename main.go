package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"

	"time"

	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/ssor/config"
	"github.com/ssor/nsq_collector/parser"
)

var ()

func main() {
	config_info, err := config.LoadConfig("./conf/config.json")
	if err != nil {
		fmt.Println("[ERR] load config file err: ", err)
		return
	}
	nsqUrl := config_info.Get("url").(string)

	endPoint := config_info.Get("endpoint").(string)
	if len(endPoint) <= 0 {
		fmt.Println("[ERR] endPoint setting err: ", err)
		return
	}

	f := func() {
		statistics := pullNsqReport(nsqUrl)
		if statistics != nil {
			PushStatisticsToMonitor(statistics, endPoint, "topic_")
		}
	}
	go RunTask(f, time.Second*60)

	fmt.Println("[OK] start task")
	f() // do one time on start

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until a signal is received.
	<-c
	fmt.Println("[OK] Quit")
}

func pullNsqReport(url string) map[string]int {
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("[ERR] pull url [%s] err: %s", url, err)
		return nil
	}
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		fmt.Printf("[ERR] read response  err: %s", err)
		return nil
	}

	report, err := parser.Parse(content)
	if err != nil {
		fmt.Printf("[ERR] parse response  err: %s", err)
		return nil
	}

	return report.Summary()
}

func PushStatisticsToMonitor(statistics map[string]int, endPoint, metricPrefix string) {
	now := time.Now()
	fmt.Println("*********** result (", now.Format(time.RFC3339), "): *************")
	messages := []*FalconMessage{}
	timestamp := int(now.Unix())
	for key, count := range statistics {
		fmt.Println("topic: ", key, " -> ", count)

		msg := New_FalconMessage(endPoint, metricPrefix+key, timestamp, 60, count)
		messages = append(messages, msg)
	}

	json_bs, err := json.Marshal(messages)
	if err != nil {
		fmt.Println("[ERR] marshal err: ", err)
		spew.Dump(messages)
		return
	}

	contentReader := bytes.NewReader(json_bs)
	req, _ := http.NewRequest("POST", "http://127.0.0.1:1988/v1/push", contentReader)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("[ERR] Post data err: ", err)
		fmt.Println(string(json_bs))
		return
	}
	if resp.StatusCode == http.StatusOK {
		fmt.Println("[OK] post  success")
	} else {
		fmt.Println("[ERR] post resp: ")
		fmt.Println(string(json_bs))
		spew.Dump(resp)
	}

	fmt.Println("******************************************************************")
}

func RunTask(f func(), duration time.Duration) {
	if f == nil {
		return
	}

	ticker := time.NewTicker(duration)
	for {
		<-ticker.C
		f()
	}
}

// exists returns whether the given file or directory exists or not
func IsFileExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

/*
{
               "endpoint": 机器名（比如www等）,string
               "metric": 指标名称 ,string
               "timestamp": 时间戳, int
               "step": 60, int  (60s上传一次)
               "value": 指标, int
               "counterType": "GAUGE", string (计数器类型 增量/全量)
               "tags": "",(可默认留空字符串)
}

{"endpoint":"www.exam","metric":"conn_mongo_xsbexam_linux","timestamp":1487151349,"step":60,"value":549,"counterType":"GAUGE","tags":""}

*/

type FalconMessage struct {
	EndPoint    string `json:"endpoint"` // www.exam
	Metric      string `json:"metric"`
	Timestamp   int    `json:"timestamp"`
	Step        int    `json:"step"`
	Value       int    `json:"value"`
	CounterType string `json:"counterType"` // GAUGE
	Tags        string `json:"tags"`
}

func New_FalconMessage(endpoint, metric string, timestamp, step, value int) *FalconMessage {
	msg := &FalconMessage{
		EndPoint:    endpoint,
		Metric:      metric,
		Timestamp:   timestamp,
		Step:        step,
		Value:       value,
		CounterType: "GAUGE",
	}
	return msg
}
