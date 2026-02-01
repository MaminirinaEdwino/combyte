package main

import (
	"flag"
	"fmt"

	colortext "github.com/MaminirinaEdwino/colorText"
	colortext "github.com/MaminirinaEdwino/colorText/colorText"
	"github.com/MaminirinaEdwino/combyte/cmd"
	"github.com/common-nighthawk/go-figure"
)

const(
	compressText = "This is the option for compressing a text file"
	extractText = "This is the option for extracting a file .combyte"
	compressionLevelText = "This option is for specifying the compression level that you want (it won't affect the extraction if you're worrie dabout it)"
)

func main() {
	// fmt.Println("Combyte CLI")
	myfigure := figure.NewColorFigure("COMBYTE CLI", "block", "BLUE", true)
	myfigure.Print()	
	fmt.Println(colortext.GreenString("By Edwino Maminirina"))

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
		fmt.Println(`
Command List
	--compress --filename="file.txt"
	--c --filename="file.txt"

	--extract --filename="file.txt.combyte"
	--e --filename="file.txt.combyte"

	--compress or --c : Compression a file, followed by the option --filename
	--extract or --e : Extract a file (.combyte file), followed by the option --filename
		`)
	}
}