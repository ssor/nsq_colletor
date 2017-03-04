package parser

import (
	"io/ioutil"
	"testing"
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
		{"add_new_task_record", 0},
		{"taskrecords_update", 1611},
	}

	summary := nmr.Summary()

	for _, topicTestCase := range topicTestCases {
		topic := nmr.findTopic(topicTestCase.TopicName)
		if topic == nil {
			t.Fatalf("topic [%s] should exist", topicTestCase.TopicName)
		}
		if topic.totalRequeueCount() != topicTestCase.RequeueCount {
			t.Fatalf("topic [%s] has %d requeue msg, but %d expected", topicTestCase.TopicName, topic.totalRequeueCount(), topicTestCase.RequeueCount)
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
