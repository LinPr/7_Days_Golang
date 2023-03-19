package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// hash maps bytes to uint32
type Hash func(data []byte) uint32

// map consistants all hashed keys
type Map struct {
	hash     Hash
	replicas int

	// 一个[]int + map 的组合数据结构当做 ordered map 来用
	keys    []int
	hashMap map[int]string
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	// 默认哈希函数为 crc32.ChecksumIEEE
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}

	// 将[]int + map 的组合数据结构当o rdered map 来用
	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	hash := int(m.hash([]byte(key)))
	idx := sort.Search(len(m.keys), func(i int) bool {
		// 在数组中找到第一个条件成立的元素的下标
		return m.keys[i] >= hash
	})

	//
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
