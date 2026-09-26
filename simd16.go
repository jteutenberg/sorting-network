package sortingnetwork

import "simd/archsimd"

// SortingNetwork16BitonicSIMD sorts the 16 int32 channels at the front of data
// with a bitonic sorter. The channels stay in two vectors for every layer:
// each layer has eight comparators, so the comparison count never drops.
// Partners sit at distances 1, 2, 4, and 8. Distance 8 is a cross-vector
// min/max; the shorter distances are in-vector shuffles. The sorted values
// already land on channels 0..15, so there is no closing permute.
func SortingNetwork16BitonicSIMD(data []int32) {
	lo := archsimd.LoadInt32x8(data)
	hi := archsimd.LoadInt32x8(data[8:])

	// Block of 2. Direction alternates up, down, up, down on each vector.
	lo = compareDist1UDUD(lo)
	hi = compareDist1UDUD(hi)

	// Block of 4. Low half of each vector sorts up, high half sorts down.
	lo = compareDist2UUDD(lo)
	hi = compareDist2UUDD(hi)
	lo = compareDist1UUDD(lo)
	hi = compareDist1UUDD(hi)

	// Block of 8. lo sorts up, hi sorts down, which makes one bitonic sequence.
	lo = compareDist4Up(lo)
	hi = compareDist4Down(hi)
	lo = compareDist2Up(lo)
	hi = compareDist2Down(hi)
	lo = compareDist1Up(lo)
	hi = compareDist1Down(hi)

	// Block of 16. Merge that bitonic sequence into sorted order, all up.
	lo, hi = compareDist8Up(lo, hi)
	lo = compareDist4Up(lo)
	hi = compareDist4Up(hi)
	lo = compareDist2Up(lo)
	hi = compareDist2Up(hi)
	lo = compareDist1Up(lo)
	hi = compareDist1Up(hi)

	lo.Store(data)
	hi.Store(data[8:])
}

// SortingNetwork16RegisterSIMD sorts the 16 int32 channels at the front of data.
// It is the size-16 bitonic merge with a different first half:
//  1. Corresponding lanes of the two vectors are ordered by one min/max.
//  2. Those eight sorted pairs are transposed into the two vectors.
//  3. Each vector is merged up to a sorted run of 8, the high run is reversed,
//     and the existing 16-wide bitonic merge finishes the sort.
func SortingNetwork16RegisterSIMD(data []int32) {
	lo := archsimd.LoadInt32x8(data)
	hi := archsimd.LoadInt32x8(data[8:])
	lo, hi = compareDist8Up(lo, hi)

	pairsLo := lo.InterleaveLoGrouped(hi)
	pairsHi := lo.InterleaveHiGrouped(hi)
	lo = pairsLo.ConcatPermute128Scalars(0, 2, pairsHi)
	hi = pairsLo.ConcatPermute128Scalars(1, 3, pairsHi)
	lo = sortedPairsTo8(lo)
	hi = sortedPairsTo8(hi)
	hi = hi.Permute(laneReverse8)

	lo, hi = compareDist8Up(lo, hi)
	lo = compareDist4Up(lo)
	hi = compareDist4Up(hi)
	lo = compareDist2Up(lo)
	hi = compareDist2Up(hi)
	lo = compareDist1Up(lo)
	hi = compareDist1Up(hi)

	lo.Store(data)
	hi.Store(data[8:])
}

// compareDist1Up compares adjacent channels and writes min, max, min, max, ...
func compareDist1Up(a archsimd.Int32x8) archsimd.Int32x8 {
	b := a.PermuteScalarsGrouped(1, 0, 3, 2)
	c := a.Min(b)
	d := a.Max(b)
	return c.ConcatPermuteScalarsGrouped(0, 4, 2, 6, d)
}

// compareDist1Down compares adjacent channels and writes max, min, max, min, ...
func compareDist1Down(a archsimd.Int32x8) archsimd.Int32x8 {
	b := a.PermuteScalarsGrouped(1, 0, 3, 2)
	c := a.Min(b)
	d := a.Max(b)
	return c.ConcatPermuteScalarsGrouped(4, 0, 6, 2, d)
}

