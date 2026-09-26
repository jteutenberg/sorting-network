package sortingnetwork

import "simd/archsimd"

// SortingNetwork32BitonicSIMD sorts the 32 int32 channels at the front of data
// with a bitonic sorter. The channels stay in four vectors for every layer:
// each layer has sixteen comparators, so the comparison count never drops.
// Partners sit at distances 1, 2, 4, 8, and 16. Distances 8 and 16 are
// cross-vector min/max; the shorter distances are the same in-vector shuffles
// as the size-16 sorter. The sorted values already land on channels 0..31,
// so there is no closing permute.
func SortingNetwork32BitonicSIMD(data []int32) {
	v0 := archsimd.LoadInt32x8(data)
	v1 := archsimd.LoadInt32x8(data[8:])
	v2 := archsimd.LoadInt32x8(data[16:])
	v3 := archsimd.LoadInt32x8(data[24:])

	// Block of 2. Direction alternates up, down, up, down inside every vector.
	v0 = compareDist1UDUD(v0)
	v1 = compareDist1UDUD(v1)
	v2 = compareDist1UDUD(v2)
	v3 = compareDist1UDUD(v3)

	// Block of 4. Low half of each vector sorts up, high half sorts down.
	v0 = compareDist2UUDD(v0)
	v1 = compareDist2UUDD(v1)
	v2 = compareDist2UUDD(v2)
	v3 = compareDist2UUDD(v3)
	v0 = compareDist1UUDD(v0)
	v1 = compareDist1UUDD(v1)
	v2 = compareDist1UUDD(v2)
	v3 = compareDist1UUDD(v3)

	// Block of 8. Vectors on channels 0 and 16 sort up; 8 and 24 sort down.
	v0 = compareDist4Up(v0)
	v1 = compareDist4Down(v1)
	v2 = compareDist4Up(v2)
	v3 = compareDist4Down(v3)
	v0 = compareDist2Up(v0)
	v1 = compareDist2Down(v1)
	v2 = compareDist2Up(v2)
	v3 = compareDist2Down(v3)
	v0 = compareDist1Up(v0)
	v1 = compareDist1Down(v1)
	v2 = compareDist1Up(v2)
	v3 = compareDist1Down(v3)

	// Block of 16. The low 16 channels merge up and the high 16 merge down.
	v0, v1 = compareDist8Up(v0, v1)
	v2, v3 = compareDist8Down(v2, v3)
	v0 = compareDist4Up(v0)
	v1 = compareDist4Up(v1)
	v2 = compareDist4Down(v2)
	v3 = compareDist4Down(v3)
	v0 = compareDist2Up(v0)
	v1 = compareDist2Up(v1)
	v2 = compareDist2Down(v2)
	v3 = compareDist2Down(v3)
	v0 = compareDist1Up(v0)
	v1 = compareDist1Up(v1)
	v2 = compareDist1Down(v2)
	v3 = compareDist1Down(v3)

	// Block of 32. Merge that bitonic sequence into sorted order, all up.
	v0, v2 = compareDist8Up(v0, v2)
	v1, v3 = compareDist8Up(v1, v3)
	v0, v1 = compareDist8Up(v0, v1)
	v2, v3 = compareDist8Up(v2, v3)
	v0 = compareDist4Up(v0)
	v1 = compareDist4Up(v1)
	v2 = compareDist4Up(v2)
	v3 = compareDist4Up(v3)
	v0 = compareDist2Up(v0)
	v1 = compareDist2Up(v1)
	v2 = compareDist2Up(v2)
	v3 = compareDist2Up(v3)
	v0 = compareDist1Up(v0)
	v1 = compareDist1Up(v1)
	v2 = compareDist1Up(v2)
	v3 = compareDist1Up(v3)

	v0.Store(data)
	v1.Store(data[8:])
	v2.Store(data[16:])
	v3.Store(data[24:])
}

