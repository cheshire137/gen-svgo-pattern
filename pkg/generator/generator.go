package generator

import (
	"fmt"
	"os"
	"strconv"

	"github.com/rustyoz/svg"
)

type Generator struct {
	packageName string
	tab         string
	typeName    string
	width       int
	height      int
}

func NewGenerator(packageName string, tab string, typeName string, width int, height int) *Generator {
	return &Generator{
		packageName: packageName,
		tab:         tab,
		typeName:    typeName,
		width:       width,
		height:      height,
	}
}

func (g *Generator) WriteSvgCode(svgFile *svg.Svg, outFile *os.File) {
	patternWidth := g.getWidth(svgFile)
	fmt.Printf("Using width %d\n", patternWidth)

	patternHeight := g.getHeight(svgFile)
	fmt.Printf("Using height %d\n", patternHeight)

	g.writeFileHeader(outFile)
	g.writeConstructor(patternWidth, patternHeight, outFile)
	g.writeDefinePatternFunction(svgFile, outFile)
	g.writeStyleFunction(outFile)
}

func (g *Generator) writeFileHeader(outFile *os.File) {
	outFile.WriteString(fmt.Sprintf("package %s\n\n", g.packageName))
	g.writeImports(outFile)
	g.writeTypeDefinition(outFile)
}

func (g *Generator) writeImports(outFile *os.File) {
	outFile.WriteString("import (\n")
	outFile.WriteString(fmt.Sprintf("%s\"fmt\"\n\n", g.tab))
	outFile.WriteString(fmt.Sprintf("%ssvg \"github.com/ajstarks/svgo\"\n", g.tab))
	outFile.WriteString(")\n\n")
}

func (g *Generator) writeTypeDefinition(outFile *os.File) {
	outFile.WriteString(fmt.Sprintf("type %s struct {\n", g.typeName))
	outFile.WriteString(fmt.Sprintf("%sID            string\n", g.tab))
	outFile.WriteString(fmt.Sprintf("%smaskID        string\n", g.tab))
	outFile.WriteString(fmt.Sprintf("%spatternWidth  int\n", g.tab))
	outFile.WriteString(fmt.Sprintf("%spatternHeight int\n", g.tab))
	outFile.WriteString("}\n\n")
}

func (g *Generator) writeConstructor(patternWidth int, patternHeight int, outFile *os.File) {
	outFile.WriteString(fmt.Sprintf("func New%s() *%s {\n", g.typeName, g.typeName))
	outFile.WriteString(fmt.Sprintf("%sreturn &%s{\n", g.tab, g.typeName))
	outFile.WriteString(fmt.Sprintf("%s%sID:            \"%s\",\n", g.tab, g.tab, g.typeName))
	outFile.WriteString(fmt.Sprintf("%s%smaskID:        \"%s-mask\",\n", g.tab, g.tab, g.typeName))
	outFile.WriteString(fmt.Sprintf("%s%spatternWidth:  %d,\n", g.tab, g.tab, patternWidth))
	outFile.WriteString(fmt.Sprintf("%s%spatternHeight: %d,\n", g.tab, g.tab, patternHeight))
	outFile.WriteString(fmt.Sprintf("%s}\n", g.tab))
	outFile.WriteString("}\n\n")
}

func (g *Generator) writeDefinePatternFunction(svgFile *svg.Svg, outFile *os.File) {
	outFile.WriteString(fmt.Sprintf("func (p *%s) DefinePattern(width int, height int, canvas *svg.SVG) {\n", g.typeName))
	outFile.WriteString(fmt.Sprintf("%scanvas.Def()\n", g.tab))
	outFile.WriteString(fmt.Sprintf("%scanvas.Pattern(p.ID, 0, 0, p.patternWidth, p.patternHeight, \"user\", \"stroke:white;stroke-linecap:square;stroke-width:1\")\n\n", g.tab))
	if len(svgFile.Elements) > 0 {
		g.writeSvgElements(svgFile, outFile)
	} else if len(svgFile.Groups) > 0 {
		g.writeSvgGroups(svgFile, outFile)
	}
	outFile.WriteString(fmt.Sprintf("%scanvas.PatternEnd()\n\n", g.tab))
	outFile.WriteString(fmt.Sprintf("%scanvas.Mask(p.maskID, 0, 0, width, height)\n", g.tab))
	fillStr := "fill:url(#%s)"
	outFile.WriteString(fmt.Sprintf("%scanvas.Rect(0, 0, width, height, fmt.Sprintf(\"%s\", p.ID))\n", g.tab, fillStr))
	outFile.WriteString(fmt.Sprintf("%scanvas.MaskEnd()\n\n", g.tab))
	outFile.WriteString(fmt.Sprintf("%scanvas.DefEnd()\n", g.tab))
	outFile.WriteString("}\n\n")
}

