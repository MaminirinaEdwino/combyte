package cmd

import (
	"slices"
	"sort"
)

type Job struct {
	ID   int
	Data []byte
}

type Result struct {
	ID      int
	Payload []byte // Contient [pIdx + length + data]
}

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

func PackBitsEncode(input []byte) []byte {
	var out []byte
	i := 0
	n := len(input)

	for i < n {
		j := i
		for j < n-1 && j-i < 127 && input[j] == input[j+1] {
			j++
		}

		if j > i {
			count := j - i + 1
			out = append(out, byte(-(count - 1)))
			out = append(out, input[i])
			i = j + 1
		} else {
			j = i
			for j < n-1 && j-i < 127 && (input[j] != input[j+1] || (j+2 < n && input[j] != input[j+2])) {
				j++
			}
			if j == n-1 {
				j = n
			} else {
				j++
			}
			count := j - i
			out = append(out, byte(count-1))
			out = append(out, input[i:j]...)
			i = j
		}
	}
	return out
}

func PackBitsDecode(input []byte) []byte {
	var out []byte
	i := 0
	for i < len(input) {
		header := int8(input[i]) // On lit le compteur comme un entier signé
		i++

		if header >= 0 {
			// Données littérales : copier (header + 1) octets
			count := int(header) + 1
			if i+count > len(input) { break }
			out = append(out, input[i:i+count]...)
			i += count
		} else if header != -128 {
			// Données répétées : répéter l'octet suivant (-header + 1) fois
			count := int(-header) + 1
			if i >= len(input) { break }
			val := input[i]
			i++
			for j := 0; j < count; j++ {
				out = append(out, val)
			}
		}
		// -128 est ignoré (NOP)
	}
	return out
}
