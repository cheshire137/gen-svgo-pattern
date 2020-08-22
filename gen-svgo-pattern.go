package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/rustyoz/svg"
)

func main() {
	var inPath string
	flag.StringVar(&inPath, "in", "",
		"Path to an SVG image file, e.g., ~/Pictures/my-pic.svg")

	var outPath string
	flag.StringVar(&outPath, "out", "",
		"Path where the Go code should be written, e.g., ~/my-go-project/pkg/patterns/my-pattern.go")

	var packageName string
	flag.StringVar(&packageName, "pkg", "patterns",
		"Name of Go package for new type")

	var typeName string
	flag.StringVar(&typeName, "name", "MyPattern",
		"Name of Go type")

	var width int
	flag.IntVar(&width, "w", 200,
		"Width of pattern in pixels")

	var height int
	flag.IntVar(&height, "h", 200,
		"Height of pattern in pixels")

	var tab string
	flag.StringVar(&tab, "tab", "\t",
		"Indentation to use in generated Go")

	flag.Parse()
	if len(inPath) < 1 || len(outPath) < 1 {
		flag.PrintDefaults()
		os.Exit(0)
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

	outFile, err := os.Create(outPath)
	if err != nil {
		fmt.Println("Could not create Go file: " + err.Error())
		os.Exit(1)
		return
	}

	fmt.Printf("Generating Go type %s...\n", typeName)
	outFile.WriteString(fmt.Sprintf("package %s\n\n", packageName))

	outFile.WriteString("import (\n")
	outFile.WriteString(fmt.Sprintf("%s\"fmt\"\n\n", tab))
	outFile.WriteString(fmt.Sprintf("%ssvg \"github.com/ajstarks/svgo\"\n", tab))
	outFile.WriteString(")\n\n")

	outFile.WriteString(fmt.Sprintf("type %s struct {\n", typeName))
	outFile.WriteString(fmt.Sprintf("%sID string\n", tab))
	outFile.WriteString("}\n\n")

	outFile.WriteString(fmt.Sprintf("func New%s() *%s {\n", typeName, typeName))
	outFile.WriteString(fmt.Sprintf("%sreturn &%s{\n", tab, typeName))
	outFile.WriteString(fmt.Sprintf("%s%sID: \"%s\",\n", tab, tab, typeName))
	outFile.WriteString(fmt.Sprintf("%s}\n", tab))
	outFile.WriteString("}\n\n")

	outFile.WriteString(fmt.Sprintf("func (p *%s) Fill() string {\n", typeName))
	fillStr := "fill:url(#%s)"
	outFile.WriteString(fmt.Sprintf("%sreturn fmt.Sprintf(\"%s\", p.ID)\n", tab, fillStr))
	outFile.WriteString("}\n\n")

	outFile.WriteString(fmt.Sprintf("func (p *%s) DefinePattern(canvas *svg.SVG) {\n", typeName))
	outFile.WriteString(fmt.Sprintf("%spw := %d\n", tab, width))
	outFile.WriteString(fmt.Sprintf("%sph := %d\n", tab, height))
	outFile.WriteString(fmt.Sprintf("%scanvas.Def()\n", tab))
	outFile.WriteString(fmt.Sprintf("%scanvas.Pattern(p.ID, 0, 0, pw, ph, \"user\")\n\n", tab))

	for _, group := range svgFile.Groups {
		if len(group.Fill) > 0 || len(group.Stroke) > 0 {
			style := fmt.Sprintf("fill:%s;stroke:%s", group.Fill, group.Stroke)
			outFile.WriteString(fmt.Sprintf("%scanvas.Gstyle(\"%s\")\n", tab, style))
		} else {
			outFile.WriteString(fmt.Sprintf("%scanvas.Gid(\"%s\")\n", tab, group.ID))
		}

		for _, groupEl := range group.Elements {
			path, ok := groupEl.(*svg.Path)
			if ok {
				outFile.WriteString(fmt.Sprintf("%scanvas.Path(\"%s\")\n", tab, path.D))
			}
		}

		outFile.WriteString(fmt.Sprintf("%scanvas.Gend()\n\n", tab))
	}

	outFile.WriteString(fmt.Sprintf("%scanvas.PatternEnd()\n", tab))
	outFile.WriteString(fmt.Sprintf("%scanvas.DefEnd()\n", tab))
	outFile.WriteString("}\n")

	fmt.Printf("Wrote %s\n", outPath)
}