// compareDist1UDUD compares adjacent channels with directions up, down, up, down.
func compareDist1UDUD(a archsimd.Int32x8) archsimd.Int32x8 {
	b := a.PermuteScalarsGrouped(1, 0, 3, 2)
	c := a.Min(b)
	d := a.Max(b)
	return c.ConcatPermuteScalarsGrouped(0, 4, 6, 2, d)
}

// compareDist1UUDD compares adjacent channels with directions up, up, down, down.
func compareDist1UUDD(a archsimd.Int32x8) archsimd.Int32x8 {
	b := a.PermuteScalarsGrouped(1, 0, 3, 2)
	c := a.Min(b)
	d := a.Max(b)
	// Lanes 0,2,5,7 keep the min. The two 128-bit halves want opposite blends,
	// which a grouped permute cannot express.
	return c.IfElse(maskDist1UUDD, d)
}

// compareDist2Up compares channels two apart and writes min, min, max, max in each half.
func compareDist2Up(a archsimd.Int32x8) archsimd.Int32x8 {
	b := a.PermuteScalarsGrouped(2, 3, 0, 1)
	c := a.Min(b)
	d := a.Max(b)
	return c.ConcatPermuteScalarsGrouped(0, 1, 6, 7, d)
}

// compareDist2Down compares channels two apart and writes max, max, min, min in each half.
func compareDist2Down(a archsimd.Int32x8) archsimd.Int32x8 {
	b := a.PermuteScalarsGrouped(2, 3, 0, 1)
	c := a.Min(b)
	d := a.Max(b)
	return c.ConcatPermuteScalarsGrouped(4, 5, 2, 3, d)
}

// compareDist2UUDD compares channels two apart, up in the low half and down in the high half.
func compareDist2UUDD(a archsimd.Int32x8) archsimd.Int32x8 {
	b := a.PermuteScalarsGrouped(2, 3, 0, 1)
	c := a.Min(b)
	d := a.Max(b)
	return c.IfElse(maskDist2UUDD, d)
}

// compareDist4Up compares each low-half channel with the high-half channel four
// places away and leaves the mins in the low half.
func compareDist4Up(a archsimd.Int32x8) archsimd.Int32x8 {
	b := a.ConcatPermute128Scalars(3, 0, a)
	c := a.Min(b)
	d := a.Max(b)
	return c.ConcatPermute128Scalars(0, 3, d)
}

// compareDist4Down compares across the 128-bit halves and leaves the maxes in the low half.
func compareDist4Down(a archsimd.Int32x8) archsimd.Int32x8 {
	b := a.ConcatPermute128Scalars(3, 0, a)
	c := a.Min(b)
	d := a.Max(b)
	return d.ConcatPermute128Scalars(0, 3, c)
}

// compareDist8Up compares corresponding lanes of two vectors. Mins return in lo.
// Adjacent vectors are channel distance 8; vectors further apart are distance 16 or 32.
func compareDist8Up(lo, hi archsimd.Int32x8) (archsimd.Int32x8, archsimd.Int32x8) {
	return lo.Min(hi), lo.Max(hi)
}

// compareDist8Down compares corresponding lanes of two vectors. Maxes return in lo.
func compareDist8Down(lo, hi archsimd.Int32x8) (archsimd.Int32x8, archsimd.Int32x8) {
	return lo.Max(hi), lo.Min(hi)
}

// reverseHigh4 reverses lanes 4..7 and leaves 0..3 in place.
var reverseHigh4 = permIndex([8]int32{0, 1, 2, 3, 7, 6, 5, 4})

// packSorted4s gathers two interleaved pairs into consecutive sorted fours:
// [a0,b0,a1,b1, c0,d0,c1,d1] becomes [a0,b0,c0,d0, a1,b1,c1,d1].
var packSorted4s = permIndex([8]int32{0, 1, 4, 5, 2, 3, 6, 7})

// mergeSorted4sUp merges two ascending runs of 4 that already sit in one vector.
func mergeSorted4sUp(a archsimd.Int32x8) archsimd.Int32x8 {
	a = a.Permute(reverseHigh4)
	a = compareDist4Up(a)
	a = compareDist2Up(a)
	return compareDist1Up(a)
}

// sortedPairsTo8 merges four ascending pairs into one ascending run of 8.
func sortedPairsTo8(a archsimd.Int32x8) archsimd.Int32x8 {
	a = a.PermuteScalarsGrouped(0, 1, 3, 2)
	a = compareDist2Up(a)
	a = compareDist1Up(a)
	return mergeSorted4sUp(a)
}

