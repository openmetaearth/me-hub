package keeper

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTallyDoesNotDeleteVotesDuringIteration(t *testing.T) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, "tally.go", nil, 0)
	require.NoError(t, err)

	foundIterateVotes := false
	foundDeleteInIterator := false

	ast.Inspect(file, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok || !isKeeperCall(call.Fun, "IterateVotes") || len(call.Args) < 3 {
			return true
		}

		callback, ok := call.Args[2].(*ast.FuncLit)
		require.True(t, ok)
		foundIterateVotes = true

		ast.Inspect(callback.Body, func(callbackNode ast.Node) bool {
			callbackCall, ok := callbackNode.(*ast.CallExpr)
			if ok && isKeeperCall(callbackCall.Fun, "DeleteVote") {
				foundDeleteInIterator = true
			}
			return true
		})

		return true
	})

	require.True(t, foundIterateVotes)
	require.False(t, foundDeleteInIterator)
}

func isKeeperCall(expr ast.Expr, method string) bool {
	selector, ok := expr.(*ast.SelectorExpr)
	if !ok || selector.Sel.Name != method {
		return false
	}

	ident, ok := selector.X.(*ast.Ident)
	return ok && ident.Name == "keeper"
}