// SortingNetwork32RegisterSIMD sorts the 32 int32 channels at the front of data.
// It is the size-32 bitonic merge with a different first half:
//  1. Each lane is sorted by min/max between the four vectors, with no shuffles.
//  2. Those eight sorted columns are packed into four vectors, two runs of 4 each.
//  3. Each vector is merged to a sorted run of 8, odd runs are reversed, and the
//     existing 16/32 bitonic merges finish the sort.
func SortingNetwork32RegisterSIMD(data []int32) {
	v0 := archsimd.LoadInt32x8(data)
	v1 := archsimd.LoadInt32x8(data[8:])
	v2 := archsimd.LoadInt32x8(data[16:])
	v3 := archsimd.LoadInt32x8(data[24:])

	v0, v1 = compareDist8Up(v0, v1)
	v2, v3 = compareDist8Down(v2, v3)
	v0, v2 = compareDist8Up(v0, v2)
	v1, v3 = compareDist8Up(v1, v3)
	v0, v1 = compareDist8Up(v0, v1)
	v2, v3 = compareDist8Up(v2, v3)

	v0, v1, v2, v3 = transpose4x8(v0, v1, v2, v3)
	v0 = mergeSorted4sUp(v0)
	v1 = mergeSorted4sUp(v1)
	v2 = mergeSorted4sUp(v2)
	v3 = mergeSorted4sUp(v3)
	v1 = v1.Permute(laneReverse8)
	v3 = v3.Permute(laneReverse8)

	// Block of 16. The low 16 channels merge up and the high 16 merge down.
	v0, v1 = compareDist8Up(v0, v1)
	v2, v3 = compareDist8Down(v2, v3)
	v0 = compareDist4Up(v0)
	v1 = compareDist4Up(v1)
	v2 = compareDist4Down(v2)
	v3 = compareDist4Down(v3)
	v0 = compareDist2Up(v0)
	v1 = compareDist2Up(v1)
	v2 = compareDist2Down(v2)
	v3 = compareDist2Down(v3)
	v0 = compareDist1Up(v0)
	v1 = compareDist1Up(v1)
	v2 = compareDist1Down(v2)
	v3 = compareDist1Down(v3)

	// Block of 32. Merge that bitonic sequence into sorted order, all up.
	v0, v2 = compareDist8Up(v0, v2)
	v1, v3 = compareDist8Up(v1, v3)
	v0, v1 = compareDist8Up(v0, v1)
	v2, v3 = compareDist8Up(v2, v3)
	v0 = compareDist4Up(v0)
	v1 = compareDist4Up(v1)
	v2 = compareDist4Up(v2)
	v3 = compareDist4Up(v3)
	v0 = compareDist2Up(v0)
	v1 = compareDist2Up(v1)
	v2 = compareDist2Up(v2)
	v3 = compareDist2Up(v3)
	v0 = compareDist1Up(v0)
	v1 = compareDist1Up(v1)
	v2 = compareDist1Up(v2)
	v3 = compareDist1Up(v3)

	v0.Store(data)
	v1.Store(data[8:])
	v2.Store(data[16:])
	v3.Store(data[24:])
}

// transpose4x8 packs eight sorted columns of 4 into four vectors.
// Column j is (v0[j], v1[j], v2[j], v3[j]). The result holds two columns per vector.
func transpose4x8(v0, v1, v2, v3 archsimd.Int32x8) (archsimd.Int32x8, archsimd.Int32x8, archsimd.Int32x8, archsimd.Int32x8) {
	t0 := v0.InterleaveLoGrouped(v1)
	t1 := v0.InterleaveHiGrouped(v1)
	t2 := v2.InterleaveLoGrouped(v3)
	t3 := v2.InterleaveHiGrouped(v3)
	return t0.ConcatPermute128Scalars(0, 2, t2).Permute(packSorted4s),
		t1.ConcatPermute128Scalars(0, 2, t3).Permute(packSorted4s),
		t0.ConcatPermute128Scalars(1, 3, t2).Permute(packSorted4s),
		t1.ConcatPermute128Scalars(1, 3, t3).Permute(packSorted4s)
}
