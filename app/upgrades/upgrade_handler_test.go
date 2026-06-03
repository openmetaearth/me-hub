package upgrades

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"testing"
)

func TestUpgradeHandlersDoNotOverwriteFromVersionMap(t *testing.T) {
	t.Parallel()

	files := []string{
		filepath.Join("v2_0_13", "upgrade.go"),
		filepath.Join("v2_0_14", "upgrade.go"),
		filepath.Join("v2.0.14.patch.2", "upgrade.go"),
	}

	for _, file := range files {
		file := file
		t.Run(file, func(t *testing.T) {
			t.Parallel()

			fset := token.NewFileSet()
			parsed, err := parser.ParseFile(fset, file, nil, 0)
			if err != nil {
				t.Fatalf("parse upgrade handler: %v", err)
			}

			ast.Inspect(parsed, func(node ast.Node) bool {
				assign, ok := node.(*ast.AssignStmt)
				if !ok {
					return true
				}

				for _, lhs := range assign.Lhs {
					index, ok := lhs.(*ast.IndexExpr)
					if !ok {
						continue
					}
					ident, ok := index.X.(*ast.Ident)
					if ok && ident.Name == "fromVM" {
						t.Fatalf("upgrade handler must not overwrite fromVM at %s", fset.Position(lhs.Pos()))
					}
				}

				return true
			})
		})
	}
}
