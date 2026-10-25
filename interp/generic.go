package interp

import (
	"strings"
	"sync/atomic"
)

// GenericFunc contains the code of a generic function.
// This is used in the `yaegi extract` command to represent generic functions
// instead of the actual value of the function, since you cannot get the
// [reflect.Value] of a generic function. This is then used to interpret the
// function when it is imported in yaegi.
type GenericFunc string

// adot produces an AST dot(1) directed acyclic graph for the given node. For debugging only.
// func (n *node) adot() { n.astDot(dotWriter(n.interp.dotCmd), n.ident) }

// genAST returns a new AST where generic types are replaced by instantiated types.
func genAST(sc *scope, root *node, types []*itype) (*node, bool, error) {
	typeParam := map[string]*node{}
	pindex := 0
	tname := ""
	rtname := ""
	recvrPtr := false
	fixNodes := []*node{}
	var gtree func(n, anc *node, inBody bool) (*node, error)
	sname := root.child[0].ident + "["
	if root.kind == funcDecl {
		sname = root.child[1].ident + "["
	}

	// Input type parameters must be resolved prior AST generation, as compilation
	// of generated AST may occur in a different scope.
	for _, t := range types {
		sname += t.id() + ","
	}
	sname = strings.TrimSuffix(sname, ",") + "]"

	gtree = func(n, anc *node, inBody bool) (*node, error) {
		nod := copyNode(n, anc, false)
		switch n.kind {
		case funcDecl, funcType:
			nod.val = nod

		case identExpr:
			// Replace generic type by instantiated one.
			nt, ok := typeParam[n.ident]
			if !ok {
				// The clone may refer to names which are visible at the
				// definition site but not at the instantiation site (types
				// and packages of the generic package). Resolve them in the
				// definition scope, so nested instantiations complete instead
				// of degrading to incomplete types, except in function bodies
				// which compile in the definition scope already. Names visible
				// at the instantiation site keep their normal resolution, as
				// well as declaration names which are never resolved.
				if !inBody && !isGenericName(root, n) {
					if _, _, found := sc.lookup(n.ident); !found {
						if typ := lookupDefinitionType(root, n.ident); typ != nil {
							nod.typ = typ
						}
					}
				}
				break
			}
			nod = copyNode(nt, anc, true)
			nod.typ = nt.typ

		case indexExpr:
			// Catch a possible recursive generic type definition
			if root.kind != typeSpec {
				break
			}
			if root.child[0].ident != n.child[0].ident {
				break
			}
			nod := copyNode(n.child[0], anc, false)
			fixNodes = append(fixNodes, nod)
			return nod, nil

		case fieldList:
			//  Node is the type parameters list of a generic function.
			if root.kind == funcDecl && n.anc == root.child[2] && childPos(n) == 0 && len(types) > 0 {
				// Fill the types lookup table used for type substitution.
				for _, c := range n.child {
					l := len(c.child) - 1
					for _, cc := range c.child[:l] {
						if pindex >= len(types) {
							return nil, cc.cfgErrorf("undefined type for %s", cc.ident)
						}
						t, err := nodeType(c.interp, sc, c.child[l])
						if err != nil {
							return nil, err
						}
						if err := checkConstraint(types[pindex], t); err != nil {
							return nil, err
						}
						typeParam[cc.ident] = copyNode(cc, cc.anc, false)
						typeParam[cc.ident].ident = types[pindex].id()
						typeParam[cc.ident].typ = types[pindex]
						pindex++
					}
				}
				// Skip type parameters specification, so generated func doesn't look generic.
				return nod, nil
			}

			// Node is the receiver of a generic method.
			if root.kind == funcDecl && n.anc == root && childPos(n) == 0 && len(n.child) > 0 {
				rtn := n.child[0].child[1]
				// Method receiver is a generic type if it takes some type parameters.
				if rtn.kind == indexExpr || rtn.kind == indexListExpr || (rtn.kind == starExpr && (rtn.child[0].kind == indexExpr || rtn.child[0].kind == indexListExpr)) {
					if rtn.kind == starExpr {
						// Method receiver is a pointer on a generic type.
						rtn = rtn.child[0]
						recvrPtr = true
					}
					rtname = rtn.child[0].ident + "["
					for _, cc := range rtn.child[1:] {
						if pindex >= len(types) {
							return nil, cc.cfgErrorf("undefined type for %s", cc.ident)
						}
						it := types[pindex]
						typeParam[cc.ident] = copyNode(cc, cc.anc, false)
						typeParam[cc.ident].ident = it.id()
						typeParam[cc.ident].typ = it
						rtname += it.id() + ","
						pindex++
					}
					rtname = strings.TrimSuffix(rtname, ",") + "]"
				}
			}

			// Node is the type parameters list of a generic type.
			if root.kind == typeSpec && n.anc == root && childPos(n) == 1 {
				// Fill the types lookup table used for type substitution.
				tname = n.anc.child[0].ident + "["
				for _, c := range n.child {
					l := len(c.child) - 1
					for _, cc := range c.child[:l] {
						if pindex >= len(types) {
							return nil, cc.cfgErrorf("undefined type for %s", cc.ident)
						}
						it := types[pindex]
						t, err := nodeType(c.interp, sc, c.child[l])
						if err != nil {
							return nil, err
						}
						if err := checkConstraint(types[pindex], t); err != nil {
							return nil, err
						}
						typeParam[cc.ident] = copyNode(cc, cc.anc, false)
						typeParam[cc.ident].ident = it.id()
						typeParam[cc.ident].typ = it
						tname += it.id() + ","
						pindex++
					}
				}
				tname = strings.TrimSuffix(tname, ",") + "]"
				return nod, nil
			}
		}

		for _, c := range n.child {
			gn, err := gtree(c, nod, inBody || c.kind == blockStmt)
			if err != nil {
				return nil, err
			}
			nod.child = append(nod.child, gn)
		}
		return nod, nil
	}

	if nod, found := root.interp.generic[sname]; found {
		return nod, true, nil
	}

	r, err := gtree(root, root.anc, false)
	if err != nil {
		return nil, false, err
	}
	root.interp.generic[sname] = r
	r.param = append(r.param, types...)
	if tname != "" {
		for _, nod := range fixNodes {
			nod.ident = tname
		}
		r.child[0].ident = tname
	}
	if rtname != "" {
		// Replace method receiver type by synthetized ident.
		nod := r.child[0].child[0].child[1]
		if recvrPtr {
			nod = nod.child[0]
		}
		nod.kind = identExpr
		nod.ident = rtname
		nod.child = nil
	}
	// r.adot() // Used for debugging only.
	return r, false, nil
}

