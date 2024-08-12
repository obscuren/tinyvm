package vm

import "fmt"

type stack struct {
	data []uint32
}

func newStack() stack {
	return stack{
		data: make([]uint32, 0),
	}
}

func (s *stack) push(d uint32) {
	s.data = append(s.data, d)
}

func (s *stack) pop() (d uint32) {
	d = s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return
}

func (s stack) len() int {
	return len(s.data)
}

func (s *stack) dump() {
	fmt.Println("stack:")
	for i, d := range s.data {
		fmt.Printf("%d: %v\n", i, d)
	}
}
