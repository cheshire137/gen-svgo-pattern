# gen-svgo-pattern

Generate Go code using the [svgo library](https://github.com/ajstarks/svgo) from a given
SVG image file. The generated code will include an
[SVG pattern](https://developer.mozilla.org/en-US/docs/Web/SVG/Tutorial/Patterns)
to reproduce the given image.

## How to run

```bash
make
bin/gen-svgo-pattern -w [WIDTH] -h [HEIGHT] -in [PATH_TO_SVG.svg] -name [NAME_OF_GO_TYPE] -out [PATH_TO_GO_FILE.go] -tab [INDENTATION_CHARACTERS]
```

Example:

```
make && bin/gen-svgo-pattern
  -h int
    	Height of pattern in pixels (default 200)
  -in string
    	Path to an SVG image file, e.g., ~/Pictures/my-pic.svg
  -name string
    	Name of Go type (default "MyPattern")
  -out string
    	Path where the Go code should be written, e.g., ~/my-go-project/pkg/patterns/my-pattern.go
  -pkg string
    	Name of Go package for new type (default "patterns")
  -tab string
    	Indentation to use in generated Go (default "\t")
  -w int
    	Width of pattern in pixels (default 200)
```

## Thanks

- [rustyoz/svg](https://github.com/rustyoz/svg)