// maskDist1UUDD keeps mins on lanes 0, 2, 5, 7.
var maskDist1UUDD = minMask([8]int32{1, 0, 1, 0, 0, 1, 0, 1})

// maskDist2UUDD keeps mins on lanes 0, 1, 6, 7.
var maskDist2UUDD = minMask([8]int32{1, 1, 0, 0, 0, 0, 1, 1})

func minMask(keepMin [8]int32) archsimd.Mask32x8 {
	return archsimd.LoadInt32x8Array(&keepMin).Equal(archsimd.BroadcastInt32x8(1))
}

// listedPairs is the size-16 listed network, in layer order. Each layer's
// eight comparators are disjoint, so the layer runs as one vector step.
// Min lands on the lower channel. Later layers compare fewer than eight pairs
// and are scalar in SortingNetwork16ListedSIMD.
var listedPairs = [4][8][2]uint8{
	{{0, 13}, {1, 12}, {2, 15}, {3, 14}, {4, 8}, {5, 6}, {7, 11}, {9, 10}},
	{{0, 5}, {1, 7}, {2, 9}, {3, 4}, {6, 13}, {8, 14}, {10, 15}, {11, 12}},
	{{0, 1}, {2, 3}, {4, 5}, {6, 8}, {7, 9}, {10, 11}, {12, 13}, {14, 15}},
	{{0, 2}, {1, 3}, {4, 10}, {5, 11}, {6, 7}, {8, 9}, {12, 14}, {13, 15}},
}

var listedLayers [4]listedLayer

func init() {
	for i, pairs := range listedPairs {
		listedLayers[i] = newListedLayer(pairs)
	}
}

// SortingNetwork16ListedSIMD sorts the 16 int32 channels at the front of data
// with the listed network. The first four layers cover every channel, so each
// one is a permute into partner lanes, a min/max, and a permute back. From the
// fifth layer onward a layer no longer covers every channel, and those
// comparators run as scalar min/max, as in the size-8 network.
func SortingNetwork16ListedSIMD(data []int32) {
	lo := archsimd.LoadInt32x8(data)
	hi := archsimd.LoadInt32x8(data[8:])
	lo, hi = listedLayers[0].apply(lo, hi)
	lo, hi = listedLayers[1].apply(lo, hi)
	lo, hi = listedLayers[2].apply(lo, hi)
	lo, hi = listedLayers[3].apply(lo, hi)

	u0 := lo.GetLo()
	u1 := lo.GetHi()
	u2 := hi.GetLo()
	u3 := hi.GetHi()
	a, b, c, d := u0.GetElem(0), u0.GetElem(1), u0.GetElem(2), u0.GetElem(3)
	e, f, g, h := u1.GetElem(0), u1.GetElem(1), u1.GetElem(2), u1.GetElem(3)
	i, j, k, l := u2.GetElem(0), u2.GetElem(1), u2.GetElem(2), u2.GetElem(3)
	m, n, o, p := u3.GetElem(0), u3.GetElem(1), u3.GetElem(2), u3.GetElem(3)

	b, c = min(b, c), max(b, c)
	d, m = min(d, m), max(d, m)
	e, g = min(e, g), max(e, g)
	f, h = min(f, h), max(f, h)
	i, k = min(i, k), max(i, k)
	j, l = min(j, l), max(j, l)
	n, o = min(n, o), max(n, o)

	b, e = min(b, e), max(b, e)
	c, g = min(c, g), max(c, g)
	f, i = min(f, i), max(f, i)
	h, k = min(h, k), max(h, k)
	j, n = min(j, n), max(j, n)
	l, o = min(l, o), max(l, o)

	c, e = min(c, e), max(c, e)
	d, g = min(d, g), max(d, g)
	j, m = min(j, m), max(j, m)
	l, n = min(l, n), max(l, n)

	d, f = min(d, f), max(d, f)
	g, i = min(g, i), max(g, i)
	h, j = min(h, j), max(h, j)
	k, m = min(k, m), max(k, m)

	d, e = min(d, e), max(d, e)
	f, g = min(f, g), max(f, g)
	h, i = min(h, i), max(h, i)
	j, k = min(j, k), max(j, k)
	l, m = min(l, m), max(l, m)

	g, h = min(g, h), max(g, h)
	i, j = min(i, j), max(i, j)

	data[0], data[1], data[2], data[3], data[4], data[5], data[6], data[7] = a, b, c, d, e, f, g, h
	data[8], data[9], data[10], data[11], data[12], data[13], data[14], data[15] = i, j, k, l, m, n, o, p
}

