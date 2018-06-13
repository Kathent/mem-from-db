package main

import (
	"bufio"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"os"
	"strconv"
	"strings"
)

const (
	managerNameFmt = "%sManager"
)

var typeTemplates, _ = template.New("managerTmp").Parse(`
package {{ .packageName }}

import "mem-from-db/manager"

type {{ .managerName }} struct {
	manager *manager.Manager
}

{{ range $k, $v := range .index}}
func (m *{{ $.managerName}}) SearchBy{{ range $v.fields }} {{ .name }} {{ end }} ({{ range $v.fields }} {{ .name }} {{ .typ }} , {{end}}) {{ .domainType }} {
	kc := comparator.KeyValueComparator{
		Keys: make([]comparator.Comparator, 0)
	}
	{{ range .fields}}
		kc.Keys = append(kc.Keys, comparator.NewComparator({{ .typ }}, {{ .name }}))
	{{ end }}
	return m.manager.ir[{{ $k }}].Tree.Search(kc)
}
{{ end }}
`)

type domain struct {
	managerName string
	fields      []field
	index       map[string]index
}

type field struct {
	name string
	typ  string
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

	err = typeTemplates.Execute(writer, map[string]interface{}{
		"packageName": m.pkg,
		"managerName": m.domain.managerName,
		"index":       m.domain.index,
	})

	if err != nil {
		return err
	}
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
		domain := createDomain(val)
		if domain == nil {
			continue
		}

		generateErr := (&managerGenerator{
			dstPath: strings.Replace(*filePath, ".", "_manager.", 1),
			pkg:     "",
			domain:  domain,
		}).generate()

		if generateErr != nil {
			panic(generateErr)
		}
	}
}

func createDomain(object *ast.Object) *domain {
	if object.Kind != ast.Typ {
		return nil
	}

	//name := object.Name
	spec, ok := object.Decl.(*ast.TypeSpec)
	if !ok {
		return nil
	}

	stt, ok := spec.Type.(*ast.StructType)
	if !ok {
		return nil
	}

	d := domain{}
	d.managerName = fmt.Sprintf(managerNameFmt, object.Name)
	for _, v := range stt.Fields.List {
		f := field{name: v.Names[0].Name, typ: v.Type.(*ast.Ident).Name}
		d.fields = append(d.fields, f)

		tagMap := parseTag(v.Tag)
		if val, exist := tagMap["index"]; exist {
			for _, v := range val {
				split := strings.Split(v, ",")
				if len(split) >= 2 {
					fieldIndex, err := strconv.Atoi(split[1])
					if pre, exist := d.index[split[0]]; exist {
						if len(pre.fields) <= fieldIndex {
							fields := make([]field, fieldIndex+1)
							copy(fields, pre.fields)
							fields[fieldIndex] = f
							pre.fields = fields
						} else {
							pre.fields[fieldIndex] = f
						}
					} else {
						if err != nil {
							continue
						}

						fields := make([]field, fieldIndex+1)
						fields[fieldIndex] = f
						d.index[split[0]] = index{name: split[0], fields: fields}
					}
				}
			}
		} else if val, exist := tagMap["id"]; exist {
			d.index[val[0]] = index{name: val[0], fields: []field{f}}
		}
	}
	return &d
}

// parseTag `index:"idx_a_b,1"` -> map["index"] = []string{"idx_a_b,1"}
func parseTag(lit *ast.BasicLit) map[string][]string {
	if lit == nil {
		return nil
	}

	res := make(map[string][]string)
	splits := strings.Fields(lit.Value)
	for _, val := range splits {
		split := strings.Split(val, ":")
		if len(split) == 2 {
			if tmp, ok := res[split[0]]; ok {
				tmp = append(tmp, strings.Split(split[1], ",")...)
			} else {
				res[split[0]] = strings.Split(split[1], ",")
			}
		}
	}

	return res
}
