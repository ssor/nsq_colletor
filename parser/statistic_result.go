package parser

type StatisticResult map[string]int

func (sr StatisticResult) Add(more StatisticResult) StatisticResult {
	if more == nil {
		return sr
	}

	for k, v := range more {
		if value, exists := sr[k]; exists {
			sr[k] = value + v
		} else {
			sr[k] = v
		}
	}
	return sr
}
