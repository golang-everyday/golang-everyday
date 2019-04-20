package hashset

import (
	"testing"
	)

//测试Set的用法
func TestHashSet(t *testing.T)  {
	s := NewSet()
	s.Add(11)
	s.Add("aa")
	for i, k := range s.List() {
		t.Log(i, k)
	}
}
