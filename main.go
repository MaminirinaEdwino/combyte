package main

import (
	"flag"
	"fmt"

	"github.com/MaminirinaEdwino/combyte/cmd"
)

const(
	compressText = "This is the option for compressing a text file"
	extractText = "This is the option for extracting a file .combyte"
	compressionLevelText = "This option is for specifying the compression level that you want (it won't affect the extraction if you're worrie dabout it)"
)

func main() {
	fmt.Println("Combyte CLI")
	
	compress := flag.Bool("compress", false, compressText)
	c := flag.Bool("c", false, compressText)

	extract := flag.Bool("extract", false, extractText)
	e := flag.Bool("e", false, extractText)

	filename := flag.String("filename", "", "THe file that you want to extract or compress")
	compressionLevel := flag.Int("level", 3, compressionLevelText)
	flag.Parse()
	switch  {
	case *compress || *c:
		if *filename == "" {
			break
		}
		cmd.Compress(*filename, *compressionLevel)
	case *extract || *e:
		if *filename == "" {
			break
		}
		cmd.Extract(*filename)
	default : 
		fmt.Println("Help documentation")
	}
}