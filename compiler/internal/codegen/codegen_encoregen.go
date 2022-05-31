package codegen

import (
	"strings"

	. "github.com/dave/jennifer/jen"

	"encr.dev/parser/est"
	"encr.dev/parser/paths"
)

// EncoreGen generates the encore.gen.go file containing user-facing
// generated code. If nothing needs to be generated it returns nil.
func (b *Builder) EncoreGen(svc *est.Service) *File {
	if len(svc.RPCs) == 0 {
		return nil
	}

	f := NewFilePathName(svc.Root.ImportPath, svc.Root.Name)
	f.ImportNames(importNames)

	f.HeaderComment("Code generated by encore. DO NOT EDIT.")

	f.Commentf("NewClient returns a Client suitable for calling the %s service.", svc.Root.Name)
	f.Func().Id("NewClient").Params().Id("Client").Block(
		Comment("The implementation is elided and generated at compile time by Encore."),
		Return(Nil()),
	)
	f.Line()

	f.Commentf("Client is a client for calling the %s service.", svc.Root.Name)
	f.Type().Id("Client").InterfaceFunc(func(g *Group) {
		for i, rpc := range svc.RPCs {
			b.encoreGenRPC(f, g, rpc, i == 0)
		}
	})

	return f
}

func (b *Builder) encoreGenRPC(f *File, g *Group, rpc *est.RPC, first bool) {
	if rpc.Doc != "" {
		// If it's not the first RPC, add a blank line between the previous RPC and this doc comment.
		if !first {
			g.Line()
		}
		for _, line := range strings.Split(strings.TrimSpace(rpc.Doc), "\n") {
			g.Comment(line)
		}
	}

	g.Id(rpc.Name).ParamsFunc(func(g *Group) {
		var names nameAllocator
		g.Id(names.Get("ctx")).Qual("context", "Context")
		for _, seg := range rpc.Path.Segments {
			if seg.Type != paths.Literal {
				g.Id(names.Get(seg.Value)).Add(b.builtinType(seg.ValueType))
			}
		}
		if rpc.Raw {
			g.Id(names.Get("req")).Op("*").Qual("net/http", "Request")
		} else if req := rpc.Request; req != nil {
			g.Id(names.Get("req")).Add(b.namedType(f, req))
		}
	}).Do(func(s *Statement) {
		if rpc.Raw {
			s.Params(Op("*").Qual("net/http", "Response"), Error())
		} else if resp := rpc.Response; resp != nil {
			s.Params(b.namedType(f, resp), Error())
		} else {
			s.Error()
		}
	})
}
