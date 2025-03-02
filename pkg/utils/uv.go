package utils

import (
	"sort"
	"sync"
)

var localCounter uvCounter

// UVStat 用于存储用户ID和对应的访问次数
type UVStat struct {
	BookId int64
	Count  int64
}

// uvCounter 结构体包含并发安全的UV计数器
type uvCounter struct {
	mu     sync.RWMutex
	counts map[int64]int64
}

func init() {
	localCounter = uvCounter{
		counts: make(map[int64]int64),
	}
}

// AddUv 增加指定用户的访问计数
func AddUv(bookId int64) {
	localCounter.mu.Lock()
	defer localCounter.mu.Unlock()
	localCounter.counts[bookId]++
}

// GetTopK 获取分页的TopN访问统计
func GetTopK(page, pageSize int) []UVStat {
	localCounter.mu.RLock()
	defer localCounter.mu.RUnlock()

	// 将map转换为切片
	stats := make([]UVStat, 0, len(localCounter.counts))
	for userID, count := range localCounter.counts {
		stats = append(stats, UVStat{userID, count})
	}

	// 排序
	sort.Slice(stats, func(i, j int) bool {
		if stats[i].Count == stats[j].Count {
			return stats[i].BookId < stats[j].BookId
		}
		return stats[i].Count > stats[j].Count
	})

	// 处理分页参数
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	start := (page - 1) * pageSize
	end := start + pageSize

	// 处理越界情况
	if start >= len(stats) {
		return []UVStat{}
	}
	if end > len(stats) {
		end = len(stats)
	}

	return stats[start:end]
}
