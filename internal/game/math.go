package game

import "math"

func Round(term float64) int {
	return int(math.Round(term))
}

func Ptr[T any](v T) *T { return &v }
