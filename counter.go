/*
* @Author:  老杨
* @Email:   xcapp1314@gmail.com
* @Date:    2023/7/13 07:01
* @Explain: ...
 */

package counter

import (
	"sync"
	"time"
)

type Counter struct {
	accessMap   map[time.Time]int64
	mutex       sync.Mutex
	ticker      *time.Ticker
	done        chan struct{}
	cleanupTime time.Duration
}

// NewCounter 创建一个统计技术器
func NewCounter(cleanupTime ...time.Duration) *Counter {
	counter := &Counter{
		accessMap:   make(map[time.Time]int64),
		ticker:      time.NewTicker(time.Minute),
		done:        make(chan struct{}),
		cleanupTime: time.Hour, // 默认保留1小时数据
	}
	if len(cleanupTime) > 0 {
		counter.cleanupTime = cleanupTime[0]
		// 启动定时器，每隔一段时间清理旧数据
		go counter.cleanupOldData()
	}

	return counter
}

func (c *Counter) Close() {
	_ = recover() //忽略关闭错误
	close(c.done)
	c.ticker.Stop()
}

func (c *Counter) Add(count int64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.accessMap[time.Now()] += count
}

func (c *Counter) Count(duration time.Duration) int64 {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// 获取当前时间
	now := time.Now()

	// 计算最早的时间点
	earliest := now.Add(-duration)

	// 统计特定时间内访问量
	var count int64
	for t, cnt := range c.accessMap {
		if t.After(earliest) && t.Before(now) {
			count += cnt
		}
	}
	return count
}

func (c *Counter) cleanupOldData() {
	for {
		select {
		case <-c.ticker.C:
			c.mutex.Lock()

			// 获取当前时间
			now := time.Now()

			// 计算最早的时间点
			earliest := now.Add(-c.cleanupTime)

			// 删除早于 cleanupTime 的数据
			for t := range c.accessMap {
				if t.Before(earliest) {
					delete(c.accessMap, t)
				}
			}

			c.mutex.Unlock()
		case <-c.done:
			return
		}
	}
}