// isGenericName reports whether n is the declared name of a generic type or
// function being cloned, which is never resolved and must keep no stamped type.
func isGenericName(root, n *node) bool {
	if n.kind != identExpr {
		return false
	}
	switch root.kind {
	case typeSpec:
		return len(root.child) > 0 && n == root.child[0]
	case funcDecl:
		index := 0
		if len(root.child) > 0 && root.child[0].kind == fieldList {
			index = 1
		}
		return len(root.child) > index && n == root.child[index]
	}
	return false
}

// lookupDefinitionType resolves a bare identifier in the scope where a generic
// definition was declared, for nodes cloned by genAST. It returns the declared
// type for visible type and package symbols with complete types, or nil when the
// name is not visible there, so callers keep the normal resolution path.
func lookupDefinitionType(root *node, ident string) *itype {
	if root == nil || root.scope == nil {
		return nil
	}
	if sym, _, found := root.scope.lookup(ident); found && sym.typ != nil && !sym.typ.incomplete &&
		(sym.kind == typeSym || sym.kind == pkgSym) {
		return sym.typ
	}
	if sym, ok := lookupPackageByBase(root.scope, ident); ok && sym.typ != nil && !sym.typ.incomplete {
		return sym.typ
	}
	return nil
}

func copyNode(n, anc *node, recursive bool) *node {
	var i interface{}
	nindex := atomic.AddInt64(&n.interp.nindex, 1)
	nod := &node{
		debug:  n.debug,
		anc:    anc,
		interp: n.interp,
		index:  nindex,
		level:  n.level,
		nleft:  n.nleft,
		nright: n.nright,
		kind:   n.kind,
		pos:    n.pos,
		action: n.action,
		gen:    n.gen,
		val:    &i,
		rval:   n.rval,
		ident:  n.ident,
		meta:   n.meta,
	}
	nod.start = nod
	if recursive {
		for _, c := range n.child {
			nod.child = append(nod.child, copyNode(c, nod, true))
		}
	}
	return nod
}

