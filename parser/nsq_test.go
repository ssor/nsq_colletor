package parser

import (
	"io/ioutil"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

const (
	dataPath = "testdata/nsq.json"
)

func TestParseLog(t *testing.T) {
	raw, err := ioutil.ReadFile(dataPath)
	if err != nil {
		t.Fatal(err)
	}

	nmr, err := Parse(raw)
	if err != nil {
		t.Fatal(err)
	}

	if len(nmr.Topics) != 2 {
		t.Fatal("topic count should be 2")
	}

	topicTestCases := []struct {
		TopicName    string
		RequeueCount int
	}{
		{"add_new_task_record_ch_depth", 1},
		{"add_new_task_record_ch_deferred", 2},
		{"taskrecords_update_ch_depth", 10},
		{"taskrecords_update_ch_deferred", 100},
	}

	summary := nmr.Summary()
	spew.Dump(summary)
	for _, topicTestCase := range topicTestCases {
		_, exists := summary[topicTestCase.TopicName]
		if exists == false {
			t.Fatalf("topic [%s] should exist", topicTestCase.TopicName)
		}

		if count, exists := summary[topicTestCase.TopicName]; exists == false {
			t.Fatalf("topic [%s] should exist in summary", topicTestCase.TopicName)
		} else {
			if count != topicTestCase.RequeueCount {
				t.Fatalf("topic summary count is %d, and %d is expected", count, topicTestCase.RequeueCount)
			}
		}
	}
}
