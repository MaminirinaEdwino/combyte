package cmd

import (
	"slices"
	"sort"
)


func BWT(input []byte) ([]byte, int) {
	n := len(input)
	indices := make([]int, n)
	for i := range indices {
		indices[i] = i
	}
	sort.Slice(indices, func(i, j int) bool {
		a, b := indices[i], indices[j]
		for k := range n {
			byteA := input[(a+k)%n]
			byteB := input[(b+k)%n]
			if byteA != byteB {
				return byteA < byteB
			}
		}
		return false
	})

	result := make([]byte, n)
	primaryIndex := 0
	for i, idx := range indices {
		if idx == 0 {
			primaryIndex = i
		}
		result[i] = input[(idx+n-1)%n]
	}
	return result, primaryIndex
}

func IBWT(bwt []byte, primaryIndex int) []byte {
	n := len(bwt)
	if n == 0 {
		return nil
	}

	firstCol := make([]byte, n)
	copy(firstCol, bwt)
	slices.Sort(firstCol)

	T := make([]int, n)
	count := make(map[byte][]int)
	for i, char := range bwt {
		count[char] = append(count[char], i)
	}

	curr := 0
	for _, char := range firstCol {
		T[curr] = count[char][0]
		count[char] = count[char][1:]
		curr++
	}

	result := make([]byte, n)
	last := primaryIndex
	for i := n - 1; i >= 0; i-- {
		result[i] = bwt[last]
		for j, val := range T {
			if val == last {
				last = j
				break
			}
		}
	}
	return result
}
