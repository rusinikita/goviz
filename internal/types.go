package internal

import (
	"go/types"
	"strings"
)

type Interface struct {
	name          string
	_type         *types.Named
	i             *types.Interface
	implementedBy []*Struct
}

func (i *Interface) NumMethods() int {
	return i.i.NumMethods()
}

func (i *Interface) Methods() (mm []string) {
	for mi := 0; mi < i.i.NumMethods(); mi++ {
		mm = append(mm, i.i.Method(mi).Name())
	}

	return mm
}

type Struct struct {
	name       string
	_type      *types.Named
	s          *types.Struct
	implements []*Interface
	includes   []string
	component  bool
}

func NewStruct(s *types.Struct, n *types.Named, pref string) *Struct {
	var includes []string
	var componentDeps int

	for i := 0; i < s.NumFields(); i++ {
		t := s.Field(i).Type()

		path := t.String()
		if !strings.HasPrefix(path, pref) {
			continue
		}

		switch t.Underlying().(type) {
		case *types.Interface, *types.Struct:
			componentDeps += 1
		}

		includes = append(includes, strings.TrimPrefix(path, pref))
	}

	return &Struct{
		name:      strings.TrimPrefix(n.String(), pref),
		_type:     n,
		s:         s,
		includes:  includes,
		component: true, //componentDeps > 0 && n.NumMethods() > 0,
	}
}

func (s *Struct) checkImplementation(i *Interface) {
	if !types.AssignableTo(s._type, i._type) &&
		!types.AssignableTo(types.NewPointer(s._type), i._type) {
		return
	}

	s.implements = append(s.implements, i)
	i.implementedBy = append(i.implementedBy, s)
}

func (s *Struct) NumMethods() int {
	return s._type.NumMethods()
}

func (s *Struct) Methods() (mm []string) {
	for mi := 0; mi < s._type.NumMethods(); mi++ {
		m := s._type.Method(mi)
		if !m.Exported() {
			continue
		}

		mm = append(mm, m.Name())
	}

	return mm
}

type Constructor struct {
	name string
	f    *types.Func
}
