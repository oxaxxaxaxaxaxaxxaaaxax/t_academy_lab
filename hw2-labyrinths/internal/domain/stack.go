package domain

import "errors"

type Stack struct {
	St []Cell
}

func (s Stack) Push(c Cell) Stack {
	s.St = append(s.St, c)
	return s
}

func (s Stack) Pop() (Stack, Cell, error) {
	stackLen := len(s.St)
	if stackLen == 0 {
		return Stack{}, Cell{}, errors.New("stack is empty")
	}

	c := s.St[stackLen-1]
	s.St = s.St[:stackLen-1]
	return s, c, nil
}

func (s Stack) Empty() bool {
	return len(s.St) == 0
}
