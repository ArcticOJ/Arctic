package main

import (
	"errors"
	"fmt"
	"github.com/ArcticOJ/blizzard/v0/logger"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path"
	"slices"
	"strings"
	"text/template"
)

var (
	ErrMalformedManifest = errors.New("malformed manifest")
)

var Template = template.Must(template.New("generated_map.tmpl").ParseFiles("cmd/gen_routes/generated_map.tmpl"))

type (
	TemplateData struct {
		Args      string
		Imports   []string
		Manifests map[string][]RouteManifest
	}
	RouteManifest struct {
		Method  string
		Path    string
		Handler string
		Flags   string
	}
)

func parseManifest(pkg, raw string) (RouteManifest, error) {
	// trim comments' prefix before processing
	fields := strings.Fields(strings.TrimSpace(strings.TrimLeft(raw, "/")))
	m := RouteManifest{}
	if len(fields) < 3 {
		return m, ErrMalformedManifest
	}
	m.Handler, m.Method, m.Path = fields[0], fields[1], fields[2]
	if pkg != "" {
		m.Handler = pkg + "." + m.Handler
	}
	if len(m.Path) > 1 {
		m.Path = strings.TrimSuffix(m.Path, "/")
	}
	var flags []string
	for i := 3; i < len(fields); i++ {
		// no need to use strings.HasPrefix lol
		flag := strings.TrimPrefix(fields[i], "@")
		// not a valid flag
		if len(flag) == 0 || flag == fields[i] {
			continue
		}
		flags = append(flags, "http.Route"+strings.Title(flag))
	}
	m.Flags = strings.Join(flags, "|")
	return m, nil
}

func main() {
	basePath, pkgPath, outFile := os.Args[1], os.Args[2], os.Args[3]
	dat := TemplateData{
		Args: strings.Join(os.Args, " "),
		Imports: []string{
			pkgPath + "/server/http",
		},
		Manifests: make(map[string][]RouteManifest),
	}
	build := func(dir string) {
		var manifests []RouteManifest
		fset := token.NewFileSet()
		_path := path.Join(basePath, dir)
		pkgs, e := parser.ParseDir(fset, _path, func(info fs.FileInfo) bool {
			return !(strings.HasPrefix(info.Name(), "._") || strings.HasSuffix(info.Name(), "generated.go"))
		}, parser.ParseComments)
		logger.Panic(e, "error reading '%s'", _path)
		if len(pkgs) > 0 && dir != "" {
			dat.Imports = append(dat.Imports, fmt.Sprintf("%s/routes/%s", pkgPath, dir))
		}
		for pkg, content := range pkgs {
			for _p, x := range content.Files {
				for _, decl := range x.Decls {
					if fn, ok := decl.(*ast.FuncDecl); ok {
						if fn.Doc == nil || len(fn.Doc.List) == 0 {
							continue
						}
						// nullify pkg for apex routes
						if dir == "" {
							pkg = ""
						}
						manifest, e := parseManifest(pkg, fn.Doc.List[0].Text)
						logger.Panic(e, "error parsing '%s'", _p)
						manifests = append(manifests, manifest)
					}
				}
			}
		}
		slices.SortStableFunc(manifests, func(a, b RouteManifest) int {
			if a.Path > b.Path {
				return 1
			} else if a.Path < b.Path {
				return -1
			}
			return 0
		})
		dat.Manifests["/"+dir] = manifests
	}
	build("")
	dirs, e := os.ReadDir(basePath)
	logger.Panic(e, "failed to read routes dir")
	for _, dir := range dirs {
		if dir.IsDir() {
			build(dir.Name())
		}
	}
	slices.Sort(dat.Imports)
	f, e := os.OpenFile(outFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
	logger.Panic(e, "error opening output file '%s'", outFile)
	logger.Panic(Template.Execute(f, dat), "error generating output file")
	logger.Global.Info().Msgf("generated '%s'", outFile)
}
