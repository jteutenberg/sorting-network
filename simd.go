package sortingnetwork

import (
	"simd/archsimd"
)

var pairUpOddEven = archsimd.LoadUint32x8Array(&[8]uint32{2, 3, 0, 1, 6, 7, 4, 5})

func SortingNetwork8SIMD(data []int32) {
	a := archsimd.LoadInt32x8(data)

	// The 8 integers referenced by index are 0,1,2,3,4,5,6,7
	// Lower-case m(0,1) indicates the minimum of values at indices 0 and 1
	// Upper-case M(0,1) indicates the maximum of values at indices 0 and 1
	// showing both the permutation of values and which operator applies to each pair.
	// This is always shown using the original input indices, with nested operations. It can be listed with one index per line.

	// The layer is also shown as the pre- and post- permutation lined up and also with its operator, e.g.
	// 0,1,2,3,4,5,6,7
	// 1,0,3,2,5,4,7,6
	// m M m M m M m M
	// which always uses the indices from the previous layer at the top

	// Layer 1

	// output of layer: m(0,1), M(0,1), m(2,3), M(2,3), m(4,5), M(4,5), m(6,7), M(6,7)

	// 0,1,2,3,4,5,6,7
	// 1,0,3,2,5,4,7,6
	// m M m M m M m M
	// this is giving four sub-graphs, each with an ordering of two nodes

	b := a.PermuteScalarsGrouped(1, 0, 3, 2)
	c := a.Min(b)
	a = a.Max(b)
	a = c.ConcatPermuteScalarsGrouped(0, 4, 2, 6, a)

	// Layer 2

	// we now want to apply the min/max to pairs that were both min after
	// applying the previous layer, or both max after the previous layer
	// e.g. indices 0,2 (both min) and1,3 (both max) etc.
	// output of layer (using initial indices):
	//  0. m(m(0,1), m(2,3))  the minimum of the lower subgraph of 4
	//  1. m(M(0,1), M(2,3))  an element of the lower subgraph
	//  2. M(m(0,1), m(2,3))  an element of the lower subgraph
	//  3. M(M(0,1), M(2,3))  the maximum of the lower subgraph of 4
	// and similarly for an upper subgraph of size 4

	// 0,1,2,3,4,5,6,7
	// 2,3,0,1,6,7,4,5
	// m m M M m m M M

	b = a.PermuteScalarsGrouped(2, 3, 0, 1) // permutes to 2,3,0,1,6,7,4,5
	c = a.Min(b)
	a = a.Max(b)
	a = c.ConcatPermuteScalarsGrouped(0, 1, 6, 7, a)

	// this is now in two sub-graphs 0,1,2,3 and 4,5,6,7

	// Layer 3

	// we now want to get the min of the two minimums, max of the two maximums
	// for the other four nodes, we want the two that came out min in the last layer:
	// 0,1,2,3,4,5,6,7
	// 4,5,6,7,0,1,2,3
	// m m m m M M M M

	b = a.ConcatPermute128Scalars(3, 0, a) // don't need the concat, but no other suitable op
	c = a.Min(b)
	a = a.Max(b)
	a = c.ConcatPermute128Scalars(0, 3, a)

	// we now have 0 holding the overall min, 7 holding the overall max
	// 1 holds the min of one remaining half (of three nodes)
	// 6 holds the max of the other remaining half
	// more completely, the lower half is (using original indices):
	// 0. m(m(m(0,1), m(2,3), m(m(4,5), m(6,7)) )
	// 1. m(m(M(0,1), M(2,3)), m(M(4,5), M(6,7)) )
	// 2. m(M(m(0,1), m(2,3)), M(m(4,5), m(6,7))  )
	// 3. m(M(M(0,1), M(2,3)), M(M(4,5), M(6,7))  )
	// and the remaining 4 being the maximums with the same arguments

	// Layer 4

	// this and the next layer work on two sets of three nodes
	// 1, 2, 3 and 4, 5, 6
	// We want to compare two mins from the last layer with each other
	// and two maxes from the last layer with each other.

	// 2. m(m(M(m(0,1), m(2,3)), M(m(4,5), m(6,7))  ), m(M(M(0,1), M(2,3)), M(M(4,5), M(6,7))  ) )
	// 3. M(m(M(m(0,1), m(2,3)), M(m(4,5), m(6,7))  ), m(M(M(0,1), M(2,3)), M(M(4,5), M(6,7))  ) )
	// 4. m(M(M(m(0,1), m(2,3)), M(m(4,5), m(6,7))  ), M(M(M(0,1), M(2,3)), M(M(4,5), M(6,7))  ) )
	// 5. M(M(M(m(0,1), m(2,3)), M(m(4,5), m(6,7))  ), M(M(M(0,1), M(2,3)), M(M(4,5), M(6,7))  ) )

	// 0,1,2,3,4,5,6,7
	// 1,0,3,2,5,4,7,6
	// m M m M m M m M Same as layer 1 though it leaves the outer 4 where they started

	/*
		However, this is 4 SIMD ops to do two min and two max operations.
		It is probably time to drop out to regular operations

		b := a.PermuteScalarsGrouped(1,0,3,2)
		c := a.Min(b)
		a = a.Max(b)
		a = c.ConcatPermuteScalarsGrouped(0, 4, 2, 6, a)
	*/

	x := c.GetLo()
	y := c.GetHi()

	b1 := x.GetElem(1)
	c1 := x.GetElem(2)
	d1 := x.GetElem(3)
	e1 := y.GetElem(0)
	f1 := y.GetElem(1)
	g1 := y.GetElem(2)

	c1, e1 = min(c1, e1), max(c1, e1)
	d1, f1 = min(d1, f1), max(d1, f1)

	// Layer 5
	b1, e1 = min(b1, e1), max(b1, e1)
	d1, g1 = min(d1, g1), max(d1, g1)

	// Layer 6
	b1, c1 = min(b1, c1), max(b1, c1)
	d1, e1 = min(d1, e1), max(d1, e1)
	f1, g1 = min(f1, g1), max(f1, g1)

	x = x.SetElem(1, b1).SetElem(2, c1).SetElem(3, d1)
	y = y.SetElem(0, e1).SetElem(1, f1).SetElem(2, g1)
	a = a.SetLo(x).SetHi(y)
	a.Store(data)
}

// SortingNetwork8RegisterSIMD sorts the 8 int32 channels at the front of data.
// One vector is already a single row, so there is no cross-vector transpose.
// The same run merge used by the larger register sorts builds sorted pairs,
// then fours, then one run of 8, and stays in the vector instead of dropping
// to scalar for the last layers.
func SortingNetwork8RegisterSIMD(data []int32) {
	a := archsimd.LoadInt32x8(data)
	a = compareDist1Up(a)
	a = sortedPairsTo8(a)
	a.Store(data)
}
