package util

import "golang.org/x/exp/constraints"

type dot struct {
	r string
	n int32
}

type DVV struct {
	d  dot
	vv map[string]int32
}

func (dvv DVV) Lt(odvv DVV) bool {
	d := dvv.d
	return d.n <= odvv.vv[d.r]
}

func sync(D1, D2 []DVV) []DVV {
	var result []DVV
	for _, x := range D1 {
		include := true
		for _, y := range D2 {
			if x.Lt(y) {
				include = false
				break
			}
		}
		if include {
			result = append(result, x)
		}
	}
	for _, x := range D2 {
		include := true
		for _, y := range D1 {
			if x.Lt(y) {
				include = false
				break
			}
		}
		if include {
			result = append(result, x)
		}
	}
	return result
}

func (dvv DVV) ids() []string {
	var result []string
	result = append(result, dvv.d.r)
	for r := range dvv.vv {
		result = append(result, r)
	}
	return result
}

func ids(dvvs []DVV) []string {
	var result []string
	for _, dvv := range dvvs {
		result = append(result, dvv.ids()...)
	}
	return result
}

func (dvv DVV) ceil(r string) int32 {
	if dvv.d.r == r {
		return max(dvv.d.n, dvv.vv[r])
	}
	return dvv.vv[r]
}

func ceil(dvvs []DVV, r string) int32 {
	m := int32(0)
	for _, dvv := range dvvs {
		m = max(m, dvv.ceil(r))
	}
	return m
}

func update(S []DVV, S_r []DVV, r string) DVV {
	result := DVV{
		d: dot{
			r: r,
			n: ceil(S_r, r) + 1,
		},
		vv: make(map[string]int32),
	}
	for _, i := range ids(S) {
		result.vv[i] = ceil(S, i)
	}
	return result
}

func max[T constraints.Ordered](t1, t2 T) T {
	if t1 > t2 {
		return t1
	}
	return t2
}
