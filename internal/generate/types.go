package generate

import "go/types"

const gaffaPackagePath = "github.com/lbuening/gaffa"

func isGaffaType(t types.Type, name string, n int) bool {
	named, ok := t.(*types.Named)
	return ok &&
		named.Obj().Pkg() != nil &&
		named.Obj().Pkg().Path() == gaffaPackagePath &&
		named.Obj().Name() == name &&
		named.TypeArgs().Len() == n
}

func isGaffaRef(t types.Type) bool {
	return isGaffaType(t, "Ref", 1)
}

func isGaffaImplements(t types.Type) bool {
	return isGaffaType(t, "Implements", 1)
}

func isGaffaMain(t types.Type) bool {
	return isGaffaType(t, "Main", 0)
}

func isContext(t types.Type) bool {
	n, ok := t.(*types.Named)
	if !ok {
		return false
	}
	return n.Obj().Pkg().Path() == "context" && n.Obj().Name() == "Context"
}
