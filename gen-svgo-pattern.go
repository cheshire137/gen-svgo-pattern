package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/rustyoz/svg"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s file.svg\n", os.Args[0])
		return
	}

	filename := os.Args[1]
	fmt.Println("Reading:", filename)
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Could not read file: " + err.Error())
		return
	}

	svgStr := string(buf)
	var scalefloat float64
	scalefloat = 1.0
	file, err := svg.ParseSvg(svgStr, filename, scalefloat)
	if err != nil {
		fmt.Println("Could not parse SVG: " + err.Error())
		return
	}

	fmt.Println("canvas.Def()")
	fmt.Print("canvas.Pattern(p.ID, 0, 0, pw, ph, \"user\")\n\n")

	for _, group := range file.Groups {
		if len(group.Fill) > 0 || len(group.Stroke) > 0 {
			style := fmt.Sprintf("fill:%s;stroke:%s", group.Fill, group.Stroke)
			fmt.Printf("canvas.Gstyle(\"%s\")\n", style)
		} else {
			fmt.Printf("canvas.Gid(\"%s\")\n", group.ID)
		}

		for _, groupEl := range group.Elements {
			path, ok := groupEl.(*svg.Path)
			if ok {
				fmt.Printf("canvas.Path(\"%s\")\n", path.D)
			}
		}

		fmt.Print("canvas.Gend()\n\n")
	}

	fmt.Println("canvas.PatternEnd()")
	fmt.Println("canvas.DefEnd()")
}
