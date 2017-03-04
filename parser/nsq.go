package parser

import "encoding/json"

func Parse(raw []byte) (*NsqMessageReport, error) {
	var nmr NsqMessageReport
	err := json.Unmarshal(raw, &nmr)
	if err != nil {
		return nil, err
	}
	return &nmr, nil
}

type NsqMessageReport struct {
	Topics TopicArray
}

func (report *NsqMessageReport) Summary() map[string]int {
	summary := make(map[string]int)
	for _, topic := range report.Topics {
		summary[topic.Name] = topic.totalRequeueCount()
	}
	return summary
}

func (report *NsqMessageReport) findTopic(name string) *Topic {
	return report.Topics.find(name)
}

type TopicArray []*Topic

func (topics TopicArray) find(name string) *Topic {
	for _, topic := range topics {
		if topic.Name == name {
			return topic
		}
	}
	return nil
}

type Topic struct {
	Name     string `json:"topic_name"`
	Channels ChannelArray
}

func (topic *Topic) totalRequeueCount() int {
	total := 0
	if topic.Channels == nil || len(topic.Channels) <= 0 {
		return total
	}

	for _, channel := range topic.Channels {
		total = total + channel.RequeueCount
	}
	return total
}

type ChannelArray []*Channel

type Channel struct {
	Node         string `json:"node"`
	HostName     string `json:"hostname"`
	TopicName    string `json:"topic_name"`
	ChannelName  string `json:"channel_name"`
	RequeueCount int    `json:"requeue_count"`
	// Clients      ClientArray `json:"clients"`
}

// type ClientArray []*Client

// type Client struct {
// }
