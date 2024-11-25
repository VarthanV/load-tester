package tester

import "golang.org/x/exp/constraints"

type Number interface {
	constraints.Integer | constraints.Float
}

func sum[T Number](val []T) T {
	sum := T(0)
	for _, v := range val {
		sum += v
	}
	return sum
}
