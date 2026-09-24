package domain

type Node[T any] struct {
	Priority int
	Value    T
}
