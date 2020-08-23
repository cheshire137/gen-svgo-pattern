package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cheshire137/go-brocade/pkg/generator"
	"github.com/iancoleman/strcase"
	"github.com/rustyoz/svg"
)

func main() {
	var inPath string
	flag.StringVar(&inPath, "in", "",
		"Path to an SVG image file, e.g., ~/Pictures/my-pic.svg")

	var outDir string
	flag.StringVar(&outDir, "out", "",
		"Directory where the Go code should be written, e.g., ~/my-go-project/pkg/patterns/")

	var packageName string
	flag.StringVar(&packageName, "pkg", "patterns",
		"Name of Go package for new type")

	var typeName string
	flag.StringVar(&typeName, "name", "MyPattern",
		"Name of Go type")

	var width int
	flag.IntVar(&width, "w", 0,
		"Width of pattern in pixels; leave 0 to use width from input SVG")

	var height int
	flag.IntVar(&height, "h", 0,
		"Height of pattern in pixels; leave 0 to use width from input SVG")

	var tab string
	flag.StringVar(&tab, "tab", "\t",
		"Indentation to use in generated Go")

	flag.Parse()
	if len(inPath) < 1 || len(outDir) < 1 {
		fmt.Printf("Usage: %s [options]\n\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
		return
	}

	outDirInfo, err := os.Stat(outDir)
	if err != nil {
		fmt.Println("Could not check out directory: " + err.Error())
		os.Exit(1)
		return
	}

	if !outDirInfo.Mode().IsDir() {
		fmt.Printf("Error: %s is not a directory\n", outDir)
		os.Exit(1)
		return
	}

	fmt.Println("Reading: ", inPath)
	buf, err := ioutil.ReadFile(inPath)
	if err != nil {
		fmt.Println("Could not read file: " + err.Error())
		os.Exit(1)
		return
	}

	svgStr := string(buf)
	var scalefloat float64
	scalefloat = 1.0
	svgFile, err := svg.ParseSvg(svgStr, inPath, scalefloat)
	if err != nil {
		fmt.Println("Could not parse SVG: " + err.Error())
		os.Exit(1)
		return
	}

	outFilename := fmt.Sprintf("%s.go", strcase.ToSnake(typeName))
	outPath := filepath.Join(outDir, outFilename)
	outFile, err := os.Create(outPath)
	if err != nil {
		fmt.Println("Could not create Go file: " + err.Error())
		os.Exit(1)
		return
	}

	fmt.Printf("Generating Go type %s...\n", typeName)

	generator := generator.NewGenerator(packageName, tab, typeName, width, height)
	generator.WriteSvgCode(svgFile, outFile)

	fmt.Printf("Wrote %s\n", outPath)
}
