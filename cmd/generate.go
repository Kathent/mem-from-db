package main

import (
	"flag"
	"go/parser"
	"go/token"
	"go/ast"
)

var typeTemplates = `
	package {% packageName %}

	type {% managerName%} struct {
		manager *Manager
	}
`

var funcTemplates = `
	func (m *{% managerName%}) {% functionName%} ({% indexElem %}) {% domainType %} {
		return m.ir[{ % indexIndex% }].Tree.Search()
	}
`

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
	if object.Kind != ast.Typ{
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
		v.End()
	}
}