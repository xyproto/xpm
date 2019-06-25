package main

import (
	"flag"
	"fmt"
	"github.com/xyproto/xpm"
	"image/png"
	"os"
	"path/filepath"
)

func main() {

	var (
		outputFilename string
		version        bool
	)

	flag.StringVar(&outputFilename, "o", "-", "output XPM filename")
	flag.BoolVar(&version, "v", false, "version")

	flag.Parse()

	if version {
		fmt.Println("png2xpm 1.2.0")
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "An input PNG filename is required.\n")
		os.Exit(1)
	}

	inputFilename := args[0]

	// Choose a name to use for the struct, in the XPM data
	imageName := ""
	if outputFilename == "-" {
		// Pick out the part of the input filename that is not the path and not the extension
		inputFilenameBase := filepath.Base(inputFilename)
		imageName = inputFilenameBase[:len(inputFilenameBase)-len(filepath.Ext(inputFilenameBase))]
	} else {
		// Pick out the part of the output filename that is not the path and not the extension
		outputFilenameBase := filepath.Base(outputFilename)
		imageName = outputFilenameBase[:len(outputFilenameBase)-len(filepath.Ext(outputFilenameBase))]
	}

	// Create a new XPM encoder
	enc := xpm.NewEncoder(imageName)

	// Open the PNG file
	f, err := os.Open(inputFilename)
	m, err := png.Decode(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
	f.Close()

	// Prepare to output the XPM data to either stdout or to file
	if outputFilename == "-" {
		f = os.Stdout
	} else {
		f, err = os.Create(outputFilename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	// Generate and output the XPM data
	err = enc.Encode(f, m)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}
