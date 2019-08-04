package counter

import (
	"sync"
	"sync/atomic"
	"time"
)

var counters *sync.Map
var tokens map[string]struct{}
var countLen int
var maxTokenLen int

const (
	TOKEN_SELECT    = "SELECT"
	TOKEN_UPDATE    = "UPDATE"
	TOKEN_INSERT    = "INSERT"
	TOKEN_DELETE    = "DELETE"
	TOKEN_EXECPROC  = "ExecProc"
	TOKEN_SLOWQUERY = "SLOWQUERY"
)

func init() {
	counters = new(sync.Map)
	tokens = make(map[string]struct{})
	tokens[TOKEN_SELECT] = struct{}{}
	tokens[TOKEN_SELECT+"_ERROR"] = struct{}{}
	tokens[TOKEN_UPDATE] = struct{}{}
	tokens[TOKEN_UPDATE+"_ERROR"] = struct{}{}
	tokens[TOKEN_INSERT] = struct{}{}
	tokens[TOKEN_INSERT+"_ERROR"] = struct{}{}
	tokens[TOKEN_DELETE] = struct{}{}
	tokens[TOKEN_DELETE+"_ERROR"] = struct{}{}
	tokens[TOKEN_EXECPROC] = struct{}{}
	tokens[TOKEN_EXECPROC+"_ERROR"] = struct{}{}
	tokens[TOKEN_SLOWQUERY] = struct{}{}
	maxTokenLen = len(tokens)
}

func GetCounterMap() map[string]Counter {
	counterMap := make(map[string]Counter)
	counters.Range(func(key, value interface{}) bool {
		counterMap[key.(string)] = value.(Counter)
		return true
	})
	return counterMap
}

// GetCounter get counter by key
func GetCounter(key string) Counter {
	if countLen >= maxTokenLen {
		_, exists := tokens[key]
		if !exists {
			panic("Illegal key " + key)
		}
	}
	var counter Counter
	loadCounter, exists := counters.Load(key)
	if !exists {
		counter = NewCounter()
		counters.Store(key, counter)
		countLen += 1
	} else {
		counter = loadCounter.(Counter)
	}
	return counter
}

// IncHandler inc handler which will check err
func IncHandler(key string, err error, val int64) {
	if err == nil {
		GetCounter(key).Inc(val)
	} else {
		key = key + "_ERROR"
		GetCounter(key).Inc(val)
	}
}

// DecHandler dec handler which will check err
func DecHandler(key string, err error, val int64) {
	if err == nil {
		GetCounter(key).Dec(val)
	} else {
		key = key + "_ERROR"
		GetCounter(key).Dec(val)
	}
}

// Counter incremented and decremented base on int64 value.
type Counter interface {
	StartTime() time.Time
	Clear()
	Count() int64
	Dec(int64)
	Inc(int64)
}

// NewCounter constructs a new StandardCounter.
func NewCounter() Counter {
	return &StandardCounter{startTime: time.Now()}
}

// StandardCounter is the standard implementation of a Counter
type StandardCounter struct {
	count     int64
	startTime time.Time
}

func (c *StandardCounter) StartTime() time.Time {
	return c.startTime
}

// Clear sets the counter to zero.
func (c *StandardCounter) Clear() {
	atomic.StoreInt64(&c.count, 0)
}

// Count returns the current count.
func (c *StandardCounter) Count() int64 {
	return atomic.LoadInt64(&c.count)
}

// Dec decrements the counter by the given amount.
func (c *StandardCounter) Dec(i int64) {
	atomic.AddInt64(&c.count, -i)
}

// Inc increments the counter by the given amount.
func (c *StandardCounter) Inc(i int64) {
	atomic.AddInt64(&c.count, i)
}
