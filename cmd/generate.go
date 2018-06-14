package main

import (
	"bufio"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strconv"
	"strings"
	"text/template"
)

const (
	managerNameFmt = "%sManager"
)

var typeTemplates *template.Template

type domain struct {
	managerName string
	fields      []field
	index       map[string]*index
	name        string
	tableName   string
	db          string
}

type field struct {
	name string
	typ  string
}

type index struct {
	name   string
	Fields []field
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
		"domainType":  m.domain.name,
		"dbName":      m.domain.db,
		"tableName":   m.domain.tableName,
	})

	writer.Flush()

	if err != nil {
		return err
	}
	return nil
}

func main() {
	filePath := flag.String("path", "domain/sample.go", "domain file path")
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
			pkg:     f.Name.Name,
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

	d := domain{index: make(map[string]*index)}
	d.name = object.Name
	d.managerName = fmt.Sprintf(managerNameFmt, object.Name)
	for _, v := range stt.Fields.List {
		f := field{name: v.Names[0].Name, typ: v.Type.(*ast.Ident).Name}
		d.fields = append(d.fields, f)

		tagMap := parseTag(v.Tag)
		if val, exist := tagMap["index"]; exist {
			for _, v := range val {
				split := strings.Split(strings.Replace(v, "\"", "", -1), ",")
				if len(split) >= 2 {
					fieldIndex, err := strconv.Atoi(split[1])
					if err != nil {
						continue
					}
					if pre, exist := d.index[split[0]]; exist {
						if len(pre.Fields) <= fieldIndex {
							fields := make([]field, fieldIndex+1)
							copy(fields, pre.Fields)
							fields[fieldIndex] = f
							pre.Fields = fields
						} else {
							pre.Fields[fieldIndex] = f
						}
					} else {
						fields := make([]field, fieldIndex+1)
						fields[fieldIndex] = f
						d.index[split[0]] = &index{name: split[0], Fields: fields}
					}
				}
			}
		} else if val, exist := tagMap["id"]; exist {
			replace := strings.Replace(val[0], "\"", "", -1)
			d.index[replace] = &index{name: replace, Fields: []field{f}}
		} else if val, exist := tagMap["db"]; exist {
			replace := strings.Replace(val[0], "\"", "", -1)
			split := strings.Split(replace, ",")
			if len(split) >= 2 {
				d.db = split[0]
				d.tableName = split[1]
			}
		}
	}
	return &d
}

// parseTag `index:"idx_a_b,1"` -> map["index"] = []string{"idx_a_b,1"}
func parseTag(lit *ast.BasicLit) map[string][]string {
	if lit == nil {
		return nil
	}

	val := strings.Replace(lit.Value, "`", "", -1)
	res := make(map[string][]string)
	splits := strings.Fields(val)
	for _, val := range splits {
		split := strings.Split(val, ":")
		if len(split) == 2 {
			if tmp, ok := res[split[0]]; ok {
				tmp = append(tmp, split[1])
			} else {
				res[split[0]] = []string{split[1]}
			}
		}
	}

	return res
}

func getFuncName(field []field) string {
	str := "SearchBy"
	for _, v := range field {
		str = str + v.name
	}

	return str
}

func getParameter(field []field) string {
	str := ""
	size := len(field)
	for k, v := range field {
		str += strings.ToLower(v.name) + " " + v.typ
		if k < size-1 {
			str += ","
		}
	}

	return str
}

func getKey(fields []field) string {
	str := ""
	for _, v := range fields {
		str += fmt.Sprintf("\t kc.Keys = append(kc.Keys, comparator.NewComparator(\"%s\", %s)) \n\r", v.typ,
			strings.ToLower(v.name))
	}

	return str
}

func init() {
	m := map[string]interface{}{
		"getFuncName":  getFuncName,
		"getParameter": getParameter,
		"getKey":       getKey,
	}

	tpl, err := template.New("managerTmp").Funcs(m).Parse(`
package {{ .packageName }}

import (
	"github.com/Kathent/mem-from-db/manager/comparator"
	"github.com/Kathent/mem-from-db/db/mysql"
	"github.com/Kathent/mem-from-db/manager"
)

type {{ .managerName }} struct {
	manager *manager.Manager
}

func New{{ .managerName }}(m *mysql.DbImpl) *{{ .managerName }} {
	mm := manager.NewManager(manager.TableConfig{
		DbName: "{{ .dbName }}",
		Name: "{{ .tableName }}",
		InitArr: make([]{{ .domainType }}, 0),
	}, m)

	return &{{ .managerName }}{
		manager: mm,
	}
}

{{ range $k, $v := .index}}
func (m *{{ $.managerName}}) {{ getFuncName $v.Fields}} ({{ getParameter $v.Fields}}) ({{ $.domainType }}, bool) {
	kc := &comparator.KeyValueComparator{
		Keys: make([]comparator.Comparator, 0),
	}

{{ getKey $v.Fields }}
	sample, ok  := m.manager.IR["{{ $k }}"].Tree.Search(kc).({{ $.domainType }})
	return sample, ok
}
{{ end }}
`)

	if err != nil {
		panic(err)
	}

	typeTemplates = tpl
}