// listedLayer aligns one perfect matching of the 16 channels, compares, and
// writes each min back to the lower channel.
type listedLayer struct {
	left, right listedGather
	toLo, toHi  listedBlend
}

func (s *listedLayer) apply(lo, hi archsimd.Int32x8) (archsimd.Int32x8, archsimd.Int32x8) {
	left := s.left.load(lo, hi)
	right := s.right.load(lo, hi)
	mins := left.Min(right)
	maxs := left.Max(right)
	return s.toLo.mix(mins, maxs), s.toHi.mix(mins, maxs)
}

type listedGather struct {
	allLo, allHi bool
	idx          archsimd.Uint32x8
	loIdx, hiIdx archsimd.Uint32x8
	takeLo       archsimd.Mask32x8
}

func (g *listedGather) load(lo, hi archsimd.Int32x8) archsimd.Int32x8 {
	if g.allLo {
		return lo.Permute(g.idx)
	}
	if g.allHi {
		return hi.Permute(g.idx)
	}
	return lo.Permute(g.loIdx).IfElse(g.takeLo, hi.Permute(g.hiIdx))
}

type listedBlend struct {
	minIdx, maxIdx archsimd.Uint32x8
	takeMin        archsimd.Mask32x8
}

func (b *listedBlend) mix(mins, maxs archsimd.Int32x8) archsimd.Int32x8 {
	return mins.Permute(b.minIdx).IfElse(b.takeMin, maxs.Permute(b.maxIdx))
}

func newListedLayer(pairs [8][2]uint8) listedLayer {
	var seen [16]bool
	var leftCh, rightCh [8]uint8
	for p, pair := range pairs {
		if pair[0] >= pair[1] || pair[1] > 15 || seen[pair[0]] || seen[pair[1]] {
			panic("listed network layer has a repeated or unordered channel")
		}
		seen[pair[0]], seen[pair[1]] = true, true
		leftCh[p], rightCh[p] = pair[0], pair[1]
	}
	for _, ok := range seen {
		if !ok {
			panic("listed network layer does not cover every channel")
		}
	}
	return listedLayer{
		left:  newListedGather(leftCh),
		right: newListedGather(rightCh),
		toLo:  newListedBlend(pairs, 0),
		toHi:  newListedBlend(pairs, 8),
	}
}

func newListedGather(ch [8]uint8) listedGather {
	var loI, hiI, sel [8]int32
	allLo, allHi := true, true
	for i, c := range ch {
		if c < 8 {
			allHi = false
			loI[i] = int32(c)
			sel[i] = 1
		} else {
			allLo = false
			hiI[i] = int32(c - 8)
		}
	}
	g := listedGather{
		loIdx:  permIndex(loI),
		hiIdx:  permIndex(hiI),
		takeLo: minMask(sel),
	}
	if allLo {
		g.allLo = true
		g.idx = g.loIdx
	} else if allHi {
		g.allHi = true
		g.idx = g.hiIdx
	}
	return g
}

func newListedBlend(pairs [8][2]uint8, base uint8) listedBlend {
	var minI, maxI, takeMin [8]int32
	for lane := uint8(0); lane < 8; lane++ {
		ch := lane + base
		found := false
		for p, pair := range pairs {
			switch ch {
			case pair[0]:
				minI[lane] = int32(p)
				takeMin[lane] = 1
				found = true
			case pair[1]:
				maxI[lane] = int32(p)
				found = true
			}
		}
		if !found {
			panic("listed network blend missed a channel")
		}
	}
	return listedBlend{
		minIdx:  permIndex(minI),
		maxIdx:  permIndex(maxI),
		takeMin: minMask(takeMin),
	}
}

func permIndex(idx [8]int32) archsimd.Uint32x8 {
	return archsimd.LoadInt32x8Array(&idx).ToBits()
}
