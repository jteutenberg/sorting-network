package sortingnetwork

import (
	"fmt"
	"math/rand"
	"slices"
	"testing"
)

func BenchmarkSortingNetwork8RegisterSIMD(b *testing.B) {
	size := 8
	b.Run(fmt.Sprintf("Size%d", size), func(b *testing.B) {
		values := make([]int32, 8*1000)
		b.ResetTimer()
		b.StopTimer()
		for i := 0; i < b.N; i++ {
			for j := range values {
				values[j] = rand.Int31()
			}
			b.StartTimer()
			for j := 0; j <= len(values)-8; j++ {
				SortingNetwork8RegisterSIMD(values[j : j+8])
			}
			b.StopTimer()
			if !slices.IsSorted(values[len(values)-8:]) {
				for _, v := range values[len(values)-8:] {
					fmt.Println(v)
				}
				fmt.Println()
				b.Fatal("sort 8 register SIMD didn't sort")
			}
		}
	})
}

func BenchmarkSortingNetwork16RegisterSIMD(b *testing.B) {
	size := 16
	b.Run(fmt.Sprintf("Size%d", size), func(b *testing.B) {
		values := make([]int32, 16*500)
		b.ResetTimer()
		b.StopTimer()
		for i := 0; i < b.N; i++ {
			for j := range values {
				values[j] = rand.Int31()
			}
			b.StartTimer()
			for j := 0; j <= len(values)-16; j++ {
				SortingNetwork16RegisterSIMD(values[j : j+16])
			}
			b.StopTimer()
			if !slices.IsSorted(values[len(values)-16:]) {
				for _, v := range values[len(values)-16:] {
					fmt.Println(v)
				}
				fmt.Println()
				b.Fatal("sort 16 register SIMD didn't sort")
			}
		}
	})
}

func BenchmarkSortingNetwork32RegisterSIMD(b *testing.B) {
	size := 32
	b.Run(fmt.Sprintf("Size%d", size), func(b *testing.B) {
		values := make([]int32, 32*250)
		b.ResetTimer()
		b.StopTimer()
		for i := 0; i < b.N; i++ {
			for j := range values {
				values[j] = rand.Int31()
			}
			b.StartTimer()
			for j := 0; j <= len(values)-32; j++ {
				SortingNetwork32RegisterSIMD(values[j : j+32])
			}
			b.StopTimer()
			if !slices.IsSorted(values[len(values)-32:]) {
				for _, v := range values[len(values)-32:] {
					fmt.Println(v)
				}
				fmt.Println()
				b.Fatal("sort 32 register SIMD didn't sort")
			}
		}
	})
}

func BenchmarkSortingNetwork64RegisterSIMD(b *testing.B) {
	size := 64
	b.Run(fmt.Sprintf("Size%d", size), func(b *testing.B) {
		values := make([]int32, 64*125)
		b.ResetTimer()
		b.StopTimer()
		for i := 0; i < b.N; i++ {
			for j := range values {
				values[j] = rand.Int31()
			}
			b.StartTimer()
			for j := 0; j <= len(values)-64; j++ {
				SortingNetwork64RegisterSIMD(values[j : j+64])
			}
			b.StopTimer()
			if !slices.IsSorted(values[len(values)-64:]) {
				for _, v := range values[len(values)-64:] {
					fmt.Println(v)
				}
				fmt.Println()
				b.Fatal("sort 64 register SIMD didn't sort")
			}
		}
	})
}

func BenchmarkSortingNetwork64BitonicSIMD(b *testing.B) {
	size := 64
	b.Run(fmt.Sprintf("Size%d", size), func(b *testing.B) {
		values := make([]int32, 64*125)
		b.ResetTimer()
		b.StopTimer()
		for i := 0; i < b.N; i++ {
			for j := range values {
				values[j] = rand.Int31()
			}
			b.StartTimer()
			for j := 0; j <= len(values)-64; j++ {
				SortingNetwork64BitonicSIMD(values[j : j+64])
			}
			b.StopTimer()
			if !slices.IsSorted(values[len(values)-64:]) {
				for _, v := range values[len(values)-64:] {
					fmt.Println(v)
				}
				fmt.Println()
				b.Fatal("sort 64 bitonic SIMD didn't sort")
			}
		}
	})
}

func BenchmarkSortingNetwork32BitonicSIMD(b *testing.B) {
	size := 32
	b.Run(fmt.Sprintf("Size%d", size), func(b *testing.B) {
		values := make([]int32, 32*250)
		b.ResetTimer()
		b.StopTimer()
		for i := 0; i < b.N; i++ {
			for j := range values {
				values[j] = rand.Int31()
			}
			b.StartTimer()
			for j := 0; j <= len(values)-32; j++ {
				SortingNetwork32BitonicSIMD(values[j : j+32])
			}
			b.StopTimer()
			if !slices.IsSorted(values[len(values)-32:]) {
				for _, v := range values[len(values)-32:] {
					fmt.Println(v)
				}
				fmt.Println()
				b.Fatal("sort 32 bitonic SIMD didn't sort")
			}
		}
	})
}

func BenchmarkSortingNetwork16BitonicSIMD(b *testing.B) {
	size := 16
	b.Run(fmt.Sprintf("Size%d", size), func(b *testing.B) {
		values := make([]int32, 16*500)
		b.ResetTimer()
		b.StopTimer()
		for i := 0; i < b.N; i++ {
			for j := range values {
				values[j] = rand.Int31()
			}
			b.StartTimer()
			for j := 0; j <= len(values)-16; j++ {
				SortingNetwork16BitonicSIMD(values[j : j+16])
			}
			b.StopTimer()
			if !slices.IsSorted(values[len(values)-16:]) {
				for _, v := range values[len(values)-16:] {
					fmt.Println(v)
				}
				fmt.Println()
				b.Fatal("sort 16 bitonic SIMD didn't sort")
			}
		}
	})
}

func BenchmarkSortingNetwork16ListedSIMD(b *testing.B) {
	size := 16
	b.Run(fmt.Sprintf("Size%d", size), func(b *testing.B) {
		values := make([]int32, 16*500)
		b.ResetTimer()
		b.StopTimer()
		for i := 0; i < b.N; i++ {
			for j := range values {
				values[j] = rand.Int31()
			}
			b.StartTimer()
			for j := 0; j <= len(values)-16; j++ {
				SortingNetwork16ListedSIMD(values[j : j+16])
			}
			b.StopTimer()
			if !slices.IsSorted(values[len(values)-16:]) {
				for _, v := range values[len(values)-16:] {
					fmt.Println(v)
				}
				fmt.Println()
				b.Fatal("sort 16 listed SIMD didn't sort")
			}
		}
	})
}

func BenchmarkSortingNetwork8SIMD(b *testing.B) {
	size := 8
	b.Run(fmt.Sprintf("Size%d", size), func(b *testing.B) {
		values := make([]int32, 8*1000)
		for i := range values {
			values[i] = rand.Int31()
		}
		b.ResetTimer()
		b.StopTimer()
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			for j := range values {
				values[j] = rand.Int31()
			}
			b.StartTimer()
			for j := 0; j <= len(values)-8; j++ {
				SortingNetwork8SIMD(values[j : j+8])
			}
			b.StopTimer()
			if !slices.IsSorted(values[len(values)-8:]) {
				for j := range values[len(values)-8:] {
					fmt.Println(j, values[j])
				}
				fmt.Println()
				b.Fatal("sort 8 SIMD didn't sort")
			}
		}
	})
}
