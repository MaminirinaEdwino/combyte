package main

import (
	"flag"
	"fmt"
)

const(
	compressText = "This is the option for compressing a text file"
	extractText = "This is the option for extracting a file .combyte"
)

func main() {
	fmt.Println("Combyte CLI")
	
	compress := flag.Bool("compress", false, compressText)
	c := flag.Bool("c", false, compressText)

	extract := flag.Bool("extract", false, extractText)
	e := flag.Bool("e", false, extractText)

	filename := flag.String("filename", "", "THe file that you want to extract or compress")
	flag.Parse()
	switch  {
	case *compress || *c:
		if *filename == "" {
			break
		}
	case *extract || *e:
		if *filename == "" {
			break
		}
	default : 
		fmt.Println("Help documentation")
	}
}