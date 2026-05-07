package game

import (
	"fmt"
)

type Queue[T any] []T

func MakeQueue[T any]() Queue[T] {
	return make(Queue[T], 0)
}

func (q *Queue[T]) Enqueue(value T) {
	*q = append(*q, value)
}

func (q *Queue[T]) Dequeue() (T, error) {
	if len(*q) == 0 {
		var zero T
		return zero, fmt.Errorf("queue is empty")
	}

	value := (*q)[0]
	var zero T
	(*q)[0] = zero
	*q = (*q)[1:]
	return value, nil
}
func (q Queue[T]) Peek() (T, error) {
	if len(q) == 0 {
		var zero T
		return zero, fmt.Errorf("queue is empty")
	}

	value := (q)[0]
	return value, nil
}
