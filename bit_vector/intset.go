package main

import (
	"bytes"
	"fmt"
)

type IntSet struct {
	words []uint64
}

func (s *IntSet) Has(x int) bool {
	word, bit := x/64, uint(x%64)
	return word <= len(s.words) && s.words[word]&(1<<bit) != 0
}

func (s *IntSet) Add(x int) {
	word, bit := x/64, uint(x%64)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

func (s *IntSet) AddAll(values ...int) {
	for _, v := range values {
		s.Add(v)
	}
}

func (s *IntSet) Remove(x int) {
	word, bit := x/64, uint(x%64)
	s.words[word] &^= 1 << bit
}

func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

func (s *IntSet) IntersectWith(t *IntSet) {
	var length, i int
	if len(s.words) < len(t.words) {
		length = len(s.words)
	} else {
		length = len(t.words)
	}
	for i = 0; i < length; i++ {
		s.words[i] &= t.words[i]
	}
	if i < len(s.words)-1 {
		s.words = s.words[:i+1]
	}
}

func (s *IntSet) DifferenceWith(t *IntSet) {
	slength := len(s.words)
	tlength := len(t.words)
	for i := 0; i < slength; i++ {
		if i >= tlength {
			break
		}
		s.words[i] &^= t.words[i]
	}
}

func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < 64; j++ {
			if word&(1<<j) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", 64*i+j)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}

func popcount(x uint64) int {
	count := 0
	for x > 0 {
		count++
		x &= x - 1
	}
	return count
}

func (s *IntSet) Len() int {
	var count int
	for _, word := range s.words {
		count += popcount(word)
	}
	return count
}

func (s *IntSet) Clear() {
	s.words = []uint64{}
}

func (s *IntSet) Copy() *IntSet {
	var result IntSet
	result.words = make([]uint64, len(s.words))
	copy(result.words, s.words)
	return &result
}

func main() {
	var x, y IntSet
	x.Add(1)
	x.Add(144)
	x.Add(9)
	fmt.Println(x.String())

	y.Add(9)
	y.Add(42)
	fmt.Println(y.String())

	//x.UnionWith(&y)
	//fmt.Println(x.String())

	//fmt.Println(x.Has(9), x.Has(123))
	//x.Remove(144)
	//x.Remove(1)
	//x.Clear()
	//x.IntersectWith(&y)
	//t := x.Copy()

	x.DifferenceWith(&y)
	fmt.Println(x.String())
	//fmt.Println(t.String())
	fmt.Println(x.Len())
}