func inferTypesFromCall(sc *scope, fun *node, args []*node) ([]*itype, error) {
	ftn := fun.typ.node
	// Fill the map of parameter types, indexed by type param ident.
	paramTypes := map[string]*itype{}
	for _, c := range ftn.child[0].child {
		typ, err := nodeType(fun.interp, sc, c.lastChild())
		if err != nil {
			return nil, err
		}
		for _, cc := range c.child[:len(c.child)-1] {
			paramTypes[cc.ident] = typ
		}
	}

	// inferred collects the concrete type for each type parameter name.
	// unify matches parameter shapes against argument types and records by
	// name, so multi-name fields (s1, s2 S) and constraint-only type params
	// (E in S ~[]E, K/V in M ~map[K]V) resolve correctly. It replaces the
	// old positional inferTypes, which dropped E/K/V entirely.
	inferred := map[string]*itype{}
	// record stores input under name and, when name looks like a suffixed
	// copy of a declared type param (e.g. V1/V2 for V in Equal[M1, M2]),
	// also under the base name so declaration-order emit finds it.
	record := func(name string, input *itype) {
		if name == "" || input == nil {
			return
		}
		if _, ok := inferred[name]; !ok {
			inferred[name] = input
		}
		// Strip trailing digits to find base (V1 -> V). Only apply when
		// the full name is not itself a declared param but base is.
		if paramTypes[name] == nil {
			base := strings.TrimRight(name, "0123456789")
			if base != name && paramTypes[base] != nil {
				if _, ok := inferred[base]; !ok {
					inferred[base] = input
				}
			}
		}
	}
	var unify func(param, input *itype)
	unify = func(param, input *itype) {
		if param == nil || input == nil {
			return
		}
		switch param.cat {
		case genericT:
			record(param.name, input)
			// Decompose the type parameter constraint as well, so
			// constraint-only type params get inferred from the same arg.
			if param.val != nil {
				unify(param.val, input)
			}
		case nilT:
			// Caller-scope view of a type parameter (the #1700 nil-type
			// case): record the argument when the name is a known type param.
			if paramTypes[param.name] != nil {
				record(param.name, input)
			}
		case chanT, ptrT, sliceT, arrayT:
			if param.val != nil && input.val != nil {
				unify(param.val, input.val)
			}
		case variadicT:
			// Variadic params (vals ...T) receive individual elements,
			// not a slice, when called as f(a, b). Infer T directly
			// from the element type; also handle f(slice...) where
			// input is a slice/array carrying the element in val.
			if param.val != nil {
				if input.val != nil && (input.cat == sliceT || input.cat == arrayT) {
					unify(param.val, input.val)
				} else {
					unify(param.val, input)
				}
			}
		case mapT:
			unify(param.key, input.key)
			unify(param.val, input.val)
		case structT:
			for i, f := range param.field {
				if i >= len(input.field) {
					break
				}
				unify(f.typ, input.field[i].typ)
			}
		case funcT:
			for i, t := range param.arg {
				if i >= len(input.arg) {
					break
				}
				unify(t, input.arg[i])
			}
			for i, t := range param.ret {
				if i >= len(input.ret) {
					break
				}
				unify(t, input.ret[i])
			}
		}
	}

	// resolveParam prefers the generic function's own scope, where type
	// parameters are visible as genericT with their constraints. The caller
	// scope sees them as nilT (the #1700 nil-type case), which cannot be
	// decomposed for constraint-only type params (e.g. E in Equal[S ~[]E]).
	resolveParam := func(n *node) *itype {
		if fun.scope != nil {
			if t, err := nodeType(fun.interp, fun.scope, n); err == nil && t != nil && t.cat == genericT {
				return t
			}
		}
		t, _ := nodeType(fun.interp, sc, n)
		return t
	}

	types := []*itype{}
	ai := 0 // index into args, advancing per expanded param name
	for _, c := range ftn.child[1].child {
		typ := resolveParam(c.lastChild())
		if typ == nil {
			continue
		}
		names := c.child[:len(c.child)-1]
		if len(names) == 0 {
			names = c.child
		}
		for range names {
			// A missing argument (fewer args than params) leaves the
			// type param uninferred; genAST then reports the first
			// undefined one, and check.arguments reports arity.
			if ai < len(args) {
				unify(typ, args[ai].typ)
				ai++
			}
		}
	}
	// Emit inferred types in type parameter declaration order, so genAST
	// maps each type param positionally. Missing entries stay missing and
	// genAST reports the first undefined one, as before.
	for _, c := range ftn.child[0].child {
		for _, cc := range c.child[:len(c.child)-1] {
			if t, ok := inferred[cc.ident]; ok {
				types = append(types, t)
			}
		}
	}

	return types, nil
}

func checkConstraint(it, ct *itype) error {
	if len(ct.constraint) == 0 && len(ct.ulconstraint) == 0 {
		return nil
	}
	for _, c := range ct.constraint {
		if it.equals(c) || it.matchDefault(c) {
			return nil
		}
	}
	for _, c := range ct.ulconstraint {
		if it.underlying().equals(c) || it.matchDefault(c) {
			return nil
		}
	}
	return it.node.cfgErrorf("%s does not implement %s", it.id(), ct.id())
}
