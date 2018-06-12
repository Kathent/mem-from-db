package main

import (
	"bufio"
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"os"
)

var typeTemplates, _ = template.New("managerTmp").Parse(`
package {% packageName %}

import "mem-from-db/manager"

type {% managerName%} struct {
	manager *manager.Manager
}
`)

var funcTemplates, _ = template.New("funcTmp").Parse(`
func (m *{% managerName%}) {% functionName%} ({% indexElem %}) {% domainType %} {
return m.ir[{ % indexName% }].Tree.Search()
}
`)

var compareTmp, _ = template.New("compares").Parse(`
	
`)

type domain struct {
	managerName string
	fields      []field
	index       map[string]index
}

type field struct {
	name string
}

type index struct {
	name   string
	fields []field
}

type managerGenerator struct {
	dstPath string
	pkg     string
	domain  *domain
}

func (m *managerGenerator) generate() error {
	var file *os.File
	var err error
	if _, err := os.Stat(m.dstPath); os.IsNotExist(err) {
		file, err = os.Create(m.dstPath)
	} else {
		file, err = os.OpenFile(m.dstPath, os.O_APPEND, 0644)
	}

	if err != nil {
		return err
	}
	writer := bufio.NewWriter(file)

	typeTemplates.Execute(writer, nil)
	return nil
}

func main() {
	filePath := flag.String("path", "cmd/sample.go", "domain file path")
	flag.Parse()

	fSet := token.NewFileSet()
	f, err := parser.ParseFile(fSet, *filePath, nil, parser.AllErrors)
	if err != nil {
		panic(err)
	}

	for _, val := range f.Scope.Objects {
		createManager(val)
	}
}

func createManager(object *ast.Object) {
	if object.Kind != ast.Typ {
		return
	}

	//name := object.Name
	spec, ok := object.Decl.(*ast.TypeSpec)
	if !ok {
		return
	}

	stt, ok := spec.Type.(*ast.StructType)
	if !ok {
		return
	}

	for _, v := range stt.Fields.List {
	}
}
