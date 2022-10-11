package util

import "fmt"

type Pair[T1, T2 any] struct {
	First  T1
	Second T2
}

func NewPair[T1, T2 any](fst T1, snd T2) Pair[T1, T2] {
	return Pair[T1, T2]{
		First:  fst,
		Second: snd,
	}
}

func (p Pair[T1, T2]) String() string {
	return fmt.Sprintf("(%v, %v)", p.First, p.Second)
}