func (g *Generator) writeStyleFunction(outFile *os.File) {
	outFile.WriteString(fmt.Sprintf("func (p *%s) Style(color string, offsetX int, offsetY int) string {\n", g.typeName))
	styleStr := "mask:url(#%s);fill:%s;transform:translate(%dpx, %dpx)"
	outFile.WriteString(fmt.Sprintf("%sreturn fmt.Sprintf(\"%s\", p.maskID, color, offsetX, offsetY)\n", g.tab, styleStr))
	outFile.WriteString("}\n")
}

func (g *Generator) getWidth(svgFile *svg.Svg) int {
	if g.width > 0 {
		return g.width
	}

	svgWidth, err := strconv.Atoi(svgFile.Width)
	if err == nil {
		return svgWidth
	}

	return 200
}

func (g *Generator) getHeight(svgFile *svg.Svg) int {
	if g.height > 0 {
		return g.height
	}

	svgHeight, err := strconv.Atoi(svgFile.Height)
	if err == nil {
		return svgHeight
	}

	return 200
}

func (g *Generator) writeSvgGroups(svgFile *svg.Svg, outFile *os.File) {
	for _, group := range svgFile.Groups {
		g.writeSvgGroup(&group, outFile)
	}
}

func (g *Generator) writeSvgElements(svgFile *svg.Svg, outFile *os.File) {
	for _, el := range svgFile.Elements {
		group, ok := el.(*svg.Group)
		if ok {
			g.writeSvgGroup(group, outFile)
		} else {
			path, ok := el.(*svg.Path)
			if ok {
				g.writeSvgPath(path, outFile)
			}
		}
	}
}

func (g *Generator) writeSvgGroup(group *svg.Group, outFile *os.File) {
	if len(group.Fill) > 0 || len(group.Stroke) > 0 {
		g.writeSvgGroupStyle(group, outFile)
	} else if len(group.ID) > 0 {
		outFile.WriteString(fmt.Sprintf("%scanvas.Gid(\"%s\")\n", g.tab, group.ID))
	} else {
		outFile.WriteString(fmt.Sprintf("%scanvas.Group(\"\")\n", g.tab))
	}

	for _, groupEl := range group.Elements {
		path, ok := groupEl.(*svg.Path)
		if ok {
			g.writeSvgPath(path, outFile)
		}
	}

	outFile.WriteString(fmt.Sprintf("%scanvas.Gend()\n\n", g.tab))
}

func (g *Generator) writeSvgGroupStyle(group *svg.Group, outFile *os.File) {
	var style string
	if len(group.Fill) > 0 && len(group.Stroke) > 0 {
		style = fmt.Sprintf("fill:%s;stroke:%s", group.Fill, group.Stroke)
	} else if len(group.Fill) > 0 {
		style = fmt.Sprintf("fill:%s", group.Fill)
	} else {
		style = fmt.Sprintf("stroke:%s", group.Stroke)
	}
	outFile.WriteString(fmt.Sprintf("%scanvas.Gstyle(\"%s\")\n", g.tab, style))
}

func (g *Generator) writeSvgPath(path *svg.Path, outFile *os.File) {
	outFile.WriteString(fmt.Sprintf("%scanvas.Path(\"%s\")\n", g.tab, path.D))
}
