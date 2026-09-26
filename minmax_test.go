package sortingnetwork

import (
	"fmt"
	"math/rand"
	"slices"
	"testing"
)

func BenchmarkSortingNetwork8MinMax(b *testing.B) {
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
				SortingNetwork8MinMax(values[j : j+8])
			}
			b.StopTimer()
			if !slices.IsSorted(values[len(values)-8:]) {
				for j := range values[len(values)-8:] {
					fmt.Println(j, values[j])
				}
				fmt.Println()
				b.Fatal("sort 8 didn't sort")
			}
		}
	})
}

func BenchmarkSortingNetwork64MinMax(b *testing.B) {
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
				SortingNetwork64MinMax(values[j : j+64])
			}
			b.StopTimer()
			if !slices.IsSorted(values[len(values)-64:]) {
				for j := range values[len(values)-64:] {
					fmt.Println(j, values[j])
				}
				fmt.Println()
				b.Fatal("sort 64 didn't sort")
			}
		}
	})
}

func BenchmarkSortingNetwork32MinMax(b *testing.B) {
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
				SortingNetwork32MinMax(values[j : j+32])
			}
			b.StopTimer()
			if !slices.IsSorted(values[len(values)-32:]) {
				for j := range values[len(values)-32:] {
					fmt.Println(j, values[j])
				}
				fmt.Println()
				b.Fatal("sort 32 didn't sort")
			}
		}
	})
}

func BenchmarkSortingNetwork16MinMax(b *testing.B) {
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
				SortingNetwork16MinMax(values[j : j+16])
			}
			b.StopTimer()
			if !slices.IsSorted(values[len(values)-16:]) {
				for j := range values[len(values)-16:] {
					fmt.Println(j, values[j])
				}
				fmt.Println()
				b.Fatal("sort 16 didn't sort")
			}
		}
	})
}
