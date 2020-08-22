package generator

import (
	"fmt"
	"os"

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
	g.writeFileHeader(outFile)
	g.writeConstructor(outFile)
	g.writeFillFunction(outFile)
	g.writeDefinePatternFunction(svgFile, outFile)
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
	outFile.WriteString(fmt.Sprintf("%sID string\n", g.tab))
	outFile.WriteString("}\n\n")
}

func (g *Generator) writeConstructor(outFile *os.File) {
	outFile.WriteString(fmt.Sprintf("func New%s() *%s {\n", g.typeName, g.typeName))
	outFile.WriteString(fmt.Sprintf("%sreturn &%s{\n", g.tab, g.typeName))
	outFile.WriteString(fmt.Sprintf("%s%sID: \"%s\",\n", g.tab, g.tab, g.typeName))
	outFile.WriteString(fmt.Sprintf("%s}\n", g.tab))
	outFile.WriteString("}\n\n")
}

func (g *Generator) writeFillFunction(outFile *os.File) {
	outFile.WriteString(fmt.Sprintf("func (p *%s) Fill() string {\n", g.typeName))
	fillStr := "fill:url(#%s)"
	outFile.WriteString(fmt.Sprintf("%sreturn fmt.Sprintf(\"%s\", p.ID)\n", g.tab, fillStr))
	outFile.WriteString("}\n\n")
}

func (g *Generator) writeDefinePatternFunction(svgFile *svg.Svg, outFile *os.File) {
	outFile.WriteString(fmt.Sprintf("func (p *%s) DefinePattern(canvas *svg.SVG) {\n", g.typeName))
	outFile.WriteString(fmt.Sprintf("%spw := %d\n", g.tab, g.width))
	outFile.WriteString(fmt.Sprintf("%sph := %d\n", g.tab, g.height))
	outFile.WriteString(fmt.Sprintf("%scanvas.Def()\n", g.tab))
	outFile.WriteString(fmt.Sprintf("%scanvas.Pattern(p.ID, 0, 0, pw, ph, \"user\")\n\n", g.tab))
	g.writeSvgElements(svgFile, outFile)
	outFile.WriteString(fmt.Sprintf("%scanvas.PatternEnd()\n", g.tab))
	outFile.WriteString(fmt.Sprintf("%scanvas.DefEnd()\n", g.tab))
	outFile.WriteString("}\n")
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
