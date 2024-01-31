/*
* @Author:  老杨
* @Email:   xcapp1314@gmail.com
* @Date:    2024/1/31 22:20:30 星期三
* @Explain: ...
 */

package counter

import (
	"sync"
	"time"
)

var (
	_counters = make(map[string]*Counter)
	_mu       sync.Mutex
)

func counter(key string) *Counter {
	_mu.Lock()
	defer _mu.Unlock()
	if _, ok := _counters[key]; !ok {
		_counters[key] = NewCounter()
	}
	return _counters[key]
}

func Add(key string, count int64) {
	counter(key).Add(count)
}

func Close(key string) {
	_mu.Lock()
	defer _mu.Unlock()
	if _counter, ok := _counters[key]; ok {
		_counter.Close()
		delete(_counters, key)
	}
}

func Count(key string, duration time.Duration) int64 {
	return counter(key).Count(duration)
}
