package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"golang.org/x/tools/go/ast/astutil"
)

var rootCmd = &cobra.Command{
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: support directory
		f, err := os.Open("testdata/main.go")
		if err != nil {
			log.Fatal(err)
		}

		if err := findStructs(bufio.NewReader(f)); err != nil {
			log.Fatal(err)
		}
	},
}

var (
	file string
	dir  string
)

func init() {
	rootCmd.PersistentFlags().StringVar(&file, "file", "file", "source file specified a sinle file")
	rootCmd.PersistentFlags().StringVar(&dir, "dir", "directory", "retrieve source files in directory")
}

func main() {
	rootCmd.Execute()
}

func findStructs(file io.Reader) error {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	fset := token.NewFileSet()
	af, err := parser.ParseFile(fset, "", data, parser.ParseComments)
	if err != nil {
		return err
	}

	fn := func(c *astutil.Cursor) bool {
		n := c.Node()
		if g, ok := n.(*ast.GenDecl); ok {
			if g.Tok == token.TYPE {
				// TODO: require check all specs
				if ts, ok := g.Specs[0].(*ast.TypeSpec); ok {
					if st, ok := ts.Type.(*ast.StructType); ok {
						// TODO: require collect all structs
						printMeta(ts.Name.Name, st)
						printMetaJson(ts.Name.Name, st)
						return false
					}
				}
			}
		}
		return true
	}
	astutil.Apply(af, fn, nil)

	return nil
}

// TODO: prototype
func printMeta(name string, node *ast.StructType) {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("## %s\n", name))
	buf.WriteString("| Field | Type | Description |\n")
	buf.WriteString("| ----- | ---- | ----------- |\n")
	for _, f := range node.Fields.List {
		v := f.Tag.Value[strings.Index(f.Tag.Value, "json:\"")+6:]
		v = v[:strings.Index(v, "\"")]
		buf.WriteString(fmt.Sprintf("| %s | %s | %s |\n", v, f.Type, strings.TrimSpace(f.Comment.Text())))
	}
	fmt.Println(buf.String())
}

// TODO: prototype
func printMetaJson(name string, node *ast.StructType) {
	pr := printer.Config{Mode: printer.TabIndent}
	pr.Fprint(os.Stdout, token.NewFileSet(), node)
}
