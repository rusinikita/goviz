package internal

import (
	"fmt"
	"go/types"
	"log"
	"strings"

	"github.com/samber/lo"
	"golang.org/x/tools/go/packages"
)

type CompileResult struct {
	interfaces   []*Interface
	structs      []*Struct
	constructors []*Constructor
}

func Compile(projectDir string) *CompileResult {
	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedTypes |
			packages.NeedModule,
		Dir: projectDir,
	}
	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		log.Fatal(err)
	}

	var interfaces []*Interface
	var structs []*Struct
	var constructors []*Constructor

	if len(pkgs) == 0 {
		return nil
	}

	pathPrefixToRemove := pkgs[0].Module.Path + "/"

	for _, pkg := range pkgs {
		for _, name := range pkg.Types.Scope().Names() {
			object := pkg.Types.Scope().Lookup(name)

			if obj, ok := object.(*types.TypeName); ok {
				objType := obj.Type()
				named, ok := objType.(*types.Named)
				if !ok {
					continue
				}

				if i, ok := named.Underlying().(*types.Interface); ok {
					interfaces = append(interfaces, &Interface{
						name:  strings.TrimPrefix(named.String(), pathPrefixToRemove),
						_type: named,
						i:     i,
					})
				}

				if s, ok := named.Underlying().(*types.Struct); ok {
					structs = append(structs, NewStruct(s, named, pathPrefixToRemove))
				}
			}

			if f, ok := object.(*types.Func); ok {
				if !strings.HasPrefix(f.Name(), "New") {
					continue
				}

				constructors = append(constructors, &Constructor{
					name: f.String(),
					f:    f,
				})
			}
		}
	}

	log.Println("collected")

	interfaces = lo.Filter(interfaces, func(item *Interface, _ int) bool {
		return item.NumMethods() > 1
	})

	structs = lo.Filter(structs, func(item *Struct, _ int) bool {
		return item.component
	})

	for _, s := range structs {
		for _, i := range interfaces {
			s.checkImplementation(i)
		}
	}

	fmt.Println("checked")

	return &CompileResult{
		interfaces:   interfaces,
		structs:      structs,
		constructors: constructors,
	}
}
