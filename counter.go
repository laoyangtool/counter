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

// Counter 计数器结构体
type Counter struct {
	accessList  []*AccessRecord
	rwMutex     sync.RWMutex
	ticker      *time.Ticker
	done        chan struct{}
	cleanupTime time.Duration
}

// AccessRecord 访问记录结构体
type AccessRecord struct {
	Time  time.Time
	Count int64
}

// NewCounter 创建一个统计计数器
func NewCounter(cleanupTime ...time.Duration) *Counter {
	counter_ := &Counter{
		accessList:  make([]*AccessRecord, 0),
		ticker:      time.NewTicker(time.Minute),
		done:        make(chan struct{}),
		cleanupTime: 2 * time.Hour, // 默认保留2小时数据
	}
	if len(cleanupTime) > 0 {
		counter_.cleanupTime = cleanupTime[0]
		go counter_.cleanupOldData()
	}
	return counter_
}

// Add 添加访问记录
func (c *Counter) Add(count int64) {
	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()
	c.accessList = append(c.accessList, &AccessRecord{Time: time.Now(), Count: count})
}

// Count 统计指定时间范围内的访问量
func (c *Counter) Count(duration time.Duration) int64 {
	c.rwMutex.RLock()
	defer c.rwMutex.RUnlock()

	now := time.Now()
	earliest := now.Add(-duration)
	var count int64

	// 从头部开始查找，直到遇到第一个不满足条件的记录
	for _, record := range c.accessList {
		if record.Time.Before(earliest) {
			count += record.Count
		} else {
			break
		}
	}

	// 清理已统计过的记录
	c.accessList = c.accessList[len(c.accessList)-len(c.accessList)+1:]

	return count
}

// cleanupOldData 定期清理旧数据
func (c *Counter) cleanupOldData() {
	for {
		select {
		case <-c.ticker.C:
			c.rwMutex.Lock()

			// 获取当前时间
			now := time.Now()

			// 计算最早的时间点
			earliest := now.Add(-c.cleanupTime)

			// 从头部开始删除满足条件的记录
			for i, record := range c.accessList {
				if record.Time.Before(earliest) {
					c.accessList = c.accessList[i+1:]
				} else {
					break
				}
			}

			c.rwMutex.Unlock()
		case <-c.done:
			return
		}
	}
}

// Close 关闭计数器
func (c *Counter) Close() {
	defer func() {
		_ = recover() // 忽略关闭错误
	}()
	close(c.done)
	c.ticker.Stop()
}
