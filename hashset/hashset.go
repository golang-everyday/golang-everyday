package hashset

import (
	"sync"
)

//golang 模拟java中的hashset的用法

type inter interface{}

//hashset结构体
type Set struct {
	m map[inter]bool
	sync.RWMutex
}

//创建set方法
func NewSet() *Set {
	return &Set{
		m: map[inter]bool{},
	}
}

//添加元素
func (s *Set) Add(item inter) {
	s.Lock()
	defer s.Unlock()
	s.m[item] = true
}

//删除元素
func (s *Set) Remove(item inter) {
	s.Lock()
	s.Unlock()
	delete(s.m, item)
}

//元素是否在set集合中
func (s *Set) Has(item inter) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.m[item]
	return ok
}

//返回set集合的长度
func (s *Set) Len() int {
	return len(s.List())
}

//清空set集合
func (s *Set) Clear() {
	s.Lock()
	defer s.Unlock()
	s.m = map[inter]bool{}
}

//判断集合是否为空
func (s *Set) IsEmpty() bool {
	if s.Len() == 0 {
		return true
	}
	return false
}

func (s *Set) List() []inter {
	s.RLock()
	defer s.RUnlock()
	list := []inter{}
	for item := range s.m {
		list = append(list, item)
	}
	return list
}
