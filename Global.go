package database

import "github.com/devfeel/database/internal/counter"

// ShowStateData return database stat data with map[string]int64
func ShowStateData() map[string]int64 {
	data := make(map[string]int64)
	counters := counter.GetCounterMap()
	for k, v := range counters {
		data[k] = v.Count()
	}
	return data
}
