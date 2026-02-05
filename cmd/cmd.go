package cmd

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
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
	count := make([]int, 256)
	for _, b := range bwt {
		count[b]++
	}
	sum := 0
	for i, c := range count {
		count[i] = sum
		sum += c
	}
	T := make([]int, n)
	for i, b := range bwt {
		T[count[b]] = i
		count[b]++
	}
	result := make([]byte, n)
	var curr int
	fmt.Println("test : ", len(T), primaryIndex, curr)
	curr = T[primaryIndex]

	for i := range n {
		result[i] = bwt[curr]
		curr = T[curr]
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
		header := int8(input[i])
		i++

		if header >= 0 {
			count := int(header) + 1
			if i+count > len(input) {
				break
			}
			out = append(out, input[i:i+count]...)
			i += count
		} else if header != -128 {
			count := int(-header) + 1
			if i >= len(input) {
				break
			}
			val := input[i]
			i++
			for range count {
				out = append(out, val)
			}
		}
	}
	return out
}

func CompressFile(r io.Reader, w io.Writer, compressionLevel int) {
	numWorkers := runtime.NumCPU()
	jobs := make(chan Job, numWorkers)
	results := make(chan Result, numWorkers)
	var wg sync.WaitGroup
	blockSize := compressionLevel * 1024
	for range numWorkers {
		wg.Go(func() {
			for job := range jobs {
				bwt, pIdx := BWT(job.Data)
				rle := PackBitsEncode(bwt)
				buf := new(bytes.Buffer)
				binary.Write(buf, binary.LittleEndian, int32(pIdx))
				binary.Write(buf, binary.LittleEndian, int32(len(rle)))
				buf.Write(rle)
				results <- Result{ID: job.ID, Payload: buf.Bytes()}
			}
		})
	}

	go func() {

		counter := 0
		for {
			buf := make([]byte, blockSize)
			n, err := r.Read(buf)
			if n > 0 {
				jobs <- Job{ID: counter, Data: buf[:n]}
				counter++
			}
			if err != nil {
				break
			}
		}
		close(jobs)
		wg.Wait()
		close(results)
	}()

	var totalByteTreated int64
	start := time.Now()

	pending := make(map[int][]byte)
	nextID := 0
	for res := range results {
		pending[res.ID] = res.Payload
		for {
			if data, ok := pending[nextID]; ok {
				w.Write(data)
				delete(pending, nextID)
				nextID++
				totalByteTreated += int64(blockSize)
				elapsed := time.Since(start).Seconds()
				if elapsed > 0 {
					mbps := float64(totalByteTreated) / float64(1024*1024) / elapsed
					fmt.Printf("\rCompressed Blocks : %d \tSpeed : %.2f Mo/s\t\tElapsed Time : %s", nextID, mbps, time.Since(start))
				}
			} else {
				break
			}
		}
	}
	fmt.Println()
}

func DecompressFile(file *os.File, destFile string) {

	extractedFile, _ := os.Create(destFile)
	defer extractedFile.Close()

	reader := bufio.NewReader(file)
	for {
		var pIdx int32
		var length int32

		err := binary.Read(reader, binary.LittleEndian, &pIdx)
		if err == io.EOF {
			break
		}

		err = binary.Read(reader, binary.LittleEndian, &length)

		rleData := make([]byte, length)
		_, err = io.ReadFull(reader, rleData)

		bwtData := PackBitsDecode(rleData)
		realData := IBWT(bwtData, int(pIdx))
		extractedFile.WriteString(string(realData))
	}
}

func Compress(filename string, compressionLevel int) {
	source, err := os.Open(filename)
	if err.Error() == fmt.Sprintf("open %s: no such file or directory", filename) {
		fmt.Printf("Fichier %s introuvable\n", filename)
	} else {
		dest, _ := os.Create(filename + ".combyte")
		defer dest.Close()

		reader := bufio.NewReader(source)
		writer := bufio.NewWriter(dest)
		defer writer.Flush()

		CompressFile(reader, writer, compressionLevel)
	}
}

func Extract(filename string) {
	source, err := os.Open(filename)
	if err.Error() == fmt.Sprintf("open %s: no such file or directory", filename) {
		fmt.Printf("Fichier %s introuvable\n", filename)
	} else {
		DecompressFile(source, strings.Replace(filename, ".combyte", "", 1))
	}
}
