// Package dynamo implements the protoc-gen-dynamo plugin.
//
// This module adds DynamoDB struct tags to generated Go protobuf code by:
// 1. Extracting dynamo annotations from proto files
// 2. Reading the generated .pb.go files
// 3. Modifying the Go AST to inject struct tags
// 4. Writing the updated files back
package godynamo

import (
	"go/parser"
	"go/printer"
	"go/token"
	"path/filepath"
	"strings"

	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"
)

type mod struct {
	*pgs.ModuleBase
	pgsgo.Context
}

func New() pgs.Module {
	return &mod{ModuleBase: &pgs.ModuleBase{}}
}

func (m *mod) InitContext(c pgs.BuildContext) {
	m.ModuleBase.InitContext(c)
	m.Context = pgsgo.InitContext(c.Parameters())
}

func (mod) Name() string {
	return "dynamo"
}

// Execute processes proto files and adds DynamoDB tags to generated Go code.
// For each file with dynamo annotations:
// 1. Extract tag mappings (field name → tag string)
// 2. Parse the corresponding .pb.go file's AST
// 3. Inject tags into struct field definitions
// 4. Write the modified Go code back
func (m mod) Execute(targets map[string]pgs.File, packages map[string]pgs.Package) []pgs.Artifact {
	outdir := m.Parameters().Str("outdir")
	extractor := newTagExtractor(m, m.Context)

	for _, f := range targets {
		tags := extractor.Extract(f)
		if len(tags) == 0 {
			continue // No dynamo annotations in this file
		}

		// Get .pb.go filename and determine read location
		gfname := m.Context.OutputPath(f).SetExt(".go").String()
		filename := gfname
		if outdir != "" {
			filename = filepath.Join(outdir, gfname)
		}

		// Parse, modify, and write back
		fs := token.NewFileSet()
		fn, err := parser.ParseFile(fs, filename, nil, parser.ParseComments)
		m.CheckErr(err)
		m.CheckErr(Retag(fn, tags))

		var buf strings.Builder
		m.CheckErr(printer.Fprint(&buf, fs, fn))
		m.OverwriteGeneratorFile(gfname, buf.String())
	}

	return m.Artifacts()
}
